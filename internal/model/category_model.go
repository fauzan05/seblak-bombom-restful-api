package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type CategoryResponse struct {
	ID          uint64                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	CreatedAt   helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt   helper_others.TimeRFC3339 `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
}

type GetCategoryRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateCategoryRequest struct {
	ID          uint64 `json:"-" validate:"required"`
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
}

type DeleteCategoryRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}
