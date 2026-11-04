package cron

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CafeteriaName struct {
	Name     string          `json:"enum_name"`
	Location CanteenLocation `json:"location"`
}

type CafeteriaWithID struct {
	Name      string `json:"name"`
	Cafeteria int64  `json:"cafeteria"`
}

type CanteenLocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Address   string  `json:"address"`
}

type CanteenDays struct {
	Days []CanteenDate `json:"days"`
}

type CanteenDate struct {
	Date  string        `json:"date"`
	Dates []CanteenDish `json:"dishes"`
}

type CanteenDish struct {
	Name     string `json:"name"`
	DishType string `json:"dish_type"`
}

// dishNameDownloadCron refreshes metadata and dish names independently, so a
// directory outage does not prevent updates for cafeterias already in the DB.
func (c *CronService) dishNameDownloadCron() error {
	namesErr := downloadCanteenNames(c)
	dishesErr := downloadDailyDishes(c)
	return errors.Join(namesErr, dishesErr)
}

func downloadDailyDishes(c *CronService) error {
	var canteens []model.Canteen
	if err := c.db.Find(&canteens).Error; err != nil {
		return err
	}
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return err
	}
	year, week := time.Now().In(berlin).ISOWeek()
	var failures []error
	for _, canteen := range canteens {
		name := strings.ReplaceAll(strings.ToLower(canteen.Name), "_", "-")
		address := fmt.Sprintf("https://tum-dev.github.io/eat-api/%s/%d/%02d.json", name, year, week)
		var menu CanteenDays
		if err := webfetch.JSON(address, &menu); err != nil {
			failures = append(failures, fmt.Errorf("menu %s: %w", name, err))
			continue
		}
		if err := storeCanteenMenu(c.db, canteen.Cafeteria, menu, year, week); err != nil {
			failures = append(failures, fmt.Errorf("storing menu %s: %w", name, err))
		}
	}
	return errors.Join(failures...)
}

func storeCanteenMenu(db *gorm.DB, cafeteriaID int64, menu CanteenDays, year, week int) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, day := range menu.Days {
			date, err := time.Parse("2006-01-02", day.Date)
			if err != nil {
				return fmt.Errorf("invalid menu date %q: %w", day.Date, err)
			}
			dateYear, dateWeek := date.ISOWeek()
			if dateYear != year || dateWeek != week {
				return fmt.Errorf("menu date %s is outside requested week", day.Date)
			}
			// Menus can omit closed days. Array indices are not weekday numbers.
			weekday := int32((int(date.Weekday()) + 6) % 7)
			for _, entry := range day.Dates {
				var dish model.Dish
				err := tx.First(&dish, "name = ? AND cafeteriaID = ?", entry.Name, cafeteriaID).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					dish = model.Dish{Name: entry.Name, Type: entry.DishType, CafeteriaID: cafeteriaID}
					if err := tx.Create(&dish).Error; err != nil {
						return err
					}
					addDishTagsToMapping(dish.Dish, dish.Name, tx)
				} else if err != nil {
					return err
				} else if dish.Type != entry.DishType {
					if err := tx.Model(&dish).Update("type", entry.DishType).Error; err != nil {
						return err
					}
				}
				// Retry partial imports and pick up menus published later in the week.
				// Include zero-valued Monday explicitly in the lookup.
				key := map[string]any{"dishID": dish.Dish, "year": year, "week": week, "day": weekday}
				weekly := model.DishesOfTheWeek{DishID: dish.Dish, Year: int32(year), Week: int32(week), Day: weekday}
				if err := tx.Where(key).FirstOrCreate(&weekly).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func downloadCanteenNames(c *CronService) error {
	var names []CafeteriaName
	if err := webfetch.JSON("https://tum-dev.github.io/eat-api/enums/canteens.json", &names); err != nil {
		return fmt.Errorf("fetching cafeteria directory: %w", err)
	}
	for _, name := range names {
		canteen := model.Canteen{Name: name.Name, Address: name.Location.Address, Latitude: name.Location.Latitude, Longitude: name.Location.Longitude}
		var existing model.Canteen
		if err := c.db.First(&existing, "name = ?", name.Name).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			if err := c.db.Create(&canteen).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else {
			if err := c.db.Model(&existing).Updates(&canteen).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// addDishTagsToMapping
// Checks whether the CanteenDish name includes one of the expressions for the excluded tags as well as the included tags.
// The corresponding tags for all identified DishNames will be saved in the table DishNameTags.
func addDishTagsToMapping(dishID int64, dishName string, db *gorm.DB) {
	lowercaseDish := strings.ToLower(dishName)
	var includedTags []int64
	if err := db.Model(&model.IncludedDishNameTagOption{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&includedTags).Error; err != nil {
		log.WithError(err).Error("Error while querying all included expressions for the CanteenDish: ", lowercaseDish)
	}

	var excludedTags []int64
	if err := db.Model(&model.ExcludedDishNameTagOption{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&excludedTags).Error; err != nil {
		log.WithError(err).Error("Error while querying all excluded expressions for the CanteenDish: ", lowercaseDish)
	}

	//set all entries in included to -1 if the excluded tag was recognised for this tag rating.
	if len(excludedTags) > 0 {
		for _, a := range excludedTags {
			if i := contains(includedTags, a); i != -1 {
				includedTags[i] = -1
			}
		}
	}

	for _, nametagID := range includedTags {
		if nametagID != -1 {
			if err := db.Create(&model.DishToDishNameTag{
				DishID:    dishID,
				NameTagID: nametagID,
			}).Error; err != nil {
				fields := log.Fields{"dishID": dishID, "nametagID": nametagID}
				log.WithError(err).WithFields(fields).Error("creating a new entry")
			}
		}
	}
}
func contains(s []int64, e int64) int64 {
	for i, a := range s {
		if a == e {
			return int64(i)
		}
	}
	return -1
}
