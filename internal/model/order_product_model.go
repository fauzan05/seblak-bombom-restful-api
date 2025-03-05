package model

import "time"

type OrderProductResponse struct {
	ID          uint64    `json:"id,omitempty"`
	OrderId     uint64    `json:"order_id,omitempty"`
	ProductId   uint64    `json:"product_id,omitempty"`
	ProductName string    `json:"product_name,omitempty"`
	Price       string    `json:"price,omitempty"`
	Quantity    int       `json:"quantity,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type CreateOrderProductRequest struct {
	OrderId     uint64  `json:"order_id"`
	ProductId   uint64  `json:"product_id" validate:"required"`
	ProductName string  `json:"product_name" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}
