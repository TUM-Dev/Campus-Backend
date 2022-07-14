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
	err := c.db.Raw("SELECT  mr.cafeteriaID as cafeteriaID, mnt.tagnameID as tagID, AVG(mnt.rating) as average, MAX(mnt.rating) as max, MIN(mnt.rating) as min, STD(mnt.rating) as std" +
		" FROM meal_rating mr" +
		" JOIN meal_name_tags mnt ON mr.id = mnt.parentRating" +
		" GROUP BY mr.cafeteriaID, mnt.tagnameID").Scan(&results).Error

	if err != nil {
		log.Println(err)
	} else {
		c.db.Where("1=1").Delete(&cafeteria_rating_models.MealNameTagsAverage{})
		for _, v := range results {
			cafeteria := cafeteria_rating_models.MealNameTagsAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				TagID:       v.TagID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&cafeteria_rating_models.MealNameTagsAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.Println(errCreate)
			}
		}
	}
}

func computeAverageForMealsInCafeteriasTags(c *CronService) {
	var results []averageMealTags
	err := c.db.Raw("SELECT mr.mealID as mealID, mr.cafeteriaID as cafeteriaID, mrt.tagID as tagID, AVG(mrt.rating) as average, MAX(mrt.rating) as max, MIN(mrt.rating) as min, STD(mrt.rating) as std" +
		" FROM meal_rating mr" +
		" JOIN meal_rating_tags mrt ON mr.id = mrt.parentRating" +
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
	err := c.db.Raw("SELECT cr.cafeteriaID as cafeteriaID, crt.tagID as tagID, AVG(crt.rating) as average, MAX(crt.rating) as max, MIN(crt.rating) as min, STD(crt.rating) as std" +
		" FROM cafeteria_rating cr" +
		" JOIN cafeteria_rating_tags crt ON cr.id = crt.parentRating" +
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
		Select("cafeteriaID, mealID, AVG(rating) as average, MAX(rating) as max, MIN(rating) as min, STD(rating) as std").
		Group("cafeteriaID,mealID").Scan(&results)

	if res.Error != nil {
		log.Println("Error in query")
		log.Println(res.Error)
	} else {
		c.db.Where("1=1").Delete(&cafeteria_rating_models.MealRatingsAverage{})
		for _, v := range results {
			cafeteria := cafeteria_rating_models.MealRatingsAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				MealID:      v.MealID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&cafeteria_rating_models.MealRatingsAverage{}).Create(&cafeteria)
			if errCreate.Error != nil {
				log.Println(errCreate.Error)
			}
		}
	}
}

func computeAverageForCafeteria(c *CronService) {
	var results []averageRatingForCafeteria
	res := c.db.Model(&cafeteria_rating_models.CafeteriaRating{}).
		Select("cafeteriaID, AVG(rating) as average, MAX(rating) as max, MIN(rating) as min, STD(rating) as std").
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
