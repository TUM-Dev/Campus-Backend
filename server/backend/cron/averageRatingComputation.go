package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	"log"
)

type averageRatingForCafeteria struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

type averageRatingForMealInCafeteria struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	MealID      int32   `gorm:"column:mealID;foreignKey:id;type:int;" json:"mealID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

type averageCafeteriaTags struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	TagID       int32   `gorm:"foreignKey:tagRatingID;column:tagID;type:int" json:"tagID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

type averageMealTags struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	TagID       int32   `gorm:"foreignKey:id;column:tagID;type:int" json:"tagID"`
	MealID      int32   `gorm:"column:mealID;foreignKey:id;type:int;" json:"mealID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

type averageMealNameTags struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	TagID       int32   `gorm:"foreignKey:id;column:tagID;type:int" json:"tagID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

/*
This cronjob precomputes average ratings of all cafeteria ratings, meal ratings and all three types of tags.
They are grouped (e.g. All Ratings for "Mensa_garching") and the computed values will then be stored in a table with the suffix "_result"
*/
func (c *CronService) averageRatingComputation() error {
	computeAverageForCafeteria(c)
	computeAverageForMealsInCafeterias(c)
	computeAverageCafeteriaTags(c)
	computeAverageForMealsInCafeteriasTags(c)
	computeAverageNameTags(c)

	return nil
}

func computeAverageNameTags(c *CronService) {
	var results []averageMealNameTags
	err := c.db.Raw("SELECT  mr.cafeteriaID as cafeteriaID, mnt.tagnameID as tagID, AVG(mnt.points) as average, MAX(mnt.points) as max, MIN(mnt.points) as min, STD(mnt.points) as std" +
		" FROM meal_rating mr" +
		" JOIN meal_name_tag mnt ON mr.mealRating = mnt.parentRating" +
		" GROUP BY mr.cafeteriaID, mnt.tagnameID").Scan(&results).Error

	if err != nil {
		log.Println(err)
	} else {
		c.db.Where("1=1").Delete(&cafeteria_rating_models.MealNameTagAverage{})
		for _, v := range results {
			cafeteria := cafeteria_rating_models.MealNameTagAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				TagID:       v.TagID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&cafeteria_rating_models.MealNameTagAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.Println(errCreate)
			}
		}
	}
}

func computeAverageForMealsInCafeteriasTags(c *CronService) {
	var results []averageMealTags
	err := c.db.Raw("SELECT mr.mealID as mealID, mr.cafeteriaID as cafeteriaID, mrt.tagID as tagID, AVG(mrt.points) as average, MAX(mrt.points) as max, MIN(mrt.points) as min, STD(mrt.points) as std" +
		" FROM meal_rating mr" +
		" JOIN meal_rating_tag mrt ON mr.mealRating = mrt.parentRating" +
		" GROUP BY mr.cafeteriaID, mrt.tagID, mr.mealID").Scan(&results).Error

	if err != nil {
		log.Println(err)
	} else {
		c.db.Where("1=1").Delete(&cafeteria_rating_models.MealRatingTagAverage{})

		for _, v := range results {
			cafeteria := cafeteria_rating_models.MealRatingTagAverage{
				CafeteriaID: v.CafeteriaID,
				MealID:      v.MealID,
				Average:     v.Average,
				TagID:       v.TagID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&cafeteria_rating_models.MealRatingTagAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.Println(errCreate)
			}
		}
	}
}

func computeAverageCafeteriaTags(c *CronService) {
	var results []averageCafeteriaTags
	err := c.db.Raw("SELECT cr.cafeteriaID as cafeteriaID, crt.tagID as tagID, AVG(crt.points) as average, MAX(crt.points) as max, MIN(crt.points) as min, STD(crt.points) as std" +
		" FROM cafeteria_rating cr" +
		" JOIN cafeteria_rating_tag crt ON cr.cafeteriaRating = crt.correspondingRating" +
		" GROUP BY cr.cafeteriaID, crt.tagID").Scan(&results).Error

	if err != nil {
		log.Println(err)
	} else {
		c.db.Where("1=1").Delete(&cafeteria_rating_models.CafeteriaRatingTagsAverage{})
		for _, v := range results {
			cafeteria := cafeteria_rating_models.CafeteriaRatingTagsAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				TagID:       v.TagID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&cafeteria_rating_models.CafeteriaRatingTagsAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.Println(errCreate)
			}

		}
	}
}

func computeAverageForMealsInCafeterias(c *CronService) {
	var results []averageRatingForMealInCafeteria
	res := c.db.Model(&cafeteria_rating_models.MealRating{}).
		Select("cafeteriaID, mealID, AVG(points) as average, MAX(points) as max, MIN(points) as min, STD(points) as std").
		Group("cafeteriaID,mealID").Scan(&results)

	if res.Error != nil {
		log.Println("Error in query")
		log.Println(res.Error)
	} else {
		c.db.Where("1=1").Delete(&cafeteria_rating_models.MealRatingAverage{})
		for _, v := range results {
			cafeteria := cafeteria_rating_models.MealRatingAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				MealID:      v.MealID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&cafeteria_rating_models.MealRatingAverage{}).Create(&cafeteria)
			if errCreate.Error != nil {
				log.Println(errCreate.Error)
			}
		}
	}
}

func computeAverageForCafeteria(c *CronService) {
	var results []averageRatingForCafeteria
	res := c.db.Model(&cafeteria_rating_models.CafeteriaRating{}).
		Select("cafeteriaID, AVG(points) as average, MAX(points) as max, MIN(points) as min, STD(points) as std").
		Group("cafeteriaID").Find(&results)

	if res.Error != nil {
		log.Println("Error in query")
		log.Println(res.Error)
	} else {
		c.db.Where("1=1").Delete(&cafeteria_rating_models.CafeteriaRatingAverage{})
		for _, v := range results {
			cafeteria := cafeteria_rating_models.CafeteriaRatingAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Create(&cafeteria)
			if errCreate.Error != nil {
				log.Println(errCreate.Error)
			}

		}
	}
}
