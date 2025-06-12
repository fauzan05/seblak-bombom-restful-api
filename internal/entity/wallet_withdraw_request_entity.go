package entity

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"time"
)

type WalletWithdrawRequests struct {
	ID               uint64                           `gorm:"primary_key;column:id"`
	UserId           uint64                           `gorm:"column:user_id"`
	Amount           float32                          `gorm:"column:amount"`
	Method           enum_state.WalletWithdrawRequest `gorm:"column:method"`
	BankName         string                           `gorm:"column:bank_name"`
	BankAcountNumber string                           `gorm:"column:bank_account_number"`
	BankAcountName   string                           `gorm:"column:bank_account_name"`
	Status           enum_state.WalletWithdrawRequest `gorm:"column:status"`
	Note             string                           `gorm:"column:note"`
	RejectionNotes   string                           `gorm:"column:rejection_notes"`
	ProcessedBy      *uint64                          `gorm:"column:processed_by"`
	ProcessedAt      *time.Time                       `gorm:"column:processed_at"`
	CreatedAt        time.Time                        `gorm:"column:created_at"`
	UpdatedAt        time.Time                        `gorm:"column:updated_at"`
	User             *User                            `gorm:"foreignKey:user_id;references:id"`
}

func (u *WalletWithdrawRequests) TableName() string {
	return "wallet_withdraw_requests"
}
