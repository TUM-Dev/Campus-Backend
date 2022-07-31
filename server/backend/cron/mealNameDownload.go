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

//fileDownloadCron Downloads all files that are not marked as finished in the database.
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
		log.Println(cafeteriaName)
		y, w := time.Now().UTC().ISOWeek()

		req := "https://tum-dev.github.io/eat-api/" + cafeteriaName + "/" + strconv.Itoa(y) + "/" + strconv.Itoa(w) + ".json"
		println(req)
		var resp, err = http.Get(req)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if body[0] == '<' {
			log.Println("Dishplan for", v, "does not exist error 404 returned.")
		} else {
			var dishes Days
			errjson := json.Unmarshal(body, &dishes)
			if errjson != nil {
				log.Println("Error in Parsing")
				log.Fatalln(errjson)
			}
			log.Println("Dishes:")
			for i := 0; i < len(dishes.Days); i++ {
				for u := 0; u < len(dishes.Days[i].Dates); u++ {

					dish := model.Dish{
						Name:        dishes.Days[i].Dates[u].Name,
						Type:        dishes.Days[i].Dates[u].DishType,
						CafeteriaID: v.Cafeteria,
					}

					res := c.db.Model(&model.Dish{}).
						Where("name = ? AND cafeteriaID = ?", dish.Name, dish.CafeteriaID)

					if res.RowsAffected == 0 {
						c.db.Model(&model.Dish{}).Create(&dish)
						addDishTagsToMapping(dish.Dish, dish.Name, c.db)
					} /*else {		//todo potentially add update logic for the weekly dishes
						c.db.Model(&cafeteria_rating_models.Cafeteria{}).
							Where("name = ?", cafeteriaNames[i].Name).
							Updates(&mensa)
					}*/
					//c.db.Model(&cafeteria_rating_models.Dish{}).Create(&dish)
				}
			}
		}
	}
}

func downloadCanteenNames(c *CronService) {
	var resp, err = http.Get("https://tum-dev.github.io/eat-api/enums/canteens.json")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var cafeteriaNames []CafeteriaName
	errjson := json.Unmarshal(body, &cafeteriaNames)
	if errjson != nil {
		log.Println("Error in Parsing")
		log.Fatalln(errjson)
	}

	for i := 0; i < len(cafeteriaNames); i++ {

		mensa := model.Cafeteria{
			Name:      cafeteriaNames[i].Name,
			Address:   cafeteriaNames[i].Location.Address,
			Latitude:  cafeteriaNames[i].Location.Latitude,
			Longitude: cafeteriaNames[i].Location.Longitude,
		}
		var cafetriaResult model.Cafeteria
		res := c.db.Model(&model.Cafeteria{}).
			Where("name = ?", cafeteriaNames[i].Name).
			First(&cafetriaResult)
		if res.RowsAffected == 0 {
			c.db.Model(&model.Cafeteria{}).Create(&mensa)
		} else {
			c.db.Model(&model.Cafeteria{}).
				Where("name = ?", cafeteriaNames[i].Name).
				Updates(&mensa)
		}
	}
}

/*
Checks whether the dish name includes one of the expressions for the excluded tas as well as the included tags.
The corresponding tags for all identified DishNames will be saved in the table DishNameTags.
*/
func addDishTagsToMapping(dishID int32, dishName string, db *gorm.DB) {
	lowercaseDish := strings.ToLower(dishName)
	var includedTags []int32
	db.Model(&model.DishNameTagOptionIncluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&includedTags)

	var excludedTags []int32
	db.Model(&model.DishNameTagOptionExcluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&excludedTags)

	//log.Println("Number of included tags: ", len(includedTags))

	//set all entries in included to -1 if the excluded tag was recognised ffor this tag rating.
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
			db.Model(&model.DishToDishNameTag{}).Create(&model.DishToDishNameTag{
				DishID:    dishID,
				NameTagID: a,
			})
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
