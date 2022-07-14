package cron

import (
	"encoding/json"
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
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
	Dates []Meal `json:"dishes"`
}

type Meal struct {
	Name     string `json:"name"`
	DishType string `json:"dish_type"`
}

//fileDownloadCron Downloads all files that are not marked as finished in the database.
func (c *CronService) mealNameDownloadCron() error {

	downloadCanteenNames(c)
	downloadDailyMeals(c)

	return nil
}

func downloadDailyMeals(c *CronService) {
	var result []CafeteriaWithID
	c.db.Model(&cafeteria_rating_models.Cafeteria{}).Select("name,id").Scan(&result)

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
			log.Println("Mealplan for", v, "does not exist error 404 returned.")
		} else {
			var meals Days
			errjson := json.Unmarshal(body, &meals)
			if errjson != nil {
				log.Println("Error in Parsing")
				log.Fatalln(errjson)
			}
			log.Println("Meals:")
			for i := 0; i < len(meals.Days); i++ {
				for u := 0; u < len(meals.Days[i].Dates); u++ {

					meal := cafeteria_rating_models.Meal{
						Name:        meals.Days[i].Dates[u].Name,
						Type:        meals.Days[i].Dates[u].DishType,
						CafeteriaID: v.Cafeteria,
					}

					res := c.db.Model(&cafeteria_rating_models.Meal{}).
						Where("name = ? AND cafeteriaID = ?", meal.Name, meal.CafeteriaID)

					if res.RowsAffected == 0 {
						c.db.Model(&cafeteria_rating_models.Meal{}).Create(&meal)
						addMealTagsToMapping(meal.Meal, meal.Name, c.db)
					} /*else {		//todo potentially add update logic for the weekly meals
						c.db.Model(&cafeteria_rating_models.Cafeteria{}).
							Where("name = ?", cafeteriaNames[i].Name).
							Updates(&mensa)
					}*/
					//c.db.Model(&cafeteria_rating_models.Meal{}).Create(&meal)
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

		mensa := cafeteria_rating_models.Cafeteria{
			Name:      cafeteriaNames[i].Name,
			Address:   cafeteriaNames[i].Location.Address,
			Latitude:  cafeteriaNames[i].Location.Latitude,
			Longitude: cafeteriaNames[i].Location.Longitude,
		}
		var cafetriaResult cafeteria_rating_models.Cafeteria
		res := c.db.Model(&cafeteria_rating_models.Cafeteria{}).
			Where("name = ?", cafeteriaNames[i].Name).
			First(&cafetriaResult)
		if res.RowsAffected == 0 {
			c.db.Model(&cafeteria_rating_models.Cafeteria{}).Create(&mensa)
		} else {
			c.db.Model(&cafeteria_rating_models.Cafeteria{}).
				Where("name = ?", cafeteriaNames[i].Name).
				Updates(&mensa)
		}
	}
}

/*
Checks whether the meal name includes one of the expressions for the excluded tas as well as the included tags.
The corresponding tags for all identified MealNames will be saved in the table MealNameTags.
*/
func addMealTagsToMapping(mealID int32, mealName string, db *gorm.DB) {
	lowercaseMeal := strings.ToLower(mealName)
	var includedTags []int32
	db.Model(&cafeteria_rating_models.MealNameTagOptionIncluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseMeal).
		Select("nameTagID").
		Scan(&includedTags)

	var excludedTags []int32
	db.Model(&cafeteria_rating_models.MealNameTagOptionExcluded{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseMeal).
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
			db.Model(&cafeteria_rating_models.MealToMealNameTag{}).Create(&cafeteria_rating_models.MealToMealNameTag{
				MealID:    mealID,
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
