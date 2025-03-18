package entity

import "time"

type XenditPayout struct {
	ID                string     `gorm:"primary_key;column:id"`
	UserID            uint64     `gorm:"column:user_id"`
	BusinessID        string     `gorm:"column:business_id"`
	ReferenceID       string     `gorm:"column:reference_id"`
	Amount            float64    `gorm:"column:amount"`
	Currency          string     `gorm:"column:currency"`
	Description       string     `gorm:"column:description"`
	ChannelCode       string     `gorm:"column:channel_code"`
	AccountNumber     string     `gorm:"column:account_number"`
	AccountHolderName string     `gorm:"column:account_holder_name"`
	Status            string     `gorm:"column:status"`
	CreatedAt         *time.Time `gorm:"column:created_at"`
	UpdatedAt         *time.Time `gorm:"column:updated_at"`
	EstimatedArrival  *time.Time `gorm:"column:estimated_arrival"`
	User              *User      `gorm:"foreignKey:user_id;references:id"`
}

func (u *XenditPayout) TableName() string {
	return "xendit_payouts"
}
