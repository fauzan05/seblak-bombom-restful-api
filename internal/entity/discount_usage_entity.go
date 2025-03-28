package entity

import (
	"time"
)

type DiscountUsage struct {
	ID         uint64    `gorm:"primary_key;column:id;autoIncrement"`
	UserId     uint64    `gorm:"column:user_id"`
	CouponId   uint64    `gorm:"column:coupon_id"`
	UsageCount int       `gorm:"column:usage_count"`
	LastUsed   time.Time `gorm:"column:last_used"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *DiscountUsage) TableName() string {
	return "discount_usage"
}
