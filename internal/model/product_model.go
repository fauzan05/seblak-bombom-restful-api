package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type ProductResponse struct {
	ID          uint64                    `json:"id"`
	Category    CategoryResponse          `json:"category"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Price       float32                   `json:"price"`
	Stock       int                       `json:"stock"`
	Images      []ImageResponse           `json:"images"`
	Reviews     []ProductReviewResponse   `json:"product_reviews"`
	IsActive    bool                      `json:"is_active"`
	CreatedAt   helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt   helper_others.TimeRFC3339 `json:"updated_at"`
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
