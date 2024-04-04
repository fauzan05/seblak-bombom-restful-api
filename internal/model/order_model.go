package model

import "seblak-bombom-restful-api/internal/helper"

type OrderResponse struct {
	ID                 uint64                `json:"id,omitempty"`
	Invoice            string                `json:"invoice,omitempty"`
	ProductId          uint64                `json:"product_id,omitempty"`
	ProductName        string                `json:"product_name,omitempty"`
	ProductDescription string                `json:"product_description,omitempty"`
	Price              string                `json:"price,omitempty"`
	Quantity           int                   `json:"quantity,omitempty"`
	Amount             string                `json:"amount,omitempty"`
	DiscountType       helper.DiscountType   `json:"discount_type,omitempty"`
	DiscountValue      float32               `json:"discount_value,omitempty"`
	UserId             uint64                `json:"user_id,omitempty"`
	FirstName          string                `json:"first_name,omitempty"`
	LastName           string                `json:"last_name,omitempty"`
	Email              string                `json:"email,omitempty"`
	Phone              string                `json:"phone,omitempty"`
	PaymentMethod      helper.PaymentMethod  `json:"payment_method,omitempty"`
	PaymentStatus      helper.PaymentStatus  `json:"payment_status,omitempty"`
	DeliveryStatus     helper.DeliveryStatus `json:"delivery_status,omitempty"`
	IsDelivery         bool                  `json:"delivery,omitempty"`
	DeliveryCost       string                `json:"delivery_cost,omitempty"`
	CategoryName       string                `json:"category_name,omitempty"`
	CompleteAddress    string                `json:"complete_address,omitempty"`
	GoogleMapLink      string                `json:"google_map_link,omitempty"`
	Distance           float32               `json:"distance,omitempty"`
	CreatedAt          string                `json:"created_at,omitempty"`
	UpdatedAt          string                `json:"updated_at,omitempty"`
}

type CreateOrderRequest struct {
	ProductId       uint64               `json:"product_id" validate:"required"`
	Quantity        int                  `json:"quantity" validate:"required"`
	DiscountCode    string               `json:"discount_code"`
	UserId          uint64               `json:"user_id" validate:"required"`
	FirstName       string               `json:"first_name" validate:"required"`
	LastName        string               `json:"last_name" validate:"required"`
	Email           string               `json:"email" validate:"required"`
	Phone           string               `json:"phone" validate:"required"`
	PaymentMethod   helper.PaymentMethod `json:"payment_method" validate:"required"`
	IsDelivery      bool                 `json:"is_delivery"`
	CompleteAddress string               `json:"complete_address" validate:"required"`
	GoogleMapLink   string               `json:"google_map_link" validate:"required"`
	Distance        float32              `json:"distance" validate:"required"`
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
