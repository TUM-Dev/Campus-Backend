package migration

import (
	"embed"
	"encoding/json"

	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"

	"gorm.io/gorm/logger"

	"github.com/go-gormigrate/gormigrate/v2"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type multiLanguageTags struct {
	MultiLanguageTags []tag `json:"tags"`
}
type tag struct {
	TagNameEnglish string `json:"tagNameEnglish"`
	TagNameGerman  string `json:"tagNameGerman"`
}

type multiLanguageNameTags struct {
	MultiLanguageNameTags []nameTag `json:"tags"`
}
type nameTag struct {
	TagNameEnglish string   `json:"tagNameEnglish"`
	TagNameGerman  string   `json:"tagNameGerman"`
	NotIncluded    []string `json:"notincluded"`
	CanBeIncluded  []string `json:"canbeincluded"`
}

type ToFillCafeteriaRatingTagOption struct {
	CafeteriaRatingsTagOption int64  `gorm:"primary_key;autoIncrement;column:cafeteriaRatingTagOption;type:int;not null;" json:"canteenRatingTagOption"`
	DE                        string `gorm:"column:DE;text;default:('de');not null;" json:"DE"`
	EN                        string `gorm:"column:EN;text;default:('en');not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *ToFillCafeteriaRatingTagOption) TableName() string {
	return "cafeteria_rating_tag_option"
}

type ToFillDishRatingTagOption struct {
	DishRatingTagOption int64  `gorm:"primary_key;autoIncrement;column:dishRatingTagOption;type:int;not null;" json:"dishRatingTagOption"`
	DE                  string `gorm:"column:DE;type:text;default:('de');not null;" json:"DE"`
	EN                  string `gorm:"column:EN;type:text;default:('en');not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *ToFillDishRatingTagOption) TableName() string {
	return "dish_rating_tag_option"
}

type ToFillDishNameTagOption struct {
	DishNameTagOption int64  `gorm:"primary_key;autoIncrement;column:dishNameTagOption;type:int;not null;" json:"dishNameTagOption"`
	DE                string `gorm:"column:DE;type:text;not null;" json:"DE"`
	EN                string `gorm:"column:EN;type:text;not null;" json:"EN"`
}

// TableName sets the insert table name for this struct type
func (n *ToFillDishNameTagOption) TableName() string {
	return "dish_name_tag_option"
}

type ToFillDishNameTagOptionExcluded struct {
	DishNameTagOptionExcluded int64  `gorm:"primary_key;autoIncrement;column:dishNameTagOptionExcluded;type:int;not null;" json:"dishNameTagOptionExcluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *ToFillDishNameTagOptionExcluded) TableName() string {
	return "dish_name_tag_option_excluded"
}

type ToFillDishNameTagOptionIncluded struct {
	DishNameTagOptionIncluded int64  `gorm:"primary_key;autoIncrement;column:dishNameTagOptionIncluded;type:int;not null;" json:"dishNameTagOptionIncluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

// TableName sets the insert table name for this struct type
func (n *ToFillDishNameTagOptionIncluded) TableName() string {
	return "dish_name_tag_option_included"
}

//go:embed static_data
var staticData embed.FS

// migrate20231003000000
// migrates the static data for the canteen rating system and adds the necessary cronjob entries
func migrate20231003000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231003000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Where("1=1").Delete(&ToFillDishNameTagOption{}, &ToFillCafeteriaRatingTagOption{}, &ToFillDishNameTagOptionIncluded{}, &ToFillDishNameTagOptionExcluded{}).Error; err != nil {
				return err
			}

			setTagTable("static_data/dishRatingTags.json", tx, model.DISH)
			setTagTable("static_data/cafeteriaRatingTags.json", tx, model.CAFETERIA)
			setNameTagOptions(tx)
			if err := SafeEnumAdd(tx, &model.Crontab{}, "type", "averageRatingComputation", "dishNameDownload"); err != nil {
				return err
			}
			if err := tx.Create(&model.Crontab{
				Interval: 300,
				Type:     null.StringFrom("averageRatingComputation"),
			}).Error; err != nil {
				return err
			}
			return tx.Create(&model.Crontab{
				Interval: 302400,
				Type:     null.StringFrom("dishNameDownload"),
			}).Error
		},
		Rollback: func(tx *gorm.DB) error {
			if err := tx.Where("1=1").Delete(&ToFillDishNameTagOption{}, &ToFillCafeteriaRatingTagOption{}, &ToFillDishNameTagOptionIncluded{}, &ToFillDishNameTagOptionExcluded{}).Error; err != nil {
				return err
			}
			if err := tx.Delete(&model.Crontab{}, "type = 'averageRatingComputation'").Error; err != nil {
				return err
			}
			if err := tx.Delete(&model.Crontab{}, "type = 'dishNameDownload'").Error; err != nil {
				return err
			}
			return SafeEnumRemove(tx, &model.Crontab{}, "type", "dishNameDownload", "averageRatingComputation")
		},
	}
}

