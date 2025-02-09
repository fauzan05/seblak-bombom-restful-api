package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

// wallet is a struct that represents a wallet entity in database table
type Wallet struct {
	ID         uint64              `gorm:"primary_key;column:id;autoIncrement"`
	UserId     uint64              `gorm:"column:user_id"`
	Balance    float32             `gorm:"column:balance"`
	Status     helper.WalletStatus `gorm:"column:status"`
	Created_At time.Time           `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At time.Time           `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User       *User               `gorm:"foreignKey:user_id;references:id"`
}

func (u *Wallet) TableName() string {
	return "wallets"
}
