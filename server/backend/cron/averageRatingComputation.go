package cron

import log "github.com/sirupsen/logrus"

//fileDownloadCron Downloads all files that are not marked as finished in the database.
func (c *CronService) averageRatingComputation() error {

	res := c.db.Table("cafeteria_rating").Raw("SELECT cafeteria,id,name, AVG(rating) AS avgrating \nFROM cafeteria_rating \nGROUP BY cafeteria;") //.Group("cafeterianame")
	log.Println("Number of affected rows", res.RowsAffected)
	//todo in welcher form bekomme ich ergebnisse
	/*
	   raw sql
	      DECLARE
	          @Avg int,
	          @StDev int

	      SELECT Meal, Canteen, @Avg = AVG(Sales), @StDev = STDEV(Sales)
	      FROM tbl_sales
	      WHERE ...



	   Groupby auf canteen -> später auch mit den dishes - für jede gruppe erstmal den average zurückgeben
	*/

	//todo save: average per group, standard deviation, how much elements per rating group(how many have rated with 5 stars?)
	return nil
}
