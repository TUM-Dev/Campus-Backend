package cron

import (
	"github.com/TUM-Dev/Campus-Backend/model"
	log "github.com/sirupsen/logrus"
)

type averageRatingForCafeteria struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

type averageRatingForDishInCafeteria struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	DishID      int32   `gorm:"column:dishID;foreignKey:id;type:int;" json:"dishID"`
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

type averageDishTags struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	TagID       int32   `gorm:"foreignKey:id;column:tagID;type:int" json:"tagID"`
	DishID      int32   `gorm:"column:dishID;foreignKey:id;type:int;" json:"dishID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

type averageDishNameTags struct {
	CafeteriaID int32   `gorm:"column:cafeteriaID;foreignKey:cafeteriaID;type:int;" json:"cafeteriaID"`
	TagID       int32   `gorm:"foreignKey:id;column:tagID;type:int" json:"tagID"`
	Average     float32 `json:"average"`
	Min         int8    `json:"min"`
	Max         int8    `json:"max"`
	Std         float32 `json:"std"`
}

// averageRatingComputation
// This cronjob precomputes average ratings of all cafeteria ratings, dish ratings and all three types of tags.
// They are grouped (e.g. All Ratings for "Mensa_garching") and the computed values will then be stored in a table with the suffix "_result"
func (c *CronService) averageRatingComputation() error {
	computeAverageForCafeteria(c)
	computeAverageForDishesInCafeterias(c)
	computeAverageCafeteriaTags(c)
	computeAverageForDishesInCafeteriasTags(c)
	computeAverageNameTags(c)

	return nil
}

func computeAverageNameTags(c *CronService) {
	var results []averageDishNameTags
	err := c.db.Raw("SELECT  mr.cafeteriaID as cafeteriaID, mnt.tagnameID as tagID, AVG(mnt.points) as average, MAX(mnt.points) as max, MIN(mnt.points) as min, STD(mnt.points) as std" +
		" FROM dish_rating mr" +
		" JOIN dish_name_tag mnt ON mr.dishRating = mnt.correspondingRating" +
		" GROUP BY mr.cafeteriaID, mnt.tagnameID").Scan(&results).Error

	if err != nil {
		log.WithError(err).Error("Error while precomputing average name tags.")
	} else {
		errDelete := c.db.Where("1=1").Delete(&model.DishNameTagAverage{}).Error
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}
		for _, v := range results {
			cafeteria := model.DishNameTagAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				TagID:       v.TagID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&model.DishNameTagAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating a new average name tag rating in the database.")
			}
		}
	}
}

func computeAverageForDishesInCafeteriasTags(c *CronService) {
	var results []averageDishTags
	err := c.db.Raw("SELECT mr.dishID as dishID, mr.cafeteriaID as cafeteriaID, mrt.tagID as tagID, AVG(mrt.points) as average, MAX(mrt.points) as max, MIN(mrt.points) as min, STD(mrt.points) as std" +
		" FROM dish_rating mr" +
		" JOIN dish_rating_tag mrt ON mr.dishRating = mrt.parentRating" +
		" GROUP BY mr.cafeteriaID, mrt.tagID, mr.dishID").Scan(&results).Error

	if err != nil {
		log.WithError(err).Error("Error while precomputing average dish tags.")
	} else {
		errDelete := c.db.Where(true).Delete(&model.DishRatingTagAverage{}).Error
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}

		for _, v := range results {
			cafeteria := model.DishRatingTagAverage{
				CafeteriaID: v.CafeteriaID,
				DishID:      v.DishID,
				Average:     v.Average,
				TagID:       v.TagID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&model.DishRatingTagAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating a new average dish tag rating in the database.")
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
		log.WithError(err).Error("Error while precomputing average cafeteria tags.")
	} else {
		errDelete := c.db.Where(true).Delete(&model.CafeteriaRatingTagsAverage{}).Error
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}
		for _, v := range results {
			cafeteria := model.CafeteriaRatingTagsAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				TagID:       v.TagID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&model.CafeteriaRatingTagsAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating a new average cafeteria tag rating in the database.")
			}
		}
	}
}

func computeAverageForDishesInCafeterias(c *CronService) {
	var results []averageRatingForDishInCafeteria
	err := c.db.Model(&model.DishRating{}).
		Select("cafeteriaID, dishID, AVG(points) as average, MAX(points) as max, MIN(points) as min, STD(points) as std").
		Group("cafeteriaID,dishID").Scan(&results).Error

	if err != nil {
		log.WithError(err).Error("Error while precomputing average dish ratings.")
	} else {
		errDelete := c.db.Where(true).Delete(&model.DishRatingAverage{}).Error
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}
		for _, v := range results {
			cafeteria := model.DishRatingAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				DishID:      v.DishID,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Model(&model.DishRatingAverage{}).Create(&cafeteria).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating a new average dish rating in the database.")
			}
		}
	}
}

func computeAverageForCafeteria(c *CronService) {
	var results []averageRatingForCafeteria
	err := c.db.Model(&model.CafeteriaRating{}).
		Select("cafeteriaID, AVG(points) as average, MAX(points) as max, MIN(points) as min, STD(points) as std").
		Group("cafeteriaID").Find(&results).Error

	if err != nil {
		log.WithError(err).Error("Error while precomputing average cafeteria ratings.")
	} else {
		errDelete := c.db.Where(true).Delete(&model.CafeteriaRatingAverage{}).Error
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}
		for _, v := range results {
			cafeteria := model.CafeteriaRatingAverage{
				CafeteriaID: v.CafeteriaID,
				Average:     v.Average,
				Min:         v.Min,
				Max:         v.Max,
				Std:         v.Std,
			}

			errCreate := c.db.Create(&cafeteria).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating a new average cafeteria rating in the database.")
			}
		}
	}
}
