package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type ProductResponse struct {
	ID          uint64                    `json:"id,omitempty"`
	Category    CategoryResponse          `json:"category,omitempty"`
	Name        string                    `json:"name,omitempty"`
	Description string                    `json:"description,omitempty"`
	Price       float32                   `json:"price,omitempty"`
	Stock       int                       `json:"stock,omitempty"`
	Images      []ImageResponse           `json:"images,omitempty"`
	Reviews     []ProductReviewResponse   `json:"product_reviews,omitempty"`
	CreatedAt   helper_others.TimeRFC3339 `json:"created_at,omitempty"`
	UpdatedAt   helper_others.TimeRFC3339 `json:"updated_at,omitempty"`
}

type CreateProductRequest struct {
	CategoryId  uint64  `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required,max=100"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
}

type GetProductRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateProductRequest struct {
	ID          uint64  `json:"-" validate:"required"`
	CategoryId  uint64  `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required,max=100"`
	Description string  `json:"description" validate:"required"`
	Price       float32 `json:"price"`
	Stock       int     `json:"stock"`
}

type DeleteProductRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}
