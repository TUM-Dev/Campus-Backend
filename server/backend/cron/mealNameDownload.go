package cron

import (
	"encoding/json"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
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

	c.db.Where("1=1").Delete(&model.Dish{}) //Remove all meals of the previous week

	var result []string
	c.db.Table("mensa").Select("name").Scan(&result)
	for _, v := range result {

		canteenName := strings.Replace(strings.ToLower(v), "_", "-", 10)
		log.Println(canteenName)
		y, w := time.Now().UTC().ISOWeek()

		req := "https://tum-dev.github.io/eat-api/" + canteenName + "/" + strconv.Itoa(y) + "/" + strconv.Itoa(w) + ".json"
		println(req)
		var resp, err = http.Get(req)
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)

		var meals Days
		errjson := json.Unmarshal(body, &meals)
		if errjson != nil {
			log.Println("Error in Parsing")
			log.Fatalln(errjson)
		}
		log.Println("Meals:")
		for i := 0; i < len(meals.Days); i++ {
			for u := 0; u < len(meals.Days[i].Dates); u++ {
				meal := model.Dish{Name: meals.Days[i].Dates[u].Name, Type: meals.Days[i].Dates[u].DishType, Canteen: v}
				c.db.Create(&meal)
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

	var canteens []CafeteriaName
	errjson := json.Unmarshal(body, &canteens)
	if errjson != nil {
		log.Println("Error in Parsing")
		log.Fatalln(errjson)
	}

	// store canteen information in mensa db
	for i := 0; i < len(canteens); i++ {
		mensa := model.Mensa{Id: int32(i), Name: canteens[i].Name, Address: canteens[i].Location.Address, Latitude: canteens[i].Location.Latitude, Longitude: canteens[i].Location.Longitude}
		if c.db.Model(&mensa).Where("name = ?", canteens[i].Name).Updates(&mensa).RowsAffected == 0 {
			c.db.Create(&mensa)
		}
	}

}
