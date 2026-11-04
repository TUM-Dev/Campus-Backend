package migration

import (
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"
)

// CafeteriaRatingAverage stores all precomputed values for the cafeteria ratings
type CafeteriaRatingAverage struct {
	CafeteriaRatingAverage int64   `gorm:"primary_key;autoIncrement;column:cafeteriaRatingAverage;type:int;not null;"`
	CafeteriaID            int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	Average                float64 `gorm:"column:average;type:double;not null;"`
	Min                    int32   `gorm:"column:min;type:int;not null;"`
	Max                    int32   `gorm:"column:max;type:int;not null;"`
	Std                    float64 `gorm:"column:std;type:double;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingAverage) TableName() string {
	return "cafeteria_rating_average"
}

// DishRatingAverage stores all precomputed values for the cafeteria ratings
type DishRatingAverage struct {
	DishRatingAverage int64   `gorm:"primary_key;autoIncrement;column:dishRatingAverage;type:int;not null;"`
	CafeteriaID       int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	DishID            int64   `gorm:"column:dishID;foreignKey:dish;type:int;not null;"`
	Average           float64 `gorm:"column:average;type:double;not null;"`
	Min               int32   `gorm:"column:min;type:int;not null;"`
	Max               int32   `gorm:"column:max;type:int;not null;"`
	Std               float64 `gorm:"column:std;type:double;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingAverage) TableName() string {
	return "dish_rating_average"
}

// DishRatingTagAverage stores all precomputed values for the cafeteria ratings
type DishRatingTagAverage struct {
	DishRatingTagsAverage int64   `gorm:"primary_key;autoIncrement;column:dishRatingTagsAverage;type:int;not null;"`
	CafeteriaID           int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;"`
	TagID                 int64   `gorm:"column:tagID;foreignKey:tagID;type:int;not null;"`
	DishID                int64   `gorm:"column:dishID;foreignKey:dishID;type:int;not null;"`
	Average               float64 `gorm:"column:average;type:double;not null;"`
	Min                   int8    `gorm:"column:min;type:int;not null;"`
	Max                   int8    `gorm:"column:max;type:int;not null;"`
	Std                   float64 `gorm:"column:std;type:double;not null;"`
}

// TableName sets the insert table name for this struct type
func (n *DishRatingTagAverage) TableName() string {
	return "dish_rating_tag_average"
}

// CafeteriaRatingTagsAverage stores all precomputed values for the cafeteria ratings
type CafeteriaRatingTagsAverage struct {
	CafeteriaRatingTagsAverage int64   `gorm:"primary_key;autoIncrement;column:cafeteriaRatingTagsAverage;type:int;not null;" json:"canteenRatingTagsAverage"`
	CafeteriaID                int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"canteenID"`
	TagID                      int64   `gorm:"column:tagID;foreignKey:cafeteriaRatingTagOption;type:int;not null;" json:"tagID"`
	Average                    float64 `gorm:"column:average;type:double;not null;" json:"average"`
	Min                        int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                        int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                        float64 `gorm:"column:std;type:double;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *CafeteriaRatingTagsAverage) TableName() string {
	return "cafeteria_rating_tag_average"
}

// DishNameTagAverage stores all precomputed values for the DishName ratings
type DishNameTagAverage struct {
	DishNameTagAverage int64   `gorm:"primary_key;autoIncrement;column:dishNameTagAverage;type:int;not null;" json:"dishNameTagAverage" `
	CafeteriaID        int64   `gorm:"column:cafeteriaID;foreignKey:cafeteria;type:int;not null;" json:"cafeteriaID"`
	TagID              int64   `gorm:"column:tagID;foreignKey:DishNameTagOption;type:int;not null;" json:"tagID"`
	Average            float64 `gorm:"column:average;type:double;not null;" json:"average" `
	Min                int8    `gorm:"column:min;type:int;not null;" json:"min"`
	Max                int8    `gorm:"column:max;type:int;not null;" json:"max"`
	Std                float64 `gorm:"column:std;type:double;not null;" json:"std"`
}

// TableName sets the insert table name for this struct type
func (n *DishNameTagAverage) TableName() string {
	return "dish_name_tag_average"
}

