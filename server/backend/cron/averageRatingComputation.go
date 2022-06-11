package cron

//fileDownloadCron Downloads all files that are not marked as finished in the database.
func (c *CronService) averageRatingComputation() error {

	//res := c.db.Table("mensa_rating").Group("meal").Select()

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
