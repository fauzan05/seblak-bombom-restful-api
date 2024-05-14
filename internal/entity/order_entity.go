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
	PaymentMethod     helper.PaymentMethod  `gorm:"column:payment_method"`
	PaymentStatus     helper.PaymentStatus  `gorm:"column:payment_status"`
	DeliveryStatus    helper.DeliveryStatus `gorm:"column:delivery_status"`
	IsDelivery        bool                  `gorm:"column:is_delivery"`
	DeliveryCost      float32               `gorm:"column:delivery_cost"`
	CompleteAddress   string                `gorm:"column:complete_address"`
	Longitude         float64               `gorm:"column:longitude"`
	Latitude          float64               `gorm:"column:latitude"`
	Distance          float32               `gorm:"column:distance"`
	Created_At        time.Time             `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At        time.Time             `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	OrderProducts     []OrderProduct        `gorm:"foreignKey:order_id;references:id"`
	MidtransSnapOrder *MidtransSnapOrder    `gorm:"foreignKey:order_id;references:id"`
}

func (o *Order) TableName() string {
	return "orders"
}
