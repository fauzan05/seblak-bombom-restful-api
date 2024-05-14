package entity

import (
	"time"
)

// token is a struct that represents a token entity in database table
type Address struct {
	ID              uint64    `gorm:"primary_key;column:id;autoIncrement"`
	UserId          uint64    `gorm:"column:user_id"`
	Regency         string    `gorm:"column:regency"`
	SubDistrict     string    `gorm:"column:subdistrict"`
	CompleteAddress string    `gorm:"column:complete_address"`
	Longitude       float64   `gorm:"column:longitude"`
	Latitude        float64   `gorm:"column:latitude"`
	IsMain          bool      `gorm:"column:is_main"`
	Created_At      time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At      time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User            *User     `gorm:"foreignKey:user_id;references:id"`
}

func (u *Address) TableName() string {
	return "addresses"
}
