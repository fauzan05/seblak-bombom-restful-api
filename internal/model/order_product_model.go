package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type OrderProductResponse struct {
	ID                        uint64                    `json:"id,omitempty"`
	OrderId                   uint64                    `json:"order_id,omitempty"`
	ProductId                 uint64                    `json:"product_id,omitempty"`
	ProductName               string                    `json:"product_name,omitempty"`
	ProductFirstImagePosition string                    `json:"product_first_image_position"`
	Category                  string                    `json:"category,omitempty"`
	Price                     float32                   `json:"price,omitempty"`
	Quantity                  int                       `json:"quantity,omitempty"`
	Product                   ProductResponse           `json:"product,omitempty"`
	CreatedAt                 helper_others.TimeRFC3339 `json:"created_at,omitempty"`
	UpdatedAt                 helper_others.TimeRFC3339 `json:"updated_at,omitempty"`
}
