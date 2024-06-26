package model

import (
	"seblak-bombom-restful-api/internal/helper"
)

type OrderResponse struct {
	ID                uint64                    `json:"id,omitempty"`
	Invoice           string                    `json:"invoice,omitempty"`
	OrderProducts     []OrderProductResponse    `json:"order_products,omitempty"`
	Amount            string                    `json:"amount,omitempty"`
	DiscountType      helper.DiscountType       `json:"discount_type,omitempty"`
	DiscountValue     float32                   `json:"discount_value,omitempty"`
	TotalDiscount     float32                   `json:"total_discount,omitempty"`
	UserId            uint64                    `json:"user_id,omitempty"`
	FirstName         string                    `json:"first_name,omitempty"`
	LastName          string                    `json:"last_name,omitempty"`
	Email             string                    `json:"email,omitempty"`
	Phone             string                    `json:"phone,omitempty"`
	PaymentMethod     helper.PaymentMethod      `json:"payment_method,omitempty"`
	PaymentStatus     helper.PaymentStatus      `json:"payment_status,omitempty"`
	DeliveryStatus    helper.DeliveryStatus     `json:"delivery_status,omitempty"`
	IsDelivery        bool                      `json:"delivery,omitempty"`
	DeliveryCost      string                    `json:"delivery_cost,omitempty"`
	CompleteAddress   string                    `json:"complete_address,omitempty"`
	Longitude         float64                   `json:"longitude,omitempty"`
	Latitude          float64                   `json:"latitude,omitempty"`
	Distance          float32                   `json:"distance,omitempty"`
	MidtransSnapOrder MidtransSnapOrderResponse `json:"midtrans_snap_order,omitempty"`
	Note              string                    `json:"note,omitempty"`
	CreatedAt         string                    `json:"created_at,omitempty"`
	UpdatedAt         string                    `json:"updated_at,omitempty"`
}

type CreateOrderRequest struct {
	DiscountCode       string                      `json:"discount_code"`
	UserId             uint64                      `json:"user_id" validate:"required"`
	FirstName          string                      `json:"first_name" validate:"required"`
	LastName           string                      `json:"last_name" validate:"required"`
	Email              string                      `json:"email" validate:"required"`
	Phone              string                      `json:"phone" validate:"required"`
	PaymentMethod      helper.PaymentMethod        `json:"payment_method" validate:"required"`
	IsDelivery         bool                        `json:"is_delivery"`
	CompleteAddress    string                      `json:"complete_address" validate:"required"`
	Longitude          float64                     `json:"longitude"`
	Latitude           float64                     `json:"latitude"`
	Distance           float32                     `json:"distance"`
	Note               string                      `json:"note"`
	OrderProducts      []CreateOrderProductRequest `json:"order_products" validate:"required"`
	IsLonLatTakeManual bool                        `json:"is_lonlat_take_manual"` // apakah lon lat mengambil dari alamat customers (otomatis) atau tidak (manual)
}

type GetOrderByCurrentRequest struct {
	ID uint64 `json:"-" validate:"required"` //user id
}

type UpdateOrderRequest struct {
	ID             uint64                `json:"-" validate:"required"` //order id
	PaymentStatus  helper.PaymentStatus  `json:"payment_status" validate:"required"`
	DeliveryStatus helper.DeliveryStatus `json:"delivery_status" validate:"required"`
}

type GetOrdersByUserIdRequest struct {
	ID uint64 `json:"-" validate:"required"` //user id
}
