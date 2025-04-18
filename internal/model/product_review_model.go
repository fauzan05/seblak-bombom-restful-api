package model

import (
	"seblak-bombom-restful-api/internal/helper"
)

type ProductReviewResponse struct {
	ID        uint64             `json:"id"`
	ProductId uint64             `json:"product_id"`
	UserId    uint64             `json:"user_id"`
	Rate      int                `json:"rate"`
	Comment   string             `json:"comment"`
	CreatedAt helper.TimeRFC3339 `json:"created_at"`
	UpdatedAt helper.TimeRFC3339 `json:"updated_at"`
}

type CreateProductReviewRequest struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	UserId    uint64 `json:"-" validate:"required"` // via token
	Rate      int    `json:"rate" validate:"required"`
	Comment   string `json:"comment" validate:"required"`
}
