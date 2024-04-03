package model

import (
	"seblak-bombom-restful-api/internal/helper"
)

type DiscountResponse struct {
	ID          uint64              `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Code        string              `json:"code,omitempty"`
	Value       float32             `json:"value,omitempty"`
	Type        helper.DiscountType `json:"type,omitempty"`
	Start       string              `json:"start,omitempty"`
	End         string              `json:"end,omitempty"`
	Status      bool                `json:"status,omitempty"`
	CreatedAt   string              `json:"created_at,omitempty"`
	UpdatedAt   string              `json:"updated_at,omitempty"`
}

type CreateDiscountRequest struct {
	Name        string              `json:"name" validate:"required,max=100"`
	Description string              `json:"description" validate:"required"`
	Code        string              `json:"code" validate:"required,max=100"`
	Value       float32             `json:"value" validate:"required"`
	Type        helper.DiscountType `json:"type" validate:"required"`
	Start       string              `json:"start" validate:"required"`
	End         string              `json:"end" validate:"required"`
	Status      bool                `json:"status" validate:"required"`
}

type GetDiscountRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateDiscountRequest struct {
	ID          uint64              `json:"-" validate:"required"`
	Name        string              `json:"name" validate:"required,max=100"`
	Description string              `json:"description" validate:"required"`
	Code        string              `json:"code" validate:"required,max=100"`
	Value       float32             `json:"value" validate:"required"`
	Type        helper.DiscountType `json:"type" validate:"required"`
	Start       string              `json:"start" validate:"required"`
	End         string              `json:"end" validate:"required"`
	Status      bool                `json:"status" validate:"required"`
}

type DeleteDiscountRequest struct {
	ID uint64 `json:"-" validate:"required"`
}
