package entity

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"time"
)

// wallet is a struct that represents a wallet entity in database table
type Wallet struct {
	ID        uint64                  `gorm:"primary_key;column:id;autoIncrement"`
	Balance   float32                 `gorm:"column:balance"`
	UserId    uint64                  `gorm:"column:user_id"`
	Status    enum_state.WalletStatus `gorm:"column:status"`
	CreatedAt time.Time               `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time               `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User      *User                   `gorm:"foreignKey:user_id;references:id"`
}

func (u *Wallet) TableName() string {
	return "wallets"
}
