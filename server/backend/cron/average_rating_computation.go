package cron

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
)

// averageRatingComputation
// This cronjob precomputes average ratings of all cafeteria ratings, dish ratings and all three types of tags.
// They are grouped (e.g. All Ratings for "Mensa_garching") and the computed values will then be stored in a table with the suffix "_result"
func (c *CronService) averageRatingComputation() error {
	computeAverageNameTags(c)

	return nil
}

func computeAverageNameTags(c *CronService) {
	var results []model.DishNameTagAverage
	err := c.db.Raw(`SELECT mr.cafeteriaID as cafeteriaID, mnt.tagnameID as tagID, AVG(mnt.points) as average, MAX(mnt.points) as max, MIN(mnt.points) as min, STD(mnt.points) as std
FROM dish_rating mr
JOIN dish_name_tag mnt ON mr.dishRating = mnt.correspondingRating
GROUP BY mr.cafeteriaID, mnt.tagnameID`).Scan(&results).Error

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
