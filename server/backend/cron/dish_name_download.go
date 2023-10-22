package cron

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CafeteriaName struct {
	Name     string          `json:"enum_name"`
	Location CanteenLocation `json:"location"`
}

type CafeteriaWithID struct {
	Name      string `json:"name"`
	Cafeteria int64  `json:"cafeteria"`
}

type CanteenLocation struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Address   string  `json:"address"`
}

type CanteenDays struct {
	Days []CanteenDate `json:"days"`
}

type CanteenDate struct {
	Dates []CanteenDish `json:"dishes"`
}

type CanteenDish struct {
	Name     string `json:"name"`
	DishType string `json:"dish_type"`
}

// fileDownloadCron
// Downloads all files that are not marked as finished in the database.
func (c *CronService) dishNameDownloadCron() error {
	downloadCanteenNames(c)
	downloadDailyDishes(c)
	return nil
}

func downloadDailyDishes(c *CronService) {
	var result []CafeteriaWithID
	if err := c.db.Model(&model.Cafeteria{}).Select("name,cafeteria").Scan(&result).Error; err != nil {
		log.WithError(err).Error("Error while querying all cafeteria names from the database.")
	}

	year, week := time.Now().UTC().ISOWeek()
	var weekliesWereAdded int64

	if err := c.db.Model(&model.DishesOfTheWeek{}).
		Where("year = ? AND week = ?", year, week).
		Count(&weekliesWereAdded).Error; err != nil {
		log.WithError(err).Error("Error while checking whether the meals of the current week have already been added to the weekly table.")
	}

	for _, v := range result {
		cafeteriaName := strings.Replace(strings.ToLower(v.Name), "_", "-", 10)

		req := fmt.Sprintf("https://tum-dev.github.io/eat-api/%s/%d/%d.json", cafeteriaName, year, week)
		log.WithField("req", req).Debug("Fetching menu")
		var resp, err = http.Get(req)
		if err != nil {
			log.WithError(err).Error("Error fetching menu.")
			continue
		}
		if resp.StatusCode != 200 {
			fields := log.Fields{
				"Name":       v.Name,
				"StatusCode": resp.StatusCode,
			}
			log.WithError(err).WithFields(fields).Error("Menu does not exist")
			continue
		}
		var dishes CanteenDays
		if err := json.NewDecoder(resp.Body).Decode(&dishes); err != nil {
			log.WithError(err).Error("Error in Parsing")
		}

		for weekDayIndex, day := range dishes.Days {
			for _, date := range day.Dates {
				dish := model.Dish{
					Name:        date.Name,
					Type:        date.DishType,
					CafeteriaID: v.Cafeteria,
				}

				var count int64
				var dishId int64
				if err := c.db.Model(&model.Dish{}).
					Where("name = ? AND cafeteriaID = ?", dish.Name, dish.CafeteriaID).
					Select("CanteenDish").First(&dishId).
					Count(&count).Error; err != nil {
					log.WithError(err).Error("Error while checking whether this is already in database")
				}
				if count == 0 {
					if err := c.db.Create(&dish).Error; err != nil {
						log.WithError(err).Error("Error while creating new CanteenDish entry with name {}. CanteenDish won't be saved", dish.Name)
					}
					addDishTagsToMapping(dish.Dish, dish.Name, c.db)
					dishId = dish.Dish
				}
				if weekliesWereAdded == 0 {
					errCreate := c.db.Create(&model.DishesOfTheWeek{
						DishID: dishId,
						Year:   int32(year),
						Week:   int32(week),
						Day:    int32(weekDayIndex),
					}).Error
					if errCreate != nil {
						log.WithError(errCreate).Error("Error while inserting CanteenDish for this weeks weekly dishes", dish.Name)
					}
				}
			}
		}
	}
}

func downloadCanteenNames(c *CronService) {
	var resp, err = http.Get("https://tum-dev.github.io/eat-api/enums/canteens.json")
	if err != nil {
		log.WithError(err).Error("Error fetching cafeteria list from eat-api.")
	}
	var cafeteriaNames []CafeteriaName
	if err := json.NewDecoder(resp.Body).Decode(&cafeteriaNames); err != nil {
		log.WithError(err).Error("Error while unmarshalling json data.")
	}

	for _, cafeteriaName := range cafeteriaNames {
		mensa := model.Cafeteria{
			Name:      cafeteriaName.Name,
			Address:   cafeteriaName.Location.Address,
			Latitude:  cafeteriaName.Location.Latitude,
			Longitude: cafeteriaName.Location.Longitude,
		}
		var cafeteriaResult model.Cafeteria
		if err := c.db.First(&cafeteriaResult, "name = ?", cafeteriaName.Name).Error; err != nil {
			if err := c.db.Create(&mensa).Error; err != nil {
				log.WithError(err).Error("Error while creating the db entry for the cafeteria ", cafeteriaName.Name)
			}
		} else {
			if err := c.db.Where("name = ?", cafeteriaName.Name).Updates(&mensa).Error; err != nil {
				log.WithError(err).Error("Error while updating the db entry for the cafeteria {}.", cafeteriaName.Name)
			}
		}
	}
}

// addDishTagsToMapping
// Checks whether the CanteenDish name includes one of the expressions for the excluded tags as well as the included tags.
// The corresponding tags for all identified DishNames will be saved in the table DishNameTags.
func addDishTagsToMapping(dishID int64, dishName string, db *gorm.DB) {
	lowercaseDish := strings.ToLower(dishName)
	var includedTags []int64
	if err := db.Model(&model.DishNameTagOptionIncluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&includedTags).Error; err != nil {
		log.WithError(err).Error("Error while querying all included expressions for the CanteenDish: ", lowercaseDish)
	}

	var excludedTags []int64
	if err := db.Model(&model.DishNameTagOptionExcluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&excludedTags).Error; err != nil {
		log.WithError(err).Error("Error while querying all excluded expressions for the CanteenDish: ", lowercaseDish)
	}

	//set all entries in included to -1 if the excluded tag was recognised for this tag rating.
	if len(excludedTags) > 0 {
		for _, a := range excludedTags {
			if i := contains(includedTags, a); i != -1 {
				includedTags[i] = -1
			}
		}
	}

	for _, nametagID := range includedTags {
		if nametagID != -1 {
			if err := db.Create(&model.DishToDishNameTag{
				DishID:    dishID,
				NameTagID: nametagID,
			}).Error; err != nil {
				fields := log.Fields{"dishID": dishID, "nametagID": nametagID}
				log.WithError(err).WithFields(fields).Error("creating a new entry")
			}
		}
	}
}
func contains(s []int64, e int64) int64 {
	for i, a := range s {
		if a == e {
			return int64(i)
		}
	}
	return -1
}
