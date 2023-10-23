package model

// Canteen stores all Available cafeterias in the format of the eat-api
type Canteen struct {
	Cafeteria int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteria;type:int;not null;" json:"canteen" `
	Name      string  `gorm:"column:name;type:mediumtext;not null;" json:"name" `
	Address   string  `gorm:"column:address;type:text;not null;" json:"address" `
	Latitude  float32 `gorm:"column:latitude;type:float;not null;" json:"latitude" `
	Longitude float32 `gorm:"column:longitude;type:float;not null;" json:"longitude"`
}
