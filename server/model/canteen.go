package model

import "strings"

// Canteen stores all canteens known to the eat-api. Name holds the eat-api enum_name (e.g. MENSA_GARCHING).
type Canteen struct {
	CanteenID int64   `gorm:"primary_key;autoIncrement;column:cafeteria;type:int;not null;" json:"canteen" `
	Name      string  `gorm:"column:name;type:text;not null;" json:"name" `
	Address   string  `gorm:"column:address;type:text;not null;" json:"address" `
	Latitude  float64 `gorm:"column:latitude;type:double;not null;" json:"latitude" `
	Longitude float64 `gorm:"column:longitude;type:double;not null;" json:"longitude"`
}

// CanteenEnumName converts an eat-api canteen identifier into the enum_name form stored in Canteen.Name.
// The eat-api uses two identifiers per canteen: canteen_id ("mensa-garching") and enum_name ("MENSA_GARCHING").
// Both forms are accepted, so clients may send either one.
func CanteenEnumName(id string) string {
	return strings.ReplaceAll(strings.ToUpper(strings.TrimSpace(id)), "-", "_")
}

// CanteenEatAPIID converts an eat-api canteen identifier into the canteen_id form ("mensa-garching").
// This is the form used in eat-api URLs and in the canteen_head_counts table.
func CanteenEatAPIID(id string) string {
	return strings.ReplaceAll(strings.ToLower(strings.TrimSpace(id)), "_", "-")
}
