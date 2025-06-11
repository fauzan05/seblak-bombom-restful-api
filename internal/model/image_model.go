package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type ImageResponse struct {
	ID        uint64             `json:"id,omitempty"`
	ProductId uint64             `json:"product_id,omitempty"`
	FileName  string             `json:"file_name,omitempty"`
	Type      string             `json:"type,omitempty"`
	Position  int                `json:"position,omitempty"`
	CreatedAt helper_others.TimeRFC3339 `json:"created_at,omitempty"`
	UpdatedAt helper_others.TimeRFC3339 `json:"updated_at,omitempty"`
}

type AddImagesRequest struct {
	Images []ImageAddRequest `json:"-" validate:"required"`
}

type ImageAddRequest struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	FileName  string `json:"file_name" validate:"required,max=100"`
	Type      string `json:"type" validate:"required"`
	Position  int    `json:"position" validate:"required"`
}

type UpdateImagesRequest struct {
	Images []ImageUpdateRequest `json:"-" validate:"required"`
}

type ImageUpdateRequest struct {
	ID       uint64 `json:"id" validate:"required"`
	Position int    `json:"position" validate:"required"`
}

type DeleteImagesRequest struct {
	Images []DeleteImageRequest `json:"-"`
}

type DeleteImageRequest struct {
	ID uint64 `json:"id" validate:"required"`
}
