package migration

import (
	"embed"
	"encoding/json"

	"gorm.io/gorm/logger"

	"github.com/TUM-Dev/Campus-Backend/server/backend"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/guregu/null"
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

// State of the models at this migration
type CafeteriaRatingTagOption struct {
	CafeteriaRatingsTagOption int64  `gorm:"primary_key;AUTO_INCREMENT;column:cafeteriaRatingTagOption;type:int;not null;" json:"canteenRatingTagOption"`
	DE                        string `gorm:"column:DE;text;default:de;not null;" json:"DE"`
	EN                        string `gorm:"column:EN;text;default:en;not null;" json:"EN"`
}

type DishRatingTagOption struct {
	DishRatingTagOption int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishRatingTagOption;type:int;not null;" json:"dishRatingTagOption"`
	DE                  string `gorm:"column:DE;type:text;default:de;not null;" json:"DE"`
	EN                  string `gorm:"column:EN;type:text;default:en;not null;" json:"EN"`
}

type DishNameTagOption struct {
	DishNameTagOption int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOption;type:int;not null;" json:"dishNameTagOption"`
	DE                string `gorm:"column:DE;type:text;not null;" json:"DE"`
	EN                string `gorm:"column:EN;type:text;not null;" json:"EN"`
}

type DishNameTagOptionExcluded struct {
	DishNameTagOptionExcluded int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOptionExcluded;type:int;not null;" json:"dishNameTagOptionExcluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

type DishNameTagOptionIncluded struct {
	DishNameTagOptionIncluded int64  `gorm:"primary_key;AUTO_INCREMENT;column:dishNameTagOptionIncluded;type:int;not null;" json:"dishNameTagOptionIncluded"`
	NameTagID                 int64  `gorm:"foreignKey:dishNameTagOption;column:nameTagID;type:int;not null;" json:"nameTagID"`
	Expression                string `gorm:"column:expression;type:text;" json:"expression"`
}

//go:embed static_data
var staticData embed.FS

// migrate20231003000000
// migrates the static data for the canteen rating system and adds the necessary cronjob entries
func (m TumDBMigrator) migrate20231003000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231003000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Delete(&DishNameTagOption{}, &CafeteriaRatingTagOption{}, &DishNameTagOptionIncluded{}, &DishNameTagOptionExcluded{}).Error; err != nil {
				return err
			}
			// todo re tables actually created?
			setTagTable("static_data/dishRatingTags.json", tx, backend.DISH)
			setTagTable("static_data/cafeteriaRatingTags.json", tx, backend.CAFETERIA)
			setNameTagOptions(tx)
			err := tx.Create(&model.Crontab{
				Interval: 300,
				Type:     null.StringFrom("averageRatingComputation"),
			}).Error
			if err != nil {
				return err
			}
			err = tx.Create(&model.Crontab{
				Interval: 302400,
				Type:     null.StringFrom("dishNameDownload"),
			}).Error
			return err
		},
		Rollback: func(tx *gorm.DB) error {
			err := tx.Delete(&DishNameTagOption{}, &CafeteriaRatingTagOption{}, &DishNameTagOptionIncluded{}, &DishNameTagOptionExcluded{}).Error
			return err
		},
	}
}

/*
Updates the list of dishtags.
If a tag with the exact german and english name does not exist yet, it will be created.
Old tags won't be removed to prevent problems with foreign keys.
*/
func setNameTagOptions(db *gorm.DB) {
	file, err := staticData.ReadFile("static_data/dishNameTags.json")
	if err != nil {
		log.WithError(err).Error("Error including json.")
	}

	var tagsNames multiLanguageNameTags
	errjson := json.Unmarshal(file, &tagsNames)
	if errjson != nil {
		log.WithError(errjson).Error("Error parsing nameTagList to json.")
	}
	for _, v := range tagsNames.MultiLanguageNameTags {
		var parentId int64

		parent := DishNameTagOption{
			DE: v.TagNameGerman,
			EN: v.TagNameEnglish,
		}

		if err := db.Model(&DishNameTagOption{}).Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Create(&parent).Error; err != nil {
			fields := log.Fields{"en": v.TagNameEnglish, "de": v.TagNameGerman}
			log.WithError(err).WithFields(fields).Error("Error while creating tag")
		}
		parentId = parent.DishNameTagOption

		addCanBeIncluded(parentId, db, v)
		addNotIncluded(parentId, db, v)
	}
}

func addNotIncluded(parentId int64, db *gorm.DB, v nameTag) {
	for _, expression := range v.NotIncluded {
		fields := log.Fields{"expression": expression, "parentId": parentId}
		err := db.Model(&DishNameTagOptionExcluded{}).
			Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Create(&DishNameTagOptionExcluded{
				Expression: expression,
				NameTagID:  parentId}).Error
		if err != nil {
			log.WithError(err).WithFields(fields).Error("Unable to create new can be excluded tag")
		}
	}
}

func addCanBeIncluded(parentId int64, db *gorm.DB, v nameTag) {
	for _, expression := range v.CanBeIncluded {
		fields := log.Fields{"expression": expression, "parentId": parentId}

		err := db.Model(&DishNameTagOptionIncluded{}).
			Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Create(&DishNameTagOptionIncluded{
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
func setTagTable(path string, db *gorm.DB, tagType backend.ModelType) {
	tagsDish := generateRatingTagListFromFile(path)

	var insertModel *gorm.DB
	if tagType == backend.DISH {
		insertModel = db.Model(&DishRatingTagOption{})
	} else {
		insertModel = db.Model(&CafeteriaRatingTagOption{})
	}

	for _, v := range tagsDish.MultiLanguageTags {
		fields := log.Fields{"de": v.TagNameGerman, "en": v.TagNameEnglish}
		var createError error
		if tagType == backend.CAFETERIA {
			element := CafeteriaRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}
			createError = insertModel.Create(&element).Error
		} else {
			element := DishRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}
			createError = insertModel.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Create(&element).Error
		}

		if createError != nil {
			log.WithError(createError).WithFields(fields).Error("Unable to create new can be excluded tag")
		} else {
			log.WithFields(fields).Info("New Entry successfully created")
		}
	}
}

func generateRatingTagListFromFile(path string) multiLanguageTags {
	file, err := staticData.ReadFile(path)
	if err != nil {
		log.WithError(err).Error("Error including json.")
	}

	var tags multiLanguageTags
	errjson := json.Unmarshal(file, &tags)
	if errjson != nil {
		log.WithError(errjson).Error("Error parsing ratingTagList to json.")
	}
	return tags
}
