package model

import "time"

type CartItemResponse struct {
	ID        uint64          `json:"id,omitempty"`
	CartId    uint64          `json:"cart_id,omitempty"`
	Product   ProductResponse `json:"product,omitempty"`
	Quantity  int             `json:"quantity,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
}