// migrate20231015000000
// migrates the static data for the canteen rating system and adds the necessary cronjob entries
func migrate20231015000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231015000000",
		Migrate: func(tx *gorm.DB) error {
			// cronjob
			if err := tx.Delete(&model.Crontab{}, "type = 'averageRatingComputation'").Error; err != nil {
				return err
			}
			if err := SafeEnumRemove(tx, &model.Crontab{}, "type", "averageRatingComputation"); err != nil {
				return err
			}
			// tables
			tables := []string{"cafeteria_rating_average", "dish_rating_average", "dish_rating_tag_average", "cafeteria_rating_tag_average", "dish_name_tag_average"}
			for _, table := range tables {
				if err := tx.Migrator().DropTable(table); err != nil {
					return err
				}
			}
			// views
			if err := tx.Exec(`CREATE VIEW cafeteria_rating_statistics AS
SELECT cafeteriaID, Avg(points) AS average, MIN(points) AS min, Max(points) AS max, STD(points) AS std
FROM cafeteria_rating
GROUP BY cafeteriaID
ORDER BY COUNT(*) DESC, average DESC`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE VIEW dish_rating_statistics AS
select d.cafeteriaID AS cafeteriaID,
       r.dishID      AS dishID,
       avg(r.points) AS average,
       max(r.points) AS max,
       min(r.points) AS min,
       std(r.points) AS std
from dish d,
     dish_rating r
WHERE r.dishID=d.dish
group by d.cafeteriaID, r.dishID
order by count(0) desc, avg(r.points) desc`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE VIEW dish_rating_tag_statistics AS
select d.cafeteriaID  AS cafeteriaID,
       mrt.tagID       AS tagID,
       mr.dishID       AS dishID,
       avg(mrt.points) AS average,
       max(mrt.points) AS max,
       min(mrt.points) AS min,
       std(mrt.points) AS std
from dish d,
    dish_rating mr,
    dish_rating_tag mrt
WHERE mr.dishRating = mrt.parentRating AND d.dish=mr.dishID
group by d.cafeteriaID, mrt.tagID, mr.dishID;`).Error; err != nil {
				return err
			}
			if err := tx.Exec(`CREATE VIEW cafeteria_rating_tag_statistics AS
SELECT cr.cafeteriaID as cafeteriaID, crt.tagID as tagID, AVG(crt.points) as average, MAX(crt.points) as max, MIN(crt.points) as min, STD(crt.points) as std
FROM cafeteria_rating cr
JOIN cafeteria_rating_tag crt ON cr.cafeteriaRating = crt.correspondingRating
GROUP BY cr.cafeteriaID, crt.tagID`).Error; err != nil {
				return err
			}
			return tx.Exec(`CREATE VIEW dish_name_tag_statistics AS
select d.cafeteriaID   AS cafeteriaID,
       mnt.tagNameID   AS tagID,
       avg(mnt.points) AS average,
       max(mnt.points) AS max,
       min(mnt.points) AS min,
       std(mnt.points) AS std
from dish d,
     dish_rating mr,
     dish_name_tag mnt
WHERE d.dish = mr.dishID
  AND mr.dishRating = mnt.correspondingRating
group by d.cafeteriaID, mnt.tagNameID;
`).Error
		},
		Rollback: func(tx *gorm.DB) error {
			// views
			createdViews := []string{"cafeteria_rating_statistics", "dish_rating_statistics", "dish_rating_tag_statistics", "cafeteria_rating_tag_statistics", "dish_name_tag_statistics"}
			for _, view := range createdViews {
				if err := tx.Exec("DROP VIEW IF EXISTS " + view).Error; err != nil {
					return err
				}
			}
			// tables
			if err := tx.AutoMigrate(&CafeteriaRatingAverage{}, &DishRatingAverage{}, &DishRatingTagAverage{}, &CafeteriaRatingTagsAverage{}, &DishNameTagAverage{}); err != nil {
				return err
			}
			// cronjob
			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "averageRatingComputation"); err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 300,
				Type:     null.StringFrom("averageRatingComputation"),
			}).Error
		},
	}
}
