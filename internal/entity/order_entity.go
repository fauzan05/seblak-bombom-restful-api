package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type Order struct {
	ID                uint64                `gorm:"primary_key;column:id;autoIncrement"`
	Invoice           string                `gorm:"column:invoice"`
	Amount            float32               `gorm:"column:amount"`
	DiscountValue     float32               `gorm:"column:discount_value"`
	DiscountType      helper.DiscountType   `gorm:"column:discount_type"`
	TotalDiscount     float32               `gorm:"column:total_discount"`
	UserId            uint64                `gorm:"column:user_id"`
	FirstName         string                `gorm:"column:first_name"`
	LastName          string                `gorm:"column:last_name"`
	Email             string                `gorm:"column:email"`
	Phone             string                `gorm:"column:phone"`
	PaymentGateway    helper.PaymentGateway `gorm:"column:payment_gateway"`
	PaymentMethod     helper.PaymentMethod  `gorm:"column:payment_method"`
	ChannelCode       helper.ChannelCode    `gorm:"channel_code"`
	PaymentStatus     helper.PaymentStatus  `gorm:"column:payment_status"`
	OrderStatus       helper.OrderStatus    `gorm:"column:order_status"`
	IsDelivery        bool                  `gorm:"column:is_delivery"`
	DeliveryCost      float32               `gorm:"column:delivery_cost"`
	CompleteAddress   string                `gorm:"column:complete_address"`
	Note              string                `gorm:"column:note"`
	Created_At        time.Time             `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At        time.Time             `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	OrderProducts     []OrderProduct        `gorm:"foreignKey:order_id;references:id"`
	XenditTransaction *XenditTransactions   `gorm:"foreignKey:order_id;references:id"`
}

func (o *Order) TableName() string {
	return "orders"
}
