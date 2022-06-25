package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	"log"
)

type averageRatingForCafeteria struct {
	CafeteriaID int32   `json:"cafeteriaID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
}

type averageRatingForMealInCafeteria struct {
	CafeteriaID int32   `json:"cafeteriaID"`
	MealID      int32   `json:"mealID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
}

type averageCafeteriaTags struct {
	cafeteria_rating_models.Cafeteria
	Id          int
	CafeteriaID int `gorm:"ForeignKey:mensaID"`
	Rating      int
}

type averageCafeteriaTagsTest struct {
	CafeteriaID int32   `json:"cafeteria_rating.cafeteriaID"`
	TagID       int32   `json:"cafeteria_rating_tags.tagID"`
	Average     float32 `json:"cafeteria_rating_tags.rating"`
}

//regularly computes the average rating for every cafeteria
func (c *CronService) averageRatingComputation() error {

	//computeAverageForCafeteria(c)
	//computeAverageForMealsInCafeterias(c)
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

	/*
		todo
		cafeteriarating x prating tags um anhand der parent rating idceafeteriRating.id jedem tagrating eine mensa id zuzuordnen, dann zweimal gruppieren
		(erst nach mensa,dann anch dem tag) und jweeils den durchschnitt bestimmen


		join:
		aus dem normalen rating muss nur die id as cafeteriaID übernimen werden, aus dem rating_tags alles bis auf das parent rating
	*/

	//cafeteria_rating_tags.rating, cafeteria_rating_tags.id,
	/*	var results []averageCafeteriaTags
		err := c.db.Model(&cafeteria_rating_models.CafeteriaRatingTags{}).
			Select("cafeteria_rating.cafeteriaID as cafeteriaID,AVG(cafeteria_rating_tags.rating) as average").
			Joins("JOIN cafeteria_rating ON cafeteria_rating.Id = cafeteria_rating_tags.parentRating").
			Group("cafeteria_rating_tags.tagID").
			Find(&results)*/

	//todo erstmal die tabelle vorbereiten, dann daraus mit dem average querien.
	/*var results []averageCafeteriaTags
	err := c.db.Model(&cafeteria_rating_models.CafeteriaRatingTags{}).
		Select("cafeteria_rating_tags.tagID as tagID,AVG(cafeteria_rating_tags.rating) as average").
		Joins("JOIN cafeteria_rating ON cafeteria_rating.Id = cafeteria_rating_tags.parentRating").
		Group("tagID,cafeteria_rating.cafeteriaID").
		Scan(&results)*/

	/*	db := c.db.Model(&cafeteria_rating_models.CafeteriaRatingTags{}).
			Select("cafeteria_rating_tags.*, cafeteria_rating.*").
			Joins("JOIN cafeteria_rating ON cafeteria_rating.Id = cafeteria_rating_tags.parentRating").
			Group("tagID,cafeteria_rating.cafeteriaID")
		//zweiter teil des groups: cafeteria_rating.cafeteriaID,cafeteria_rating_tags.id

		var result []averageCafeteriaTagsTest
		test := db.First(&result)
		println(test.Error)
		println(result)*/

	/*var res []averageCafeteriaTags
	err := c.db.Raw("SELECT cr.cafeteriaID, crt.tagID, AVG(crt.rating) as average " +
		"FROM cafeteria_rating cr " +
		"JOIN cafeteria_rating_tags crt ON cr.id = crt.parentRating" +
		"GROUP BY cr.cafeteriaID").Scan(&res).Error*/

	var res []averageCafeteriaTagsTest
	err := c.db.Raw("SELECT cafeteriaID, tagID, AVG(rating) as average" +
		" FROM (SELECT * FROM cafeteria_rating cr" +
		" JOIN cafeteria_rating_tags crt ON cr.id = crt.parentRating) table").Scan(&res).Error

	/*
		+
				" GROUP BY cafeteriaID"
	*/
	/*	var res []averageCafeteriaTags
		err := c.db.Debug().Raw("SELECT cr.id, cr.cafeteriaID, cr.rating" +
			" FROM cafeteria_rating cr" +
			" JOIN cafeteria_rating_tags crt ON cr.id = crt.parentRating) table" +
			" WHERE cr.cafeteriaID = 1").Scan(&res)
	*/
	/*
		todo lässt es sich nur nicht auslesen, da es ein foreing key ist?
	*/
	/*
		 +
				"JOIN cafeteria_rating_tags crt ON cr.id = crt.parentRating"
	*/
	if err != nil {
		log.Println(err)
	}
}

func computeAverageForMealsInCafeterias(c *CronService) {
	var results []averageRatingForMealInCafeteria
	res := c.db.Model(&cafeteria_rating_models.MealRating{}).
		Select("cafeteriaID, mealID, AVG(rating) as average, MAX(rating) as max, MIN(rating) as min").
		Group("cafeteriaID,mealID").Scan(&results)

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
			testDish := c.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).
				Where("cafeteriaID = ? AND mealID = ?", cafeteria.CafeteriaID, cafeteria.MealID).
				First(&existing)

			if testDish.RowsAffected == 1 {
				errUpdate := c.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).
					Where("cafeteriaID = ? AND mealID = ?", cafeteria.CafeteriaID, cafeteria.MealID).
					Updates(cafeteria)

				if errUpdate.Error != nil {
					log.Println(errUpdate.Error)
				}
			} else {
				log.Println("New average rating will be created for cafeteria with ID: ", v.CafeteriaID)
				errCreate := c.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).Create(&cafeteria)
				if errCreate.Error != nil {
					log.Println(errCreate.Error)
				}
			}
		}
	}
}

func computeAverageForCafeteria(c *CronService) {
	var results []averageRatingForCafeteria
	res := c.db.Model(&cafeteria_rating_models.CafeteriaRating{}).
		Select("cafeteriaID, AVG(rating) as average, MAX(rating) as max, MIN(rating) as min").
		Group("cafeteriaID").Find(&results)

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
			testDish := c.db.Model(&cafeteria_rating_models.CafeteriaRatingsAverage{}).Where("cafeteriaID = ?", cafeteria.CafeteriaID).First(&existing)

			if testDish.RowsAffected == 1 {
				errUpdate := c.db.Model(&cafeteria_rating_models.CafeteriaRatingsAverage{}).
					Where("cafeteriaID = ?", cafeteria.CafeteriaID).
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
