package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type Payout struct {
	ID                uint64              `gorm:"primary_key;column:id;autoIncrement"`
	UserId            uint64              `gorm:"column:user_id"`
	XenditPayoutId    string              `gorm:"column:xendit_payout_id"`
	Amount            float32             `gorm:"column:amount"`
	Currency          string              `gorm:"column:currency"`
	Method            helper.PayoutMethod `gorm:"column:method"`
	Status            helper.PayoutStatus `gorm:"column:status"`
	Notes             string              `gorm:"column:notes"`
	CancellationNotes string              `gorm:"column:cancellation_notes"`
	Created_At        time.Time           `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At        time.Time           `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User              *User               `gorm:"foreignKey:user_id;references:id"`
	XenditPayout      *XenditPayout       `gorm:"foreignKey:xendit_payout_id;references:id"`
}

func (u *Payout) TableName() string {
	return "payouts"
}
