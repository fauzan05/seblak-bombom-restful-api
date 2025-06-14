package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type CartItemResponse struct {
	ID        uint64                    `json:"id,omitempty"`
	CartId    uint64                    `json:"cart_id,omitempty"`
	Product   ProductResponse           `json:"product,omitempty"`
	Quantity  int                       `json:"quantity,omitempty"`
	CreatedAt helper_others.TimeRFC3339 `json:"created_at,omitempty"`
	UpdatedAt helper_others.TimeRFC3339 `json:"updated_at,omitempty"`
}
