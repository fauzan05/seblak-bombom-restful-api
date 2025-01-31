package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type DiscountCoupon struct {
	ID          uint64              `gorm:"primary_key;column:id;autoIncrement"`
	Name        string              `gorm:"column:name"`
	Description string              `gorm:"column:description"`
	Code        string              `gorm:"column:code"`
	Value       float32             `gorm:"column:value"`
	Type        helper.DiscountType `gorm:"column:type"`
	Start       time.Time           `gorm:"column:start"`
	End         time.Time           `gorm:"column:end"`
	Status      bool                `gorm:"column:status"` // enable/disable = true/false
	Created_At  time.Time           `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At  time.Time           `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *DiscountCoupon) TableName() string {
	return "discount_coupons"
}
