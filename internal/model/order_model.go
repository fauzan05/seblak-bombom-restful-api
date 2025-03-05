package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type OrderResponse struct {
	ID                uint64                     `json:"id"`
	Invoice           string                     `json:"invoice"`
	OrderProducts     []OrderProductResponse     `json:"order_products"`
	Amount            string                     `json:"amount"`
	DiscountType      helper.DiscountType        `json:"discount_type"`
	DiscountValue     float32                    `json:"discount_value"`
	TotalDiscount     float32                    `json:"total_discount"`
	UserId            uint64                     `json:"user_id"`
	FirstName         string                     `json:"first_name"`
	LastName          string                     `json:"last_name"`
	Email             string                     `json:"email"`
	Phone             string                     `json:"phone"`
	PaymentMethod     helper.PaymentMethod       `json:"payment_method"`
	PaymentStatus     helper.PaymentStatus       `json:"payment_status"`
	OrderStatus       helper.OrderStatus         `json:"order_status"`
	IsDelivery        bool                       `json:"delivery"`
	DeliveryCost      string                     `json:"delivery_cost"`
	CompleteAddress   string                     `json:"complete_address"`
	Note              string                     `json:"note"`
	CreatedAt         time.Time                  `json:"created_at"`
	UpdatedAt         time.Time                  `json:"updated_at"`
	XenditTransaction *XenditTransactionResponse `json:"xendit_transaction_response,omitempty"`
}

type CreateOrderRequest struct {
	DiscountId      uint64                      `json:"discount_id"`
	UserId          uint64                      `json:"user_id" validate:"required"`
	FirstName       string                      `json:"first_name" validate:"required"`
	LastName        string                      `json:"last_name" validate:"required"`
	Email           string                      `json:"email" validate:"required"`
	Phone           string                      `json:"phone" validate:"required"`
	PaymentMethod   helper.PaymentMethod        `json:"payment_method" validate:"required"`
	ChannelCode     helper.ChannelCode          `json:"channel_code" validate:"required"`
	PaymentGateway  helper.PaymentGateway       `json:"payment_gateway" validate:"required"`
	IsDelivery      bool                        `json:"is_delivery"`
	DeliveryId      uint64                      `json:"delivery_id"`
	CompleteAddress string                      `json:"complete_address" validate:"required"`
	Note            string                      `json:"note"`
	CurrentBalance  float32                     `json:"current_balance"`
	OrderProducts   []CreateOrderProductRequest `json:"order_products" validate:"required"`
}

type GetOrderByCurrentRequest struct {
	ID uint64 `json:"-" validate:"required"` //user id
}

type UpdateOrderRequest struct {
	ID          uint64             `json:"-" validate:"required"` //order id
	OrderStatus helper.OrderStatus `json:"order_status" validate:"required"`
}

type GetOrdersByUserIdRequest struct {
	ID uint64 `json:"-" validate:"required"` //user id
}
