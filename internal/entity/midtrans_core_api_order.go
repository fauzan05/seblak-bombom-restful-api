package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type MidtransCoreAPIOrder struct {
	ID                uint64                   `gorm:"primaryKey;column:id;autoIncrement"`
	OrderId           uint64                   `gorm:"column:order_id"`
	MidtransOrderId   string                   `gorm:"column:midtrans_order_id"`
	StatusCode        string                   `gorm:"column:status_code"`
	StatusMessage     string                   `gorm:"column:status_message"`
	TransactionId     string                   `gorm:"column:transaction_id"`
	GrossAmount       float32                  `gorm:"column:gross_amount"`
	Currency          string                   `gorm:"column:currency"`
	PaymentType       string                   `gorm:"column:payment_type"`
	TransactionTime   time.Time                `gorm:"column:transaction_time"`
	ExpiryTime        time.Time                `gorm:"column:expiry_time"`
	TransactionStatus helper.TransactionStatus `gorm:"column:transaction_status"`
	FraudStatus       string                   `gorm:"column:fraud_status"`
	CreatedAt         time.Time                `gorm:"autoCreateTime;<-:create"`
	UpdatedAt         time.Time                `gorm:"autoCreateTime;autoUpdateTime"`
	Order             *Order                   `gorm:"foreignKey:order_id;references:id"`
	Actions           []Action                 `gorm:"foreignKey:midtrans_core_api_orders_id;references:id"`
}

func (m *MidtransCoreAPIOrder) TableName() string {
	return "midtrans_core_api_orders"
}
