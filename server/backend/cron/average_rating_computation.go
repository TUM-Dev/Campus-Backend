package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

// averageRatingComputation
// This cronjob precomputes average ratings of all cafeteria ratings, dish ratings and all three types of tags.
// They are grouped (e.g. All Ratings for "Mensa_garching") and the computed values will then be stored in a table with the suffix "_result"
func (c *CronService) averageRatingComputation() error {
	computeAverageCafeteriaTags(c)
	computeAverageForDishesInCafeteriasTags(c)
	computeAverageNameTags(c)

	return nil
}

func computeAverageNameTags(c *CronService) {
	var results []model.DishNameTagAverage
	err := c.db.Raw("SELECT mr.cafeteriaID as cafeteriaID, mnt.tagnameID as tagID, AVG(mnt.points) as average, MAX(mnt.points) as max, MIN(mnt.points) as min, STD(mnt.points) as std" +
		" FROM dish_rating mr" +
		" JOIN dish_name_tag mnt ON mr.dishRating = mnt.correspondingRating" +
		" GROUP BY mr.cafeteriaID, mnt.tagnameID").Scan(&results).Error

	if err != nil {
		log.WithError(err).Error("while precomputing average name tags.")
	} else if len(results) > 0 {
		errDelete := c.db.Where("1=1").Delete(&model.DishNameTagAverage{}).Error // Does not work with "true"
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}
		err := c.db.Model(&model.DishNameTagAverage{}).Create(&results).Error
		if err != nil {
			log.WithError(err).Error("while creating a new average name tag rating in the database.")
		}
	}
}

func computeAverageForDishesInCafeteriasTags(c *CronService) {
	var results []model.DishRatingTagAverage //todo namen im select anpassen
	err := c.db.Raw("SELECT mr.dishID as dishID, mr.cafeteriaID as cafeteriaID, mrt.tagID as tagID, AVG(mrt.points) as average, MAX(mrt.points) as max, MIN(mrt.points) as min, STD(mrt.points) as std" +
		" FROM dish_rating mr" +
		" JOIN dish_rating_tag mrt ON mr.dishRating = mrt.parentRating" +
		" GROUP BY mr.cafeteriaID, mrt.tagID, mr.dishID").Scan(&results).Error

	if err != nil {
		log.WithError(err).Error("while precomputing average dish tags.")
	} else if len(results) > 0 {
		errDelete := c.db.Where("1=1").Delete(&model.DishRatingTagAverage{}).Error
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}

		err := c.db.Model(&model.DishRatingTagAverage{}).Create(&results).Error
		if err != nil {
			log.WithError(err).Error("while creating a new average dish tag rating in the database.")
		}

	}
}

func computeAverageCafeteriaTags(c *CronService) {
	var results []model.CafeteriaRatingTagsAverage
	err := c.db.Raw("SELECT cr.cafeteriaID as cafeteriaID, crt.tagID as tagID, AVG(crt.points) as average, MAX(crt.points) as max, MIN(crt.points) as min, STD(crt.points) as std" +
		" FROM cafeteria_rating cr" +
		" JOIN cafeteria_rating_tag crt ON cr.cafeteriaRating = crt.correspondingRating" +
		" GROUP BY cr.cafeteriaID, crt.tagID").Scan(&results).Error

	if err != nil {
		log.WithError(err).Error("while precomputing average cafeteria tags.")
	} else if len(results) > 0 {
		errDelete := c.db.Where("1=1").Delete(&model.CafeteriaRatingTagsAverage{}).Error
		if errDelete != nil {
			log.WithError(errDelete).Error("Error while deleting old averages in the table.")
		}

		err := c.db.Model(&model.CafeteriaRatingTagsAverage{}).Create(&results).Error
		if err != nil {
			log.WithError(err).Error("while creating a new average cafeteria tag rating in the database.")
		}
	}
}
