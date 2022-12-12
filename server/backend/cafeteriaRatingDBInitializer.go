package backend

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

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

/*
Writes all available tags from the json file into tables in order to make them easier to use.
Will be executed once while the server is started.
*/
func initTagRatingOptions(db *gorm.DB) {
	updateTagTable("backend/static_data/dishRatingTags.json", db, DISH)
	updateTagTable("backend/static_data/cafeteriaRatingTags.json", db, CAFETERIA)
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
		log.WithError(err).Error("Error while checking if cronjob with name {} already exists in database", cronName)
	} else if count == 0 {
		errCreate := db.Model(&model.Crontab{}).
			Create(&model.Crontab{
				Interval: interval,
				Type:     null.String{NullString: sql.NullString{String: cronName, Valid: true}},
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
	absPathDishNames, _ := filepath.Abs("backend/static_data/dishNameTags.json")
	tagsNames := generateNameTagListFromFile(absPathDishNames)
	for _, v := range tagsNames.MultiLanguageNameTags {
		var parentId int32
		res := db.Model(&model.DishNameTagOption{}).
			Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
			Select("DishNameTagOption").
			Scan(&parentId)
		if res.Error != nil {
			log.WithError(res.Error).Error("Unable to load tag with En {} and De {} ", v.TagNameEnglish, v.TagNameGerman)
		}
		if res.RowsAffected == 0 || res.Error != nil {
			parent := model.DishRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}

			errCreate := db.Model(&model.DishNameTagOption{}).
				Create(&parent).Error
			if errCreate != nil {
				log.WithError(errCreate).Error("Error while creating tag {}, {}.", v.TagNameGerman, v.TagNameEnglish)
			}
			parentId = parent.DishRatingTagOption
		}

		addCanBeIncluded(parentId, db, v)
		addNotIncluded(parentId, db, v)

	}
}

func addNotIncluded(parentId int32, db *gorm.DB, v nameTag) {
	var count int64
	for _, u := range v.NotIncluded {
		errorLoadingIncluded := db.Model(&model.DishNameTagOptionExcluded{}).
			Where("expression LIKE ? AND NameTagID = ?", u, parentId).
			Select("DishNameTagOptionExcluded").
			Count(&count).Error
		if errorLoadingIncluded != nil {
			log.WithError(errorLoadingIncluded).Error("Unable to load can be excluded tag with expression {} and parentId {} ", u, parentId)
		} else {
			if count == 0 {
				createError := db.Model(&model.DishNameTagOptionExcluded{}).
					Create(&model.DishNameTagOptionExcluded{
						Expression: u,
						NameTagID:  parentId}).Error
				if createError != nil {
					log.WithError(errorLoadingIncluded).Error("Unable to create new can be excluded tag with expression {} and parentId {} ", u, parentId)
				}
			}
		}
	}
}

func addCanBeIncluded(parentId int32, db *gorm.DB, v nameTag) {
	var count int64
	for _, u := range v.CanBeIncluded {
		errorLoadingIncluded := db.Model(&model.DishNameTagOptionIncluded{}).
			Where("expression LIKE ? AND NameTagID = ?", u, parentId).
			Select("DishNameTagOptionIncluded").
			Count(&count).Error
		if errorLoadingIncluded != nil {
			log.WithError(errorLoadingIncluded).Error("Unable to load can be included tag with expression {} and parentId {} ", u, parentId)
		} else {
			if count == 0 {
				createError := db.Model(&model.DishNameTagOptionIncluded{}).
					Create(&model.DishNameTagOptionIncluded{
						Expression: u,
						NameTagID:  parentId,
					}).Error
				if createError != nil {
					log.WithError(errorLoadingIncluded).Error("Unable to create new can be excluded tag with expression {} and parentId {} ", u, parentId)
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
	absPathDish, _ := filepath.Abs(path)
	tagsDish := generateRatingTagListFromFile(absPathDish)
	insertModel := getTagModel(tagType, db)
	for _, v := range tagsDish.MultiLanguageTags {
		var count int64

		if tagType == CAFETERIA {
			countError := db.Model(&model.CafeteriaRatingTagOption{}).
				Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
				Select("cafeteriaRatingTagOption").Count(&count).Error
			if countError != nil {
				log.WithError(countError).Error("Unable to find cafeteria rating tag with En {} and De {} ", v.TagNameGerman, v.TagNameEnglish)
			}
		} else {
			countError := db.Model(&model.DishRatingTagOption{}).
				Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
				Select("dishRatingTagOption").Count(&count).Error
			if countError != nil {
				log.WithError(countError).Error("Unable to find dish rating tag with En {} and De {} ", v.TagNameGerman, v.TagNameEnglish)
			}
		}

		if count == 0 {
			element := model.DishRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish,
			}
			createError := insertModel.Create(&element).Error
			if createError != nil {
				log.WithError(createError).Error("Unable to create new can be excluded tag with En {} and De {} ", v.TagNameGerman, v.TagNameEnglish)
			} else {
				log.Info("New Entry with En ", v.TagNameEnglish, " and De ", v.TagNameGerman, " successfully created.")
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
	file := readFromFile(path)

	var tags multiLanguageNameTags
	errjson := json.NewDecoder(file).Decode(&tags)
	if errjson != nil {
		log.WithError(errjson).Error("Error while reading the file.")
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.WithError(err).Error("Error in parsing json.")
		}
	}(file)
	return tags
}

func generateRatingTagListFromFile(path string) multiLanguageTags {
	file := readFromFile(path)
	var tags multiLanguageTags
	errjson := json.NewDecoder(file).Decode(&tags)
	if errjson != nil {
		log.WithError(errjson).Error("Error while reading or parsing the file.")
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.WithError(err).Error("Error in parsing json.")
		}
	}(file)
	return tags
}

func readFromFile(path string) *os.File {
	jsonFile, err := os.Open(path)

	if err != nil {
		log.WithError(err).Error("Unable to open file with path: {}", path)
	}

	return jsonFile
}
