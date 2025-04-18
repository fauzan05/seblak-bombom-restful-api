package model

import (
	"seblak-bombom-restful-api/internal/helper"
)

type MidtransSnapOrderResponse struct {
	ID          uint64             `json:"id"`
	OrderId     uint64             `json:"order_id"`
	Token       string             `json:"token"`
	RedirectUrl string             `json:"redirect_url"`
	CreatedAt   helper.TimeRFC3339 `json:"created_at"`
	UpdatedAt   helper.TimeRFC3339 `json:"updated_at"`
}

type CreateMidtransSnapOrderRequest struct {
	OrderId uint64 `json:"order_id" validate:"required"`
}

type GetMidtransSnapOrderRequest struct {
	OrderId uint64 `json:"order_id" validate:"required"`
}
