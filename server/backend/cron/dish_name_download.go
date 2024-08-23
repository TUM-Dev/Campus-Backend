package cron

import (
	"encoding/json"
	"errors"
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
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
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
	var results []model.Canteen
	if err := c.db.Find(&results).Error; err != nil {
		log.WithError(err).Error("Error while querying all cafeteria names from the database.")
		return
	}

	year, week := time.Now().UTC().ISOWeek()
	var weekliesWereAdded int64

	if err := c.db.Model(&model.DishesOfTheWeek{}).
		Where("year = ? AND week = ?", year, week).
		Count(&weekliesWereAdded).Error; err != nil {
		log.WithError(err).Error("Error while checking whether the meals of the current week have already been added to the weekly table.")
		return
	}

	for _, v := range results {
		cafeteriaName := strings.Replace(strings.ToLower(v.Name), "_", "-", 10)

		req := fmt.Sprintf("https://tum-dev.github.io/eat-api/%s/%d/%02d.json", cafeteriaName, year, week)
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
			// sometimes the eat-api does not have the required data.
			// This might be because of a lot of factors, but most commonly no menu was posted for a certain week at the time of scraping
			log.WithError(err).WithFields(fields).Info("Menu does not exist")
			continue
		}
		var dishes CanteenDays
		if err := json.NewDecoder(resp.Body).Decode(&dishes); err != nil {
			log.WithError(err).Error("Error in Parsing")
			return
		}

		for weekDayIndex, day := range dishes.Days {
			for _, date := range day.Dates {
				dish := model.Dish{
					Name:        date.Name,
					Type:        date.DishType,
					CafeteriaID: v.Cafeteria,
				}

				var dbDish model.Dish
				if err := c.db.First(&dbDish, "name = ? AND cafeteriaID = ?", dish.Name, dish.CafeteriaID).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
					if err := c.db.Create(&dish).Error; err != nil {
						log.WithError(err).WithField("name", dish.Name).Error("Error while creating new CanteenDish entry. CanteenDish won't be saved")
					}
					addDishTagsToMapping(dish.Dish, dish.Name, c.db)
					dbDish = dish
				} else if err != nil {
					log.WithError(err).Error("Error while checking whether the dish is already in database")
				}

				if dbDish.Type != dish.Type {
					if err := c.db.Where("dish = ?", dbDish.Dish).Updates(&dish).Error; err != nil {
						log.WithError(err).WithField("from", dish.Type).WithField("to", dish.Type).Error("Error while updating dish to new type")
					}
				}

				if weekliesWereAdded == 0 {
					errCreate := c.db.Create(&model.DishesOfTheWeek{
						DishID: dbDish.Dish,
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
		mensa := model.Canteen{
			Name:      cafeteriaName.Name,
			Address:   cafeteriaName.Location.Address,
			Latitude:  cafeteriaName.Location.Latitude,
			Longitude: cafeteriaName.Location.Longitude,
		}
		var cafeteriaResult model.Canteen
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
	if err := db.Model(&model.IncludedDishNameTagOption{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&includedTags).Error; err != nil {
		log.WithError(err).Error("Error while querying all included expressions for the CanteenDish: ", lowercaseDish)
	}

	var excludedTags []int64
	if err := db.Model(&model.ExcludedDishNameTagOption{}).
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
