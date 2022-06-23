package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	"log"
)

type averageRatingForCafeteria struct {
	CafeteriaID int32   `json:"cafeteria"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
}

type averageRatingForMealInCafeteria struct {
	CafeteriaID int32   `json:"cafeteria"`
	MealID      int32   `json:"meal"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
}

//regularly computes the average rating for every cafeteria
func (c *CronService) averageRatingComputation() error {

	computeAverageForCafeteria(c)
	computeAverageForMealsInCafeterias(c)
	computeAverageCafeteriaTags(c)
	return nil
}

func computeAverageCafeteriaTags(c *CronService) {

	/*Todo
				alle ratings einer cafeteria sammeln -> nur die IDs merken, diese dann auf das tagrating tabelle anwenden - join
				-> alle tagratings zu einer cafeteria in einer Gruppe
				-> pro gruppe nach den tags gruppieren und den Durchschnitt bereschnen und in einemr result tabelle Speichern


			für alle drei tagarten berechnen
		nameratingtag ist das komplizierte, die beiden anderen können auf der tagrating tabelle bestimmt werden
	-> zusammenführen durch das parent rating um zu erfahren, zu welcher mensa die gerichte gehören
	*/

	/*
		err := s.db.Raw("SELECT r.*, a.campus, a.name "+
				"FROM roomfinder_rooms r "+
				"LEFT JOIN roomfinder_building2area a ON a.building_nr = r.building_nr "+
				"WHERE MATCH(room_code, info, address) AGAINST(?)", req.Query).Scan(&res).Error
	*/

	//nach der tagID gruppieren

	/*c.db.Model(&cafeteria_rating_models.CafeteriaRating{}).
	Select("id,rating,cafeteria").
	Joins("left join emails on emails.user_id = users.id").
	Scan(&result{})
	*/
	/*res, err := c.db.Model(cafeteria_rating_models.MealRatingsTags{}).
	Select("cafeteria, meal, AVG(rating) as average, MAX(rating) as max, MIN(rating) as min").
	Group("cafeteria,meal").Joins().Rows()
	*/
	/*
			Schtitte; erstmal das jion verstehen
		Meal anme tags passen noch nciht ganz -> sollten final keinen namen enthalten, sondern nur den key
	*/

	/*if err != nil {
		println("Error in query")
	}

	println(res.ColumnTypes())
	*/
}

func computeAverageForMealsInCafeterias(c *CronService) {
	var results []averageRatingForMealInCafeteria
	res := c.db.Model(cafeteria_rating_models.MealRating{}).
		Select("cafeteria, meal, AVG(rating) as average, MAX(rating) as max, MIN(rating) as min").
		Group("cafeteria,meal").Find(&results)

	if res.Error != nil {
		log.Println("Error in query")
		log.Println(res.Error)
	} else {
		for _, v := range results {
			cafeteria := cafeteria_rating_models.MealRatingsAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     float32(v.Average),
				MealID:      v.MealID,
				Min:         v.Min,
				Max:         v.Max,
			} //todo add standard deviation

			var existing *cafeteria_rating_models.MealRatingsAverage
			testDish := c.db.Model(cafeteria_rating_models.MealRatingsAverage{}).
				Where("cafeteria = ?", cafeteria.CafeteriaID).
				Where("meal = ?", cafeteria.MealID).
				First(&existing)

			if testDish.RowsAffected == 1 {
				errUpdate := c.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).
					Where("cafeteria = ?", cafeteria.CafeteriaID).
					Where("meal = ?", cafeteria.MealID).
					Updates(cafeteria)

				if errUpdate.Error != nil {
					log.Println(errUpdate.Error)
				}
			} else {
				log.Println("New average rating will be created for cafeteria with ID: ", v.CafeteriaID)
				errCreate := c.db.Create(&cafeteria)
				if errCreate.Error != nil {
					log.Println(errCreate.Error)
				}
			}
		}
	}
}

func computeAverageForCafeteria(c *CronService) {
	var results []averageRatingForCafeteria
	res := c.db.Model(cafeteria_rating_models.CafeteriaRating{}).
		Select("cafeteria, AVG(rating) as average, MAX(rating) as max, MIN(rating) as min").
		Group("cafeteria").Find(&results)

	if res.Error != nil {
		log.Println("Error in query")
		log.Println(res.Error)
	} else {
		for _, v := range results {
			cafeteria := cafeteria_rating_models.CafeteriaRatingsAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				Min:         v.Min,
				Max:         v.Max,
			} //todo add standard deviation

			var existing *cafeteria_rating_models.CafeteriaRatingsAverage
			testDish := c.db.Model(cafeteria_rating_models.CafeteriaRatingsAverage{}).Where("cafeteria = ?", cafeteria.CafeteriaID).First(&existing)

			if testDish.RowsAffected == 1 {
				errUpdate := c.db.Model(&cafeteria_rating_models.CafeteriaRatingsAverage{}).
					Where("cafeteria = ?", cafeteria.CafeteriaID).
					Updates(cafeteria)

				if errUpdate.Error != nil {
					log.Println(errUpdate.Error)
				}
			} else {
				log.Println("New rating will be created for cafeteria: ", v.CafeteriaID)
				errCreate := c.db.Create(&cafeteria)
				if errCreate.Error != nil {
					log.Println(errCreate.Error)
				}
			}
		}
	}
}
