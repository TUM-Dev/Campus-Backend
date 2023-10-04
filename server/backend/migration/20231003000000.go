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

//go:embed static_data
var staticData embed.FS

// migrate20231003000000
// migrates the static data for the canteen rating system and adds the necessary cronjob entries
func (m TumDBMigrator) migrate20231003000000() *gormigrate.Migration {
	return &gormigrate.Migration{
		ID: "20231003000000",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Delete(&model.DishNameTagOption{}, &model.CafeteriaRatingTagOption{}, &model.DishNameTagOptionIncluded{}, &model.DishNameTagOptionExcluded{}).Error; err != nil {
				return err
			}
			setTagTable("static_data/dishRatingTags.json", tx, backend.DISH)
			setTagTable("static_data/cafeteriaRatingTags.json", tx, backend.CAFETERIA)
			setNameTagOptions(tx)
			errRating := addEntriesForCronJob(tx, "averageRatingComputation", 300)
			if errRating != nil {
				return errRating
			}
			errDish := addEntriesForCronJob(tx, "dishNameDownload", 302400)
			if errDish != nil {
				return errDish
			}
			return nil
		},
	}
}

func addEntriesForCronJob(tx *gorm.DB, cronName string, interval int32) error {
	return tx.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Create(&model.Crontab{
		Interval: interval,
		Type:     null.StringFrom(cronName),
		LastRun:  0,
	}).Error
}

/*
Updates the list of dishtags.
If a tag with the exact german and english name does not exist yet, it will be created.
Old tags won't be removed to prevent problems with foreign keys.
*/
func setNameTagOptions(db *gorm.DB) {
	tagsNames := generateNameTagListFromFile("static_data/dishNameTags.json")
	for _, v := range tagsNames.MultiLanguageNameTags {
		var parentId int64

		parent := model.DishNameTagOption{
			DE: v.TagNameGerman,
			EN: v.TagNameEnglish,
		}

		if err := db.Model(&model.DishNameTagOption{}).Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Create(&parent).Error; err != nil {
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
		err := db.Model(&model.DishNameTagOptionExcluded{}).
			Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Create(&model.DishNameTagOptionExcluded{
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

		err := db.Model(&model.DishNameTagOptionIncluded{}).
			Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
			Create(&model.DishNameTagOptionIncluded{
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
	insertModel := getTagModel(tagType, db)
	for _, v := range tagsDish.MultiLanguageTags {
		fields := log.Fields{"de": v.TagNameGerman, "en": v.TagNameEnglish}
		var createError error
		if tagType == backend.CAFETERIA {
			element := model.CafeteriaRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}
			createError = insertModel.Create(&element).Error
		} else {
			element := model.DishRatingTagOption{
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

func getTagModel(tagType backend.ModelType, db *gorm.DB) *gorm.DB {
	if tagType == backend.DISH {
		return db.Model(&model.DishRatingTagOption{})
	} else {
		return db.Model(&model.CafeteriaRatingTagOption{})
	}
}

func generateNameTagListFromFile(path string) multiLanguageNameTags {
	file, err := staticData.ReadFile(path)
	if err != nil {
		log.WithError(err).Error("Error including json.")
	}

	var tags multiLanguageNameTags
	errjson := json.Unmarshal(file, &tags)
	if errjson != nil {
		log.WithError(errjson).Error("Error parsing nameTagList to json.")
	}
	return tags
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
