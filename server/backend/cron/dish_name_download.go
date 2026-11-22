package cron

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/TUM-Dev/Campus-Backend/server/internal/fcm"
	"github.com/TUM-Dev/Campus-Backend/server/internal/webfetch"
	"github.com/TUM-Dev/Campus-Backend/server/model"
	"github.com/guregu/null"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// eatAPICanteen is one entry of the eat-api canteen directory (enums/canteens.json).
// The eat-api names every canteen twice: EnumName ("MENSA_GARCHING") and CanteenID ("mensa-garching").
// The canteens table stores the EnumName; use model.CanteenEnumName / model.CanteenEatAPIID to convert.
type eatAPICanteen struct {
	EnumName  string          `json:"enum_name"`
	CanteenID string          `json:"canteen_id"`
	Location  CanteenLocation `json:"location"`
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

type CanteenDishPrices struct {
	Students *CanteenDishPrice `json:"students,omitempty"`
}

type CanteenDishPrice struct {
	BasePrice *float64 `json:"base_price,omitempty"`
}

type CanteenDish struct {
	Name     string             `json:"name"`
	DishType string             `json:"dish_type"`
	Labels   []string           `json:"labels,omitempty"`
	Prices   *CanteenDishPrices `json:"prices,omitempty"`
}

func (d CanteenDish) encodedLabels() null.String {
	if len(d.Labels) == 0 {
		return null.String{}
	}
	cleaned := make([]string, 0, len(d.Labels))
	for _, l := range d.Labels {
		if t := strings.ToUpper(strings.TrimSpace(l)); t != "" {
			cleaned = append(cleaned, t)
		}
	}
	if len(cleaned) == 0 {
		return null.String{}
	}
	joined := strings.Join(cleaned, ",")
	if len(joined) > 255 {
		joined = joined[:255]
	}
	return null.StringFrom(joined)
}

func (d CanteenDish) studentPriceCents() null.Int {
	if d.Prices == nil || d.Prices.Students == nil || d.Prices.Students.BasePrice == nil {
		return null.Int{}
	}
	price := *d.Prices.Students.BasePrice
	if price < 0 || math.IsNaN(price) || math.IsInf(price, 0) {
		return null.Int{}
	}
	return null.IntFrom(int64(math.Round(price * 100)))
}

// dishNameDownloadCron refreshes metadata and dish names independently, so a
// directory outage does not prevent updates for canteens already in the DB.
func (c *CronService) dishNameDownloadCron() error {
	namesErr := downloadCanteenNames(c)
	dishesErr := downloadDailyDishes(c)
	return errors.Join(namesErr, dishesErr)
}

type isoWeekTarget struct {
	year     int
	week     int
	optional bool
}

func downloadDailyDishes(c *CronService) error {
	berlin, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return err
	}
	return downloadDailyDishesAt(c, time.Now().In(berlin))
}

func downloadDailyDishesAt(c *CronService, now time.Time) error {
	var canteens []model.Canteen
	if err := c.db.Find(&canteens).Error; err != nil {
		return err
	}
	currYear, currWeek := now.ISOWeek()
	nextYear, nextWeek := now.AddDate(0, 0, 7).ISOWeek()
	targets := []isoWeekTarget{
		{year: currYear, week: currWeek, optional: false},
		{year: nextYear, week: nextWeek, optional: true},
	}

	var failures []error
	for _, canteen := range canteens {
		name := model.CanteenEatAPIID(canteen.Name)
		for _, target := range targets {
			address := fmt.Sprintf("https://tum-dev.github.io/eat-api/%s/%d/%02d.json", name, target.year, target.week)
			var menu CanteenDays
			if err := webfetch.JSON(address, &menu); err != nil {
				var statusErr *webfetch.StatusError
				if target.optional && errors.As(err, &statusErr) && statusErr.StatusCode == http.StatusNotFound {
					continue
				}
				failures = append(failures, fmt.Errorf("menu %s %d-W%02d: %w", name, target.year, target.week, err))
				continue
			}
			if err := storeCanteenMenu(c.db, canteen.CanteenID, menu, target.year, target.week); err != nil {
				failures = append(failures, fmt.Errorf("storing menu %s %d-W%02d: %w", name, target.year, target.week, err))
			}
		}
	}
	if sender, err := fcm.NewFromEnv(context.Background()); err != nil {
		log.WithError(err).Warn("FCM client unavailable for favorite dish alerts")
	} else if sender != nil {
		for _, target := range targets {
			if err := notifyFavoriteDishAlerts(context.Background(), c.db, sender, target.year, target.week); err != nil {
				log.WithError(err).Warn("failed to dispatch favorite dish push alerts")
			}
		}
	}
	return errors.Join(failures...)
}

