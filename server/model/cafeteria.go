package model

// Cafeteria stores all Available cafeterias in the format of the eat-api
type Cafeteria struct {
	Cafeteria int32   `gorm:"primary_key;AUTO_INCREMENT;column:cafeteria;type:int;not null;" json:"cafeteria" `
	Name      string  `gorm:"column:name;type:mediumtext;not null;" json:"name" `
	Address   string  `gorm:"column:address;type:text;not null;" json:"address" `
	Latitude  float32 `gorm:"column:latitude;type:float;not null;" json:"latitude" `
	Longitude float32 `gorm:"column:longitude;type:float;not null;" json:"longitude"`
}

// TableName sets the insert table name for this struct type
func (n *Cafeteria) TableName() string {
	return "cafeteria"
}
