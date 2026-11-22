package model

import "time"

// DishSoldOutReport records a crowdsourced "sold out" report for a dish on a specific Berlin calendar day.
type DishSoldOutReport struct {
	ID         int64     `gorm:"primary_key;autoIncrement;column:id;type:bigint;not null;" json:"id"`
	DishID     int64     `gorm:"column:dish_id;type:bigint;not null;uniqueIndex:idx_dish_sold_out_unique,priority:1;index:idx_dish_sold_out_lookup,priority:1" json:"dish_id"`
	Dish       Dish      `gorm:"foreignKey:DishID;references:Dish;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DeviceHash string    `gorm:"column:device_hash;type:varchar(64);not null;uniqueIndex:idx_dish_sold_out_unique,priority:2" json:"device_hash"`
	ReportDate string    `gorm:"column:report_date;type:char(10);not null;uniqueIndex:idx_dish_sold_out_unique,priority:3;index:idx_dish_sold_out_lookup,priority:2" json:"report_date"`
	ReportedAt time.Time `gorm:"column:reported_at;type:timestamp;not null;default:current_timestamp();" json:"reported_at"`
}
