package model

import (
	"seblak-bombom-restful-api/internal/helper"
)

type DiscountCouponResponse struct {
	ID              uint64              `json:"id"`
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Code            string              `json:"code"`
	Value           float32             `json:"value"`
	Type            helper.DiscountType `json:"type"`
	Start           helper.TimeRFC3339  `json:"start"`
	End             helper.TimeRFC3339  `json:"end"`
	Status          bool                `json:"status"`
	MaxUsagePerUser int                 `json:"max_usage_per_user"`
	UsedCount       int                 `json:"used_count"`
	MinOrderValue   float32             `json:"min_order_value"`
	CreatedAt       helper.TimeRFC3339  `json:"created_at"`
	UpdatedAt       helper.TimeRFC3339  `json:"updated_at"`
}

type CreateDiscountCouponRequest struct {
	Name            string              `json:"name" validate:"required,max=100"`
	Description     string              `json:"description" validate:"required"`
	Code            string              `json:"code" validate:"required,max=100"`
	Value           float32             `json:"value" validate:"required"`
	Type            helper.DiscountType `json:"type" validate:"required"`
	Start           helper.TimeRFC3339  `json:"start" validate:"required"`
	End             helper.TimeRFC3339  `json:"end" validate:"required"`
	MaxUsagePerUser int                 `json:"max_usage_per_user" validate:"required"`
	UsedCount       int                 `json:"used_count"`
	MinOrderValue   float32             `json:"min_order_value"`
	Status          bool                `json:"status"`
}

type GetDiscountCouponRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateDiscountCouponRequest struct {
	ID              uint64              `json:"-" validate:"required"`
	Name            string              `json:"name" validate:"required,max=100"`
	Description     string              `json:"description" validate:"required"`
	Code            string              `json:"code" validate:"required,max=100"`
	Value           float32             `json:"value" validate:"required"`
	Type            helper.DiscountType `json:"type" validate:"required"`
	Start           helper.TimeRFC3339  `json:"start" validate:"required"`
	End             helper.TimeRFC3339  `json:"end" validate:"required"`
	MaxUsagePerUser int                 `json:"max_usage_per_user" validate:"required"`
	UsedCount       int                 `json:"used_count"`
	MinOrderValue   float32             `json:"min_order_value"`
	Status          bool                `json:"status"`
}

type DeleteDiscountCouponRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}
