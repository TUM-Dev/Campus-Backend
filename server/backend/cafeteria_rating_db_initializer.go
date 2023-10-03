package backend

import (
	"embed"
	"encoding/json"

	"github.com/TUM-Dev/Campus-Backend/server/model"
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

/*
Writes all available tags from the json file into tables in order to make them easier to use.
Will be executed once while the server is started.
*/
func initTagRatingOptions(db *gorm.DB) {
	updateTagTable("static_data/dishRatingTags.json", db, DISH)
	updateTagTable("static_data/cafeteriaRatingTags.json", db, CAFETERIA)
	updateNameTagOptions(db)
	addEntriesForCronJob(db, "averageRatingComputation", 300)
	addEntriesForCronJob(db, "dishNameDownload", 302400) //run twice every week
}

func addEntriesForCronJob(db *gorm.DB, cronName string, interval int32) {
	var count int64
	err := db.Model(&model.Crontab{}).
		Where("type LIKE ?", cronName).
		Count(&count).
		Error

	if err != nil {
		log.WithError(err).WithField("cronName", cronName).Error("Error while checking if cronjob already exists in database")
	} else if count == 0 {
		errCreate := db.Model(&model.Crontab{}).
			Create(&model.Crontab{
				Interval: interval,
				Type:     null.StringFrom(cronName),
				LastRun:  0,
			}).Error
		if errCreate != nil {
			log.WithError(errCreate).Error("Error while creating cronjob with name: ", cronName)
		}
	}
}

/*
Updates the list of dishtags.
If a tag with the exact german and english name does not exist yet, it will be created.
Old tags won't be removed to prevent problems with foreign keys.
*/
func updateNameTagOptions(db *gorm.DB) {
	tagsNames := generateNameTagListFromFile("static_data/dishNameTags.json")
	for _, v := range tagsNames.MultiLanguageNameTags {
		var parentId int64
		res := db.Model(&model.DishNameTagOption{}).
			Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
			Select("DishNameTagOption").
			Scan(&parentId)
		fields := log.Fields{"en": v.TagNameEnglish, "de": v.TagNameGerman}
		if res.Error != nil {
			log.WithError(res.Error).WithFields(fields).Error("Unable to load tag")
		}
		if res.RowsAffected == 0 || res.Error != nil {
			parent := model.DishNameTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}

			if err := db.Model(&model.DishNameTagOption{}).Create(&parent).Error; err != nil {
				log.WithError(err).WithFields(fields).Error("Error while creating tag")
			}
			parentId = parent.DishNameTagOption
		}

		addCanBeIncluded(parentId, db, v)
		addNotIncluded(parentId, db, v)
	}
}

func addNotIncluded(parentId int64, db *gorm.DB, v nameTag) {
	var count int64
	for _, expression := range v.NotIncluded {
		fields := log.Fields{"expression": expression, "parentId": parentId}
		err := db.Model(&model.DishNameTagOptionExcluded{}).
			Where("expression LIKE ? AND NameTagID = ?", expression, parentId).
			Select("DishNameTagOptionExcluded").
			Count(&count).Error
		if err != nil {
			log.WithError(err).WithFields(fields).Error("Unable to load can be excluded tag")
		} else {
			if count == 0 {
				err := db.Model(&model.DishNameTagOptionExcluded{}).
					Create(&model.DishNameTagOptionExcluded{
						Expression: expression,
						NameTagID:  parentId}).Error
				if err != nil {
					log.WithError(err).WithFields(fields).Error("Unable to create new can be excluded tag")
				}
			}
		}
	}
}

func addCanBeIncluded(parentId int64, db *gorm.DB, v nameTag) {
	var count int64
	for _, expression := range v.CanBeIncluded {
		fields := log.Fields{"expression": expression, "parentId": parentId}
		err := db.Model(&model.DishNameTagOptionIncluded{}).
			Where("expression LIKE ? AND NameTagID = ?", expression, parentId).
			Select("DishNameTagOptionIncluded").
			Count(&count).Error
		if err != nil {
			log.WithError(err).WithFields(fields).Error("Unable to load can be included tag")
		} else {
			if count == 0 {
				err := db.Model(&model.DishNameTagOptionIncluded{}).
					Create(&model.DishNameTagOptionIncluded{
						Expression: expression,
						NameTagID:  parentId,
					}).Error
				if err != nil {
					log.WithError(err).WithFields(fields).Error("Unable to create new can be excluded tag")
				}
			}
		}
	}
}

/*
Reads the json file at the given path and checks whether the values have already been inserted into the corresponding table.
If an entry with the same German and English name exists, the entry won't be added.
The TagType is used to identify the corresponding model
*/
func updateTagTable(path string, db *gorm.DB, tagType modelType) {
	tagsDish := generateRatingTagListFromFile(path)
	insertModel := getTagModel(tagType, db)
	for _, v := range tagsDish.MultiLanguageTags {
		var count int64
		fields := log.Fields{"de": v.TagNameGerman, "en": v.TagNameEnglish}
		if tagType == CAFETERIA {
			countError := db.Model(&model.CafeteriaRatingTagOption{}).
				Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
				Select("cafeteriaRatingTagOption").Count(&count).Error
			if countError != nil {
				log.WithError(countError).WithFields(fields).Error("Unable to find cafeteria rating tag")
			}
		} else {
			countError := db.Model(&model.DishRatingTagOption{}).
				Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
				Select("dishRatingTagOption").Count(&count).Error
			if countError != nil {
				log.WithError(countError).WithFields(fields).Error("Unable to find dish rating tag")
			}
		}

		if count == 0 {
			var createError error
			if tagType == CAFETERIA {
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
				createError = insertModel.Create(&element).Error
			}

			if createError != nil {
				log.WithError(createError).WithFields(fields).Error("Unable to create new can be excluded tag")
			} else {
				log.WithFields(fields).Info("New Entry successfully created")
			}
		}
	}
}

func getTagModel(tagType modelType, db *gorm.DB) *gorm.DB {
	if tagType == DISH {
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
