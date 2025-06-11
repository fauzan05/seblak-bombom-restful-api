package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type CartResponse struct {
	ID        uint64                    `json:"id,omitempty"`
	UserID    uint64                    `json:"user_id,omitempty"`
	CartItems []CartItemResponse        `json:"cart_items"`
	CreatedAt helper_others.TimeRFC3339 `json:"created_at,omitempty"`
	UpdatedAt helper_others.TimeRFC3339 `json:"updated_at,omitempty"`
}

type CreateCartRequest struct {
	UserID    uint64 `json:"user_id" validate:"required"`
	ProductID uint64 `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
}

type GetAllCartByCurrentUserRequest struct {
	UserID uint64 `json:"-" validate:"required"`
}

type UpdateCartRequest struct {
	CartItemID uint64 `json:"-" validate:"required"`
	Quantity   int    `json:"quantity" validate:"required"`
}

type DeleteCartRequest struct {
	CartItemID uint64 `json:"-" validate:"required"`
}
