package model

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"time"
)

type OrderResponse struct {
	ID                uint64                     `json:"id"`
	Invoice           string                     `json:"invoice"`
	DiscountType      enum_state.DiscountType    `json:"discount_type"`
	DiscountValue     float32                    `json:"discount_value"`
	TotalDiscount     float32                    `json:"total_discount"`
	UserId            uint64                     `json:"user_id"`
	FirstName         string                     `json:"first_name"`
	LastName          string                     `json:"last_name"`
	Email             string                     `json:"email"`
	Phone             string                     `json:"phone"`
	PaymentGateway    enum_state.PaymentGateway  `json:"payment_gateway"`
	PaymentMethod     enum_state.PaymentMethod   `json:"payment_method"`
	PaymentStatus     enum_state.PaymentStatus   `json:"payment_status"`
	ChannelCode       enum_state.ChannelCode     `json:"channel_code"`
	OrderStatus       enum_state.OrderStatus     `json:"order_status"`
	IsDelivery        bool                       `json:"delivery"`
	DeliveryCost      float32                    `json:"delivery_cost"`
	CompleteAddress   string                     `json:"complete_address"`
	Note              string                     `json:"note"`
	ServiceFee        float32                    `json:"service_fee"`
	TotalProductPrice float32                    `json:"total_product_price"`
	TotalFinalPrice   float32                    `json:"total_final_price"`
	CancellationNotes string                     `json:"cancellation_notes"`
	CreatedAt         helper_others.TimeRFC3339  `json:"created_at"`
	UpdatedAt         helper_others.TimeRFC3339  `json:"updated_at"`
	OrderProducts     []OrderProductResponse     `json:"order_products"`
	XenditTransaction *XenditTransactionResponse `json:"xendit_transaction_response,omitempty"`
}

type CreateOrderRequest struct {
	DiscountId      uint64                    `json:"discount_id"`
	UserId          uint64                    `json:"user_id" validate:"required"`
	FirstName       string                    `json:"first_name" validate:"required"`
	LastName        string                    `json:"last_name" validate:"required"`
	Email           string                    `json:"email" validate:"required"`
	Phone           string                    `json:"phone" validate:"required"`
	PaymentMethod   enum_state.PaymentMethod  `json:"payment_method" validate:"required"`
	ChannelCode     enum_state.ChannelCode    `json:"channel_code" validate:"required"`
	PaymentGateway  enum_state.PaymentGateway `json:"payment_gateway" validate:"required"`
	IsDelivery      bool                      `json:"is_delivery"`
	DeliveryId      uint64                    `json:"delivery_id"`
	CompleteAddress string                    `json:"complete_address" validate:"required"`
	Note            string                    `json:"note"`
	CurrentBalance  float32                   `json:"current_balance"`
	OrderProducts   []OrderProductResponse    `json:"order_products" validate:"required"`
	Lang            enum_state.Languange      `json:"-"`
	TimeZone        time.Location             `json:"-"`
	BaseFrontEndURL string                    `json:"-"`
}

type GetOrderByCurrentRequest struct {
	ID uint64 `json:"-" validate:"required"` //user id
}

type UpdateOrderRequest struct {
	ID                uint64                 `json:"-" validate:"required"` //order id
	OrderStatus       enum_state.OrderStatus `json:"order_status"`
	CancellationNotes string                 `json:"cancellation_notes"`
	Lang              enum_state.Languange   `json:"-"`
	TimeZone          time.Location          `json:"-"`
	BaseFrontEndURL   string                 `json:"-"`
}

type GetOrdersByUserIdRequest struct {
	ID uint64 `json:"-" validate:"required"` //user id
}
