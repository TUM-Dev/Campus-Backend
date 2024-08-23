package model

// Cafeteria stores all Available cafeterias in the format of the eat-api
type Cafeteria struct {
	Cafeteria int64   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteria;type:int;not null;" json:"canteen" `
	Name      string  `gorm:"column:name;type:text;not null;" json:"name" `
	Address   string  `gorm:"column:address;type:text;not null;" json:"address" `
	Latitude  float64 `gorm:"column:latitude;type:double;not null;" json:"latitude" `
	Longitude float64 `gorm:"column:longitude;type:double;not null;" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (n *Cafeteria) TableName() string {
	return "cafeteria"
}
