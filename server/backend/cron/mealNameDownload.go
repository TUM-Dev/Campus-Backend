package cron

import (
	"encoding/json"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CafeteriaName struct {
	Name     string   `json:"enum_name"`
	Location Location `json:"location"`
}

type CafeteriaWithID struct {
	Name      string `json:"name"`
	Cafeteria int32  `json:"cafeteria"`
}

type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Address   string  `json:"address"`
}

type Days struct {
	Days []Date `json:"days"`
}

type Date struct {
	Dates []Dish `json:"dishes"`
}

type Dish struct {
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
	c.db.Model(&model.Cafeteria{}).Select("name,id").Scan(&result)

	for _, v := range result {

		cafeteriaName := strings.Replace(strings.ToLower(v.Name), "_", "-", 10)
		y, w := time.Now().UTC().ISOWeek()

		req := "https://tum-dev.github.io/eat-api/" + cafeteriaName + "/" + strconv.Itoa(y) + "/" + strconv.Itoa(w) + ".json"
		log.Info("Fetching menu from: {}", req)
		var resp, err = http.Get(req)
		if err != nil {
			log.WithError(err).Error("Error fetching menu.")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if body[0] == '<' {
			log.WithError(err).Error("Menu for", v, "does not exist error 404 returned.")
		} else {
			var dishes Days
			errjson := json.Unmarshal(body, &dishes)
			if errjson != nil {
				log.WithError(err).Error("Error in Parsing")
			}
			for i := 0; i < len(dishes.Days); i++ {
				for u := 0; u < len(dishes.Days[i].Dates); u++ {

					dish := model.Dish{
						Name:        dishes.Days[i].Dates[u].Name,
						Type:        dishes.Days[i].Dates[u].DishType,
						CafeteriaID: v.Cafeteria,
					}

					var count int64
					errCount := c.db.Model(&model.Dish{}).
						Where("name = ? AND cafeteriaID = ?", dish.Name, dish.CafeteriaID).Count(&count).Error
					if errCount != nil {
						log.WithError(errCount).Error("Error while checking whether dis is already in database")
					}
					if count == 0 {
						errCreate := c.db.Model(&model.Dish{}).Create(&dish).Error
						if errCreate != nil {
							log.WithError(errCreate).Error("Error while creating new dish entry with name {}. Dish won't be saved", dish.Name)
						}
						addDishTagsToMapping(dish.Dish, dish.Name, c.db)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Error reading json data.")
	}

	var cafeteriaNames []CafeteriaName
	errjson := json.Unmarshal(body, &cafeteriaNames)
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
			log.WithError(resExists.Error).Error("Error while unmarshalling json data.")
		}

		if resExists.RowsAffected == 0 {
			errCreate := c.db.Model(&model.Cafeteria{}).Create(&mensa).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating the db entry for the cafeteria {}", cafeteriaNames[i].Name)
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
func addDishTagsToMapping(dishID int32, dishName string, db *gorm.DB) {
	lowercaseDish := strings.ToLower(dishName)
	var includedTags []int32
	errIncluded := db.Model(&model.DishNameTagOptionIncluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&includedTags).Error
	if errIncluded != nil {
		log.WithError(errIncluded).Error("Error while querying all included expressions for the dish: {}", lowercaseDish)
	}

	var excludedTags []int32
	errExcluded := db.Model(&model.DishNameTagOptionExcluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&excludedTags).Error
	if errExcluded != nil {
		log.WithError(errExcluded).Error("Error while querying all excluded expressions for the dish: {}", lowercaseDish)
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

	for _, a := range includedTags {
		if a != -1 {
			err := db.Model(&model.DishToDishNameTag{}).Create(&model.DishToDishNameTag{
				DishID:    dishID,
				NameTagID: a,
			}).Error
			if err != nil {
				log.WithError(err).Error("Error while creating a new entry with dish {} and nametag {}", dishID, a)
			}
		}
	}
}
func contains(s []int32, e int32) int32 {
	for i, a := range s {
		if a == e {
			return int32(i)
		}
	}
	return -1
}
