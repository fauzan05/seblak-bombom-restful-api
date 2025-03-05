package model

import "time"

type MidtransSnapOrderResponse struct {
	ID          uint64    `json:"id"`
	OrderId     uint64    `json:"order_id"`
	Token       string    `json:"token"`
	RedirectUrl string    `json:"redirect_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateMidtransSnapOrderRequest struct {
	OrderId uint64 `json:"order_id" validate:"required"`
}

type GetMidtransSnapOrderRequest struct {
	OrderId uint64 `json:"order_id" validate:"required"`
}
