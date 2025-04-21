package entity

import "time"

type XenditTransactions struct {
	ID              string    `gorm:"primary_key;column:id"`
	OrderId         uint64    `gorm:"column:order_id"`
	ReferenceId     string    `gorm:"column:reference_id"`
	Amount          float64   `gorm:"column:amount"`
	Currency        string    `gorm:"column:currency"`
	PaymentMethod   string    `gorm:"column:payment_method"`
	PaymentMethodId string    `gorm:"column:payment_method_id"`
	ChannelCode     string    `gorm:"column:channel_code"`
	QrString        string    `gorm:"column:qr_string"`
	Status          string    `gorm:"column:status"`
	Description     string    `gorm:"column:description"`
	FailureCode     string    `gorm:"column:failure_code"`
	Metadata        []byte    `gorm:"column:metadata"`
	ExpiresAt       time.Time `gorm:"column:expires_at"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
	Order           *Order    `gorm:"foreignKey:order_id;references:id"`
}

func (u *XenditTransactions) TableName() string {
	return "xendit_transactions"
}