func storeCanteenMenu(db *gorm.DB, canteenID int64, menu CanteenDays, year, week int) error {
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
				labels := entry.encodedLabels()
				priceCents := entry.studentPriceCents()
				var dish model.Dish
				err := tx.First(&dish, "name = ? AND cafeteriaID = ?", entry.Name, canteenID).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					dish = model.Dish{
						Name:              entry.Name,
						Type:              entry.DishType,
						CanteenID:         canteenID,
						Labels:            labels,
						PriceStudentCents: priceCents,
					}
					if err := tx.Create(&dish).Error; err != nil {
						return err
					}
					addDishTagsToMapping(dish.Dish, dish.Name, tx)
				} else if err != nil {
					return err
				} else {
					updates := map[string]any{}
					if dish.Type != entry.DishType {
						updates["type"] = entry.DishType
					}
					if labels.Valid && dish.Labels != labels {
						updates["labels"] = labels
					}
					if priceCents.Valid && dish.PriceStudentCents != priceCents {
						updates["price_student_cents"] = priceCents
					}
					if len(updates) > 0 {
						if err := tx.Model(&dish).Updates(updates).Error; err != nil {
							return err
						}
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
	var entries []eatAPICanteen
	if err := webfetch.JSON("https://tum-dev.github.io/eat-api/enums/canteens.json", &entries); err != nil {
		return fmt.Errorf("fetching canteen directory: %w", err)
	}
	for _, entry := range entries {
		name := entry.EnumName
		if name == "" {
			name = model.CanteenEnumName(entry.CanteenID)
		}
		canteen := model.Canteen{Name: name, Address: entry.Location.Address, Latitude: entry.Location.Latitude, Longitude: entry.Location.Longitude}
		var existing model.Canteen
		if err := c.db.First(&existing, "name = ?", name).Error; errors.Is(err, gorm.ErrRecordNotFound) {
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
// Checks whether the dish name includes any included expressions while excluding excluded ones,
// deduplicates matched tag IDs, and idempotently stores mappings in dish_to_dish_name_tags.
func addDishTagsToMapping(dishID int64, dishName string, db *gorm.DB) {
	lowercaseDish := strings.ToLower(dishName)
	var includedTags []int64
	if err := db.Model(&model.IncludedDishNameTagOption{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&includedTags).Error; err != nil {
		log.WithError(err).Error("Error while querying all included expressions for the CanteenDish: ", lowercaseDish)
		return
	}

	var excludedTags []int64
	if err := db.Model(&model.ExcludedDishNameTagOption{}).
		Where("? LIKE CONCAT('%', expression ,'%')", lowercaseDish).
		Select("nameTagID").
		Scan(&excludedTags).Error; err != nil {
		log.WithError(err).Error("Error while querying all excluded expressions for the CanteenDish: ", lowercaseDish)
		return
	}

	excludedSet := make(map[int64]bool, len(excludedTags))
	for _, id := range excludedTags {
		excludedSet[id] = true
	}

	seen := make(map[int64]bool, len(includedTags))
	for _, nametagID := range includedTags {
		if excludedSet[nametagID] || seen[nametagID] {
			continue
		}
		seen[nametagID] = true
		mapping := model.DishToDishNameTag{DishID: dishID, NameTagID: nametagID}
		if err := db.Where(map[string]any{"dishID": dishID, "nameTagID": nametagID}).FirstOrCreate(&mapping).Error; err != nil {
			fields := log.Fields{"dishID": dishID, "nametagID": nametagID}
			log.WithError(err).WithFields(fields).Error("creating dish name tag mapping")
		}
	}
}

// BackfillMissingDishNameTags ensures every dish in the database has all currently matching
// name-tag mappings without creating duplicates.
func BackfillMissingDishNameTags(db *gorm.DB) error {
	var dishes []model.Dish
	if err := db.Select("dish, name").Find(&dishes).Error; err != nil {
		return err
	}
	for _, d := range dishes {
		addDishTagsToMapping(d.Dish, d.Name, db)
	}
	return nil
}
