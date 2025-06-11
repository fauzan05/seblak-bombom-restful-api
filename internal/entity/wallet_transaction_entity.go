package entity

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"time"
)

type WalletTransactions struct {
	ID        uint64                             `gorm:"primary_key;column:id"`
	UserId    uint64                             `gorm:"column:user_id"`
	OrderId   *uint64                            `gorm:"column:order_id"`
	Amount    float32                            `gorm:"column:amount"`
	Type      enum_state.WalletTransactionType   `gorm:"column:type"`
	Source    enum_state.WalletTransactionSource `gorm:"column:source"`
	Note      string                             `gorm:"column:note"`
	CreatedAt time.Time                          `gorm:"column:created_at"`
	UpdatedAt time.Time                          `gorm:"column:updated_at"`
	Order     *Order                             `gorm:"foreignKey:order_id;references:id"`
	User      *User                              `gorm:"foreignKey:user_id;references:id"`
}

func (u *WalletTransactions) TableName() string {
	return "wallet_transactions"
}
