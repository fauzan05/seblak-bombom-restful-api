package model

import (
	"mime/multipart"
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type CategoryResponse struct {
	ID            uint64                    `json:"id"`
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	ImageFilename string                    `json:"image_filename"`
	IsActive      bool                      `json:"is_active"`
	CreatedAt     helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt     helper_others.TimeRFC3339 `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name        string                `json:"name" validate:"required,max=100"`
	Description string                `json:"description"`
	Image       *multipart.FileHeader `json:"image"`
}

type GetCategoryRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateCategoryRequest struct {
	ID          uint64                `json:"-" validate:"required"`
	Name        string                `json:"name" validate:"required,max=100"`
	Description string                `json:"description"`
	Image       *multipart.FileHeader `json:"image"`
}

type DeleteCategoryRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}
