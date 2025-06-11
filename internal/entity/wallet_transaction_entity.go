package entity

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"time"
)

type WalletTransactions struct {
	ID              uint64                             `gorm:"primary_key;column:id"`
	UserId          uint64                             `gorm:"column:user_id"`
	OrderId         *uint64                            `gorm:"column:order_id"`
	Amount          float32                            `gorm:"column:amount"`
	FlowType        enum_state.WalletFlowType          `gorm:"column:flow_type"`
	TransactionType enum_state.WalletTransactionType   `gorm:"column:transaction_type"`
	PaymentMethod   enum_state.PaymentMethod           `gorm:"column:payment_method"`
	Status          enum_state.WalletTransactionStatus `gorm:"column:status"`
	ReferenceNumber string                             `gorm:"column:reference_number"`
	Note            string                             `gorm:"column:note"`
	AdminNote       string                             `gorm:"column:admin_note"`
	ProcessedBy     *uint64                            `gorm:"column:processed_by"`
	ProcessedAt     *time.Time                         `gorm:"column:processed_at"`
	CreatedAt       time.Time                          `gorm:"column:created_at"`
	UpdatedAt       time.Time                          `gorm:"column:updated_at"`
	Order           *Order                             `gorm:"foreignKey:order_id;references:id"`
	User            *User                              `gorm:"foreignKey:user_id;references:id"`
}

func (u *WalletTransactions) TableName() string {
	return "wallet_transactions"
}
