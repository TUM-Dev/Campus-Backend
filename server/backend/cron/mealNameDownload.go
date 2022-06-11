package cron

import (
	"encoding/json"
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

//fileDownloadCron Downloads all files that are not marked as finished in the database.
func (c *CronService) mealNameDownloadCron() error {

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

	for i := 0; i < len(canteens); i++ {
		log.Println(canteens[i].Name)

		mensa := model.Mensa{Id: int32(i), Name: canteens[i].Name, Address: canteens[i].Location.Address, Latitude: canteens[i].Location.Latitude, Longitude: canteens[i].Location.Longitude}
		if c.db.Model(&mensa).Where("name = ?", canteens[i].Name).Updates(&mensa).RowsAffected == 0 {
			c.db.Create(&mensa)
			log.Println("Created")
		}

		//c.db.Model(&mensa).Where("name = ?", canteens[i].Name).Update("address", "test")
	}

	//c.db.Model(&model.{})
	//weiter dish to mensa etc ausfüllen
	//dihes of the day lin db
	/*
	   Save to file:

	   	f, err := os.OpenFile("backend/static_data/cafeteriaNames.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	   	if err != nil {
	   		log.Fatal(err)
	   	}
	   	for i := 0; i < len(canteens); i++ {
	   		log.Println(canteens[i].Name)
	   		//log.Println(canteens[i].Location.Longitude)
	   		_, err := f.WriteString(canteens[i].Name)
	   		if err != nil {
	   			return err
	   		}
	   		//y[i] = tags.MultiLanguageTags[i].TagNameEnglish
	   	}


	   	f.Close()
	*/

	/*var tags CafeteriaNameList
	json.Unmarshal(body, &tags)

	var helper = len(tags.CafeteriaNameList)
	y := make([]string, helper)
	for i := 0; i < len(tags.CafeteriaName); i++ {
		y[i] = tags.MultiLanguageTags[i].TagNameEnglish
	}*/
	//	log.Println(sb)
	//	log.Println("mealnameDownloader called")
	return nil
}

type CafeteriaName struct {
	Name     string   `json:"enum_name"`
	Location Location `json:"location"`
}

type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Address   string  `json:"address"`
}

func (c *CronService) mealNameDownload(file model.Files) {
	log.Println("mealnameDownload called")
}
