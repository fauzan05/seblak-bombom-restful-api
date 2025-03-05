package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"

	"gorm.io/gorm"
)

// user is a struct that represents a user entity in database table
type User struct {
	ID        uint64         `gorm:"primary_key;column:id;autoIncrement"`
	Name      Name           `gorm:"embedded"`
	Email     string         `gorm:"column:email"`
	Phone     string         `gorm:"column:phone"`
	Password  string         `gorm:"column:password"`
	Role      helper.Role    `gorm:"column:role"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
	Token     Token          `gorm:"foreignKey:user_id;references:id"`
	Addresses []Address      `gorm:"foreignKey:user_id;references:id"`
	Cart      *Cart          `gorm:"foreignKey:user_id;references:id"`
	Wallet    *Wallet        `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}

type Name struct {
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
}
