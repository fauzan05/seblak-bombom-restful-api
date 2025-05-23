package entity

import (
	"time"

	"gorm.io/gorm"
)

// token is a struct that represents a token entity in database table
type Address struct {
	ID              uint64         `gorm:"primary_key;column:id;autoIncrement"`
	UserId          uint64         `gorm:"column:user_id"`
	DeliveryId      uint64         `gorm:"column:delivery_id"`
	CompleteAddress string         `gorm:"column:complete_address"`
	GoogleMapsLink  string         `gorm:"column:google_maps_link"`
	IsMain          bool           `gorm:"column:is_main"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at"`
	User            *User          `gorm:"foreignKey:user_id;references:id"`
	Delivery        *Delivery      `gorm:"foreignKey:delivery_id;references:id"`
}

func (u *Address) TableName() string {
	return "addresses"
}
