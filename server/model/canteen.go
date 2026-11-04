package model

// Canteen stores all Available cafeterias in the format of the eat-api
type Canteen struct {
	Cafeteria int64   `gorm:"primary_key;autoIncrement;column:cafeteria;type:int;not null;" json:"canteen" `
	Name      string  `gorm:"column:name;type:text;not null;" json:"name" `
	Address   string  `gorm:"column:address;type:text;not null;" json:"address" `
	Latitude  float64 `gorm:"column:latitude;type:double;not null;" json:"latitude" `
	Longitude float64 `gorm:"column:longitude;type:double;not null;" json:"longitude"`
}
