package model

import "seblak-bombom-restful-api/internal/helper"

type OrderResponse struct {
	ID                 uint64                `json:"id,omitempty"`
	Invoice            string                `json:"invoice,omitempty"`
	ProductId          uint64                `json:"product_id,omitempty"`
	ProductName        string                `json:"product_name,omitempty"`
	ProductDescription string                `json:"product_description,omitempty"`
	Price              int                   `json:"price,omitempty"`
	Quantity           int                   `json:"quantity,omitempty"`
	Amount             int                   `json:"amount,omitempty"`
	DiscountValue      int                   `json:"discount_value,omitempty"`
	DiscountType       helper.DiscountType   `json:"discount_type,omitempty"`
	UserId             uint64                `json:"user_id,omitempty"`
	FirstName          string                `json:"first_name,omitempty"`
	LastName           string                `json:"last_name,omitempty"`
	Email              string                `json:"email,omitempty"`
	Phone              string                `json:"phone,omitempty"`
	PaymentMethod      helper.PaymentMethod  `json:"payment_method,omitempty"`
	PaymentStatus      helper.PaymentStatus  `json:"payment_status,omitempty"`
	DeliveryStatus     helper.DeliveryStatus `json:"delivery_status,omitempty"`
	IsDelivery         bool                  `json:"delivery,omitempty"`
	DeliveryCost       int                   `json:"delivery_cost,omitempty"`
	CategoryName       string                `json:"category_name,omitempty"`
	CompleteAddress    string                `json:"complete_address,omitempty"`
	GoogleMapLink      string                `json:"google_map_link,omitempty"`
	CreatedAt          string                `json:"created_at,omitempty"`
	UpdatedAt          string                `json:"updated_at,omitempty"`
}

type CreateOrderRequest struct {
	ProductId          uint64                `json:"product_id" validate:"required"`
	ProductName        string                `json:"product_name" validate:"required"`
	ProductDescription string                `json:"product_description" validate:"required"`
	Price              int                   `json:"price" validate:"required"`
	Quantity           int                   `json:"quantity" validate:"required"`
	Amount             int                   `json:"amount" validate:"required"`
	DiscountCode       string                `json:"discount_code"`
	DiscountValue      int                   `json:"discount_value"`
	DiscountType       helper.DiscountType   `json:"discount_type"`
	UserId             uint64                `json:"user_id" validate:"required"`
	FirstName          string                `json:"first_name" validate:"required"`
	LastName           string                `json:"last_name" validate:"required"`
	Email              string                `json:"email" validate:"required"`
	Phone              string                `json:"phone" validate:"required"`
	PaymentMethod      helper.PaymentMethod  `json:"payment_method" validate:"required"`
	PaymentStatus      helper.PaymentStatus  `json:"payment_status" validate:"required"`
	DeliveryStatus     helper.DeliveryStatus `json:"delivery_status" validate:"required"`
	IsDelivery         bool                  `json:"is_delivery" validate:"required"`
	DeliveryCost       int                   `json:"delivery_cost" validate:"required"`
	CategoryName       string                `json:"category_name" validate:"required"`
	CompleteAddress    string                `json:"complete_address" validate:"required"`
	GoogleMapLink      string                `json:"google_map_link" validate:"required"`
	Distance           int                   `json:"distance" validate:"required"`
}