/*
setNameTagOptions updates the list of dishtags.
If a tag with the exact german and english name does not exist yet, it will be created.
Old tags won't be removed to prevent problems with foreign keys.
*/
func setNameTagOptions(db *gorm.DB) {
	file, err := staticData.ReadFile("static_data/dishNameTags.json")
	if err != nil {
		log.WithError(err).Error("Error including json.")
	}

	var tagsNames multiLanguageNameTags
	if err := json.Unmarshal(file, &tagsNames); err != nil {
		log.WithError(err).Error("Error parsing nameTagList to json.")
	}
	for _, v := range tagsNames.MultiLanguageNameTags {
		parent := ToFillDishNameTagOption{
			DE: v.TagNameGerman,
			EN: v.TagNameEnglish,
		}

		if err := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Create(&parent).Error; err != nil {
			fields := log.Fields{"en": v.TagNameEnglish, "de": v.TagNameGerman}
			log.WithError(err).WithFields(fields).Error("Error while creating tag")
		}

		addCanBeIncluded(parent.DishNameTagOption, db, v)
		addNotIncluded(parent.DishNameTagOption, db, v)
	}
}

func addNotIncluded(parentId int64, db *gorm.DB, v nameTag) {
	for _, expression := range v.NotIncluded {
		fields := log.Fields{"expression": expression, "parentId": parentId}
		err := db.
			Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Create(&ToFillDishNameTagOptionExcluded{
				Expression: expression,
				NameTagID:  parentId}).Error
		if err != nil {
			log.WithError(err).WithFields(fields).Error("Unable to create new DishNameTagOptionExcluded")
		}
	}
}

func addCanBeIncluded(parentId int64, db *gorm.DB, v nameTag) {
	for _, expression := range v.CanBeIncluded {
		fields := log.Fields{"expression": expression, "parentId": parentId}

		err := db.
			Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Create(&ToFillDishNameTagOptionIncluded{
				Expression: expression,
				NameTagID:  parentId,
			}).Error
		if err != nil {
			log.WithError(err).WithFields(fields).Error("Unable to create new can be excluded tag")
		}
	}
}

/*
Reads the json file at the given path and checks whether the values have already been inserted into the corresponding table.
If an entry with the same German and English name exists, the entry won't be added.
The TagType is used to identify the corresponding model
*/
func setTagTable(path string, db *gorm.DB, tagType model.ModelType) {
	tagsDish := generateRatingTagListFromFile(path)

	for _, v := range tagsDish.MultiLanguageTags {
		fields := log.Fields{"de": v.TagNameGerman, "en": v.TagNameEnglish}
		var err error
		if tagType == model.CAFETERIA {
			element := ToFillCafeteriaRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}
			err = db.Create(&element).Error
		} else {
			element := ToFillDishRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}
			err = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Create(&element).Error
		}

		if err != nil {
			log.WithError(err).WithFields(fields).Error("Unable to create new can be excluded tag")
		} else {
			log.WithFields(fields).Debug("New Tag successfully created")
		}
	}
}

func generateRatingTagListFromFile(path string) multiLanguageTags {
	file, err := staticData.ReadFile(path)
	if err != nil {
		log.WithError(err).Error("Error including json.")
	}

	var tags multiLanguageTags
	if err := json.Unmarshal(file, &tags); err != nil {
		log.WithError(err).Error("Error parsing ratingTagList to json.")
	}
	return tags
}
