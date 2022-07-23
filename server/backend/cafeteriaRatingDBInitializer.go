package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/TUM-Dev/Campus-Backend/model"
	"github.com/TUM-Dev/Campus-Backend/model/cafeteria_rating_models"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"path/filepath"
)

type MultiLanguageTags struct {
	MultiLanguageTags []Tag `json:"tags"`
}
type Tag struct {
	TagNameEnglish string `json:"tagNameEnglish"`
	TagNameGerman  string `json:"tagNameGerman"`
}

type MultiLanguageNameTags struct {
	MultiLanguageNameTags []NameTag `json:"tags"`
}
type NameTag struct {
	TagNameEnglish string   `json:"tagNameEnglish"`
	TagNameGerman  string   `json:"tagNameGerman"`
	Notincluded    []string `json:"notincluded"`
	Canbeincluded  []string `json:"canbeincluded"`
}

/*
Writes all available tags from the json file into tables in order to make them easier to use.
Will be executed once while the server is started.
*/
func initTagRatingOptions(db *gorm.DB) {
	updateTagTable("backend/static_data/dishRatingTags.json", db, MEAL)
	updateTagTable("backend/static_data/cafeteriaRatingTags.json", db, CAFETERIA)
	updateNameTagOptions(db)
	addEntriesForCronJob(db, "averageRatingComputation", 300)
	addEntriesForCronJob(db, "dishNameDownload", 302400) //run twice every week
}

func addEntriesForCronJob(db *gorm.DB, cronName string, interval int32) {
	var exists bool
	res := db.Model(&model.Crontab{}).
		Select("count(*) > 0").
		Where("type LIKE ?", cronName).
		Find(&exists).
		Error

	if res != nil {
		log.Error(res.Error)
	} else if !exists {

		db.Model(&model.Crontab{}).
			Create(&model.Crontab{
				Interval: interval,
				Type:     null.String{NullString: sql.NullString{String: cronName, Valid: true}},
				LastRun:  0,
			})
	}
}

/*
Updates the list of dishtags.
If a Tag with the exact german and english name does not exist yet, it will be created.
Old tags won't be removed to prevent problems with foreign keys.
*/
func updateNameTagOptions(db *gorm.DB) {
	absPathDishNames, _ := filepath.Abs("backend/static_data/dishNameTags.json")
	tagsNames := generateNameTagListFromFile(absPathDishNames)
	var elementID int32
	for _, v := range tagsNames.MultiLanguageNameTags {
		var parentID int32

		potentialTag := db.Model(&cafeteria_rating_models.DishNameTagOption{}).
			Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
			Select("DishNameTagOption").
			Scan(&parentID)

		if potentialTag.RowsAffected == 0 {
			parent := cafeteria_rating_models.DishRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish}

			db.Model(&cafeteria_rating_models.DishNameTagOption{}).
				Create(&parent)
			parentID = parent.DishRatingTagOption
		}

		for _, u := range v.Canbeincluded {
			resultIncluded := db.Model(&cafeteria_rating_models.DishNameTagOptionIncluded{}).
				Where("expression LIKE ? AND NameTagID = ?", u, parentID).
				Select("DishNameTagOptionIncluded").
				Scan(&elementID)
			if resultIncluded.RowsAffected == 0 {
				db.Model(&cafeteria_rating_models.DishNameTagOptionIncluded{}).
					Create(&cafeteria_rating_models.DishNameTagOptionIncluded{
						Expression: u,
						NameTagID:  parentID})
			}
		}
		for _, u := range v.Notincluded {
			resultIncluded := db.Model(&cafeteria_rating_models.DishNameTagOptionExcluded{}).
				Where("expression LIKE ? AND NameTagID = ?", u, parentID).
				Select("DishNameTagOptionExcluded").
				Scan(&elementID)
			if resultIncluded.RowsAffected == 0 {
				db.Model(&cafeteria_rating_models.DishNameTagOptionExcluded{}).
					Create(&cafeteria_rating_models.DishNameTagOptionExcluded{
						Expression: u,
						NameTagID:  parentID})
			}
		}
	}
}

/*
Reads the json file at the given path and checks whether the values have already been inserted into the corresponding table.
If an entry with the same German and English name exists, the entry won't be added.
The TagType is used to identify the corresponding model
*/
func updateTagTable(path string, db *gorm.DB, tagType int) {
	absPathDish, _ := filepath.Abs(path)
	tagsDish := generateRatingTagListFromFile(absPathDish)
	insertModel := getTagModel(tagType, db)
	for _, v := range tagsDish.MultiLanguageTags {
		var result int32
		var affectedRows = 0
		if tagType == CAFETERIA {
			affectedRows = int(db.Model(&cafeteria_rating_models.CafeteriaRatingTagOption{}).
				Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
				Select("cafeteriaRatingTagOption").
				Scan(&result).RowsAffected)
		} else {
			affectedRows = int(db.Model(&cafeteria_rating_models.DishRatingTagOption{}).
				Where("EN LIKE ? AND DE LIKE ?", v.TagNameEnglish, v.TagNameGerman).
				Select("dishRatingTagOption").
				Scan(&result).RowsAffected)
		}

		if affectedRows == 0 {
			println("New entry inserted to Rating Tag Options")
			element := cafeteria_rating_models.DishRatingTagOption{
				DE: v.TagNameGerman,
				EN: v.TagNameEnglish}
			insertModel.
				Create(&element)
		}
	}
}

func getTagModel(tagType int, db *gorm.DB) *gorm.DB {
	if tagType == MEAL {
		return db.Model(&cafeteria_rating_models.DishRatingTagOption{})
	} else {
		return db.Model(&cafeteria_rating_models.CafeteriaRatingTagOption{})
	}
}

func generateNameTagListFromFile(path string) MultiLanguageNameTags {
	byteValue := readFromFile(path)

	var tags MultiLanguageNameTags
	errorUnmarshal := json.Unmarshal(byteValue, &tags)
	if errorUnmarshal != nil {
		log.Error("Error in parsing json:", errorUnmarshal)
	}
	return tags
}

func generateRatingTagListFromFile(path string) MultiLanguageTags {
	byteValue := readFromFile(path)

	var tags MultiLanguageTags
	errorUnmarshal := json.Unmarshal(byteValue, &tags)
	if errorUnmarshal != nil {
		log.Error("Error in parsing json:", errorUnmarshal)
	}
	return tags
}

func readFromFile(path string) []byte {
	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Error("Error in parsing json:", err)
		}
	}(jsonFile)

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	return byteValue
}
