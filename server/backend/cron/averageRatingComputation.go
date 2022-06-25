package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	"gorm.io/gorm"
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
	gorm.Model
	Id           int
	CafeteriaID  int `gorm:"foreignKey:cafeteriaId"`
	Rating       int
	Comment      string
	ParentRating int `gorm:"ForeignKey:id"`
	TagID        int
}

type averageCafeteriaTagsTest struct {
	Id           int32 `gorm:"primary_key;AUTO_INCREMENT;column:id;type:int;" json:"id"`
	CafeteriaID  int32 `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	Rating       int
	Comment      string
	ParentRating int32 `gorm:"foreignKey:cafeteriaRatingID;column:parentRating;type:int;" json:"parentRating"`
}

//regularly computes the average rating for every cafeteria
func (c *CronService) averageRatingComputation() error {

	//computeAverageForCafeteria(c)
	//computeAverageForMealsInCafeterias(c)
	computeAverageCafeteriaTags(c)
	return nil
}

func computeAverageCafeteriaTags(c *CronService) {

	var res []averageCafeteriaTagsTest
	err := c.db.Raw("SELECT cr.*, crt.*" +
		" FROM cafeteria_rating cr" +
		" JOIN cafeteria_rating_tags crt ON cr.id = crt.parentRating").Scan(&res)

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
