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

type cafeteriaName struct {
	Name     string   `json:"enum_name"`
	Location location `json:"location"`
}

type cafeteriaWithID struct {
	Name      string `json:"name"`
	Cafeteria int64  `json:"cafeteria"`
}

type location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Address   string  `json:"address"`
}

type days struct {
	Days []date `json:"days"`
}

type date struct {
	Dates []dish `json:"dishes"`
}

type dish struct {
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
	var result []cafeteriaWithID
	errQueryCafeterias := c.db.Model(&model.Cafeteria{}).Select("name,cafeteria").Scan(&result).Error
	if errQueryCafeterias != nil {
		log.WithError(errQueryCafeterias).Error("Error while querying all cafeteria names from the database.")
	}

	year, week := time.Now().UTC().ISOWeek()
	var weekliesWereAdded int64
	errExistsQuery := c.db.Model(&model.DishesOfTheWeek{}).
		Where("year = ? AND week = ?", year, week).
		Count(&weekliesWereAdded).Error
	if errExistsQuery != nil {
		log.WithError(errExistsQuery).Error("Error while checking whether the meals of the current week have already been added to the weekly table.")
	}

	for _, v := range result {
		cafeteriaName := strings.Replace(strings.ToLower(v.Name), "_", "-", 10)

		req := fmt.Sprintf("https://tum-dev.github.io/eat-api/%s/%d/%d.json", cafeteriaName, year, week)
		log.WithField("req", req).Debug("Fetching menu")
		var resp, err = http.Get(req)
		if err != nil {
			log.WithError(err).Error("Error fetching menu.")
		}
		if resp.StatusCode != 200 {
			fields := log.Fields{
				"Name":       v.Name,
				"StatusCode": resp.StatusCode,
			}
			log.WithError(err).WithFields(fields).Error("Menu does not exist")
		} else {
			var dishes days
			errJson := json.NewDecoder(resp.Body).Decode(&dishes)
			if errJson != nil {
				log.WithError(err).Error("Error in Parsing")
			}

			for weekDayIndex := 0; weekDayIndex < len(dishes.Days); weekDayIndex++ {
				for u := 0; u < len(dishes.Days[weekDayIndex].Dates); u++ {
					dish := model.Dish{
						Name:        dishes.Days[weekDayIndex].Dates[u].Name,
						Type:        dishes.Days[weekDayIndex].Dates[u].DishType,
						CafeteriaID: v.Cafeteria,
					}

					var count int64
					var dishId int64
					errCount := c.db.Model(&model.Dish{}).
						Where("name = ? AND cafeteriaID = ?", dish.Name, dish.CafeteriaID).
						Select("dish").First(&dishId).
						Count(&count).Error
					if errCount != nil {
						log.WithError(errCount).Error("Error while checking whether this is already in database")
					}
					if count == 0 {
						errCreate := c.db.Model(&model.Dish{}).Create(&dish).Error
						if errCreate != nil {
							log.WithError(errCreate).Error("Error while creating new dish entry with name {}. dish won't be saved", dish.Name)
						}
						addDishTagsToMapping(dish.Dish, dish.Name, c.db)
						dishId = dish.Dish
					}
					if weekliesWereAdded == 0 {
						errCreate := c.db.Model(&model.DishesOfTheWeek{}).
							Create(&model.DishesOfTheWeek{
								DishID: dishId,
								Year:   int32(year),
								Week:   int32(week),
								Day:    int32(weekDayIndex),
							}).Error
						if errCreate != nil {
							log.WithError(errCreate).Error("Error while inserting dish for this weeks weekly dishes", dish.Name)
						}
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
	var cafeteriaNames []cafeteriaName
	errjson := json.NewDecoder(resp.Body).Decode(&cafeteriaNames)

	if errjson != nil {
		log.WithError(errjson).Error("Error while unmarshalling json data.")
	}

	for i := 0; i < len(cafeteriaNames); i++ {

		mensa := model.Cafeteria{
			Name:      cafeteriaNames[i].Name,
			Address:   cafeteriaNames[i].Location.Address,
			Latitude:  cafeteriaNames[i].Location.Latitude,
			Longitude: cafeteriaNames[i].Location.Longitude,
		}
		var cafeteriaResult model.Cafeteria
		resExists := c.db.Model(&model.Cafeteria{}).
			Where("name = ?", cafeteriaNames[i].Name).
			First(&cafeteriaResult)

		if resExists.Error != nil {
			errCreate := c.db.Model(&model.Cafeteria{}).Create(&mensa).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating the db entry for the cafeteria ", cafeteriaNames[i].Name)
			}
		} else {
			errUpdate := c.db.Model(&model.Cafeteria{}).
				Where("name = ?", cafeteriaNames[i].Name).
				Updates(&mensa).Error
			if errUpdate != nil {
				log.WithError(errUpdate).Error("Error while updating the db entry for the cafeteria {}.", cafeteriaNames[i].Name)
			}
		}
	}
}

// addDishTagsToMapping
// Checks whether the dish name includes one of the expressions for the excluded tags as well as the included tags.
// The corresponding tags for all identified DishNames will be saved in the table DishNameTags.
func addDishTagsToMapping(dishID int64, dishName string, db *gorm.DB) {
	lowercaseDish := strings.ToLower(dishName)
	var includedTags []int64
	errIncluded := db.Model(&model.DishNameTagOptionIncluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&includedTags).Error
	if errIncluded != nil {
		log.WithError(errIncluded).Error("Error while querying all included expressions for the dish: ", lowercaseDish)
	}

	var excludedTags []int64
	err := db.Model(&model.DishNameTagOptionExcluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&excludedTags).Error
	if err != nil {
		log.WithError(err).Error("Error while querying all excluded expressions for the dish: ", lowercaseDish)
	}

	//set all entries in included to -1 if the excluded tag was recognised for this tag rating.
	if len(excludedTags) > 0 {
		for _, a := range excludedTags {
			i := contains(includedTags, a)
			if i != -1 {
				includedTags[i] = -1
			}
		}
	}

	for _, nametagID := range includedTags {
		if nametagID != -1 {
			err := db.Model(&model.DishToDishNameTag{}).Create(&model.DishToDishNameTag{
				DishID:    dishID,
				NameTagID: nametagID,
			}).Error
			if err != nil {
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
