package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type Order struct {
	ID                 uint64                `gorm:"primary_key;column:id;autoIncrement"`
	Invoice            string                `gorm:"column:invoice"`
	ProductId          uint64                `gorm:"column:product_id"`
	ProductName        string                `gorm:"column:product_name"`
	ProductDescription string                `gorm:"column:product_description"`
	Price              int                   `gorm:"column:price"`
	Quantity           int                   `gorm:"column:quantity"`
	Amount             int                   `gorm:"column:amount"`
	DiscountValue      int                   `gorm:"column:discount_value"`
	DiscountType       helper.DiscountType   `gorm:"column:discount_type"`
	UserId             uint64                `gorm:"column:user_id"`
	FirstName          string                `gorm:"column:first_name"`
	LastName           string                `gorm:"column:last_name"`
	Email              string                `gorm:"column:email"`
	Phone              string                `gorm:"column:phone"`
	PaymentMethod      helper.PaymentMethod  `gorm:"column:payment_method"`
	PaymentStatus      helper.PaymentStatus  `gorm:"column:payment_status"`
	DeliveryStatus     helper.DeliveryStatus `gorm:"column:delivery_status"`
	IsDelivery         bool                  `gorm:"column:is_delivery"`
	DeliveryCost       int                   `gorm:"column:delivery_cost"`
	CategoryName       string                `gorm:"column:category_name"`
	CompleteAddress    string                `gorm:"column:complete_address"`
	GoogleMapLink      string                `gorm:"column:google_map_link"`
	Distance           int                   `gorm:"column:distance"`
	Created_At         time.Time             `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At         time.Time             `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (o *Order) TableName() string {
	return "orders"
}