package entity

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"time"
)

type Order struct {
	ID                uint64                    `gorm:"primary_key;column:id;autoIncrement"`
	Invoice           string                    `gorm:"column:invoice"`
	DiscountType      enum_state.DiscountType   `gorm:"column:discount_type"`
	DiscountValue     float32                   `gorm:"column:discount_value"`
	TotalDiscount     float32                   `gorm:"column:total_discount"`
	UserId            uint64                    `gorm:"column:user_id"`
	FirstName         string                    `gorm:"column:first_name"`
	LastName          string                    `gorm:"column:last_name"`
	Email             string                    `gorm:"column:email"`
	Phone             string                    `gorm:"column:phone"`
	PaymentGateway    enum_state.PaymentGateway `gorm:"column:payment_gateway"`
	PaymentMethod     enum_state.PaymentMethod  `gorm:"column:payment_method"`
	PaymentStatus     enum_state.PaymentStatus  `gorm:"column:payment_status"`
	ChannelCode       enum_state.ChannelCode    `gorm:"channel_code"`
	OrderStatus       enum_state.OrderStatus    `gorm:"column:order_status"`
	IsDelivery        bool                      `gorm:"column:is_delivery"`
	DeliveryCost      float32                   `gorm:"column:delivery_cost"`
	CompleteAddress   string                    `gorm:"column:complete_address"`
	Note              string                    `gorm:"column:note"`
	ServiceFee        float32                   `gorm:"column:service_fee"`
	TotalProductPrice float32                   `gorm:"column:total_product_price"`
	TotalFinalPrice   float32                   `gorm:"column:total_final_price"`
	CancellationNotes string                    `gorm:"cancellation_notes"`
	RejectionNotes    string                    `gorm:"rejection_notes"`
	CreatedAt         time.Time                 `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt         time.Time                 `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	OrderProducts     []OrderProduct            `gorm:"foreignKey:order_id;references:id"`
	XenditTransaction *XenditTransactions       `gorm:"foreignKey:order_id;references:id"`
}

func (o *Order) TableName() string {
	return "orders"
}
