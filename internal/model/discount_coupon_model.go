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
	Start           string              `json:"start"`
	End             string              `json:"end"`
	Status          bool                `json:"status"`
	TotalMaxUsage   int                 `json:"total_max_usage"`
	MaxUsagePerUser int                 `json:"max_usage_per_user"`
	UsedCount       int                 `json:"used_count"`
	MinOrderValue   int                 `json:"min_order_value"`
	CreatedAt       string              `json:"created_at"`
	UpdatedAt       string              `json:"updated_at"`
}

type CreateDiscountCouponRequest struct {
	Name            string              `json:"name" validate:"required,max=100"`
	Description     string              `json:"description" validate:"required"`
	Code            string              `json:"code" validate:"required,max=100"`
	Value           float32             `json:"value" validate:"required"`
	Type            helper.DiscountType `json:"type" validate:"required"`
	Start           string              `json:"start" validate:"required"`
	End             string              `json:"end" validate:"required"`
	TotalMaxUsage   int                 `json:"total_max_usage" validate:"required"`
	MaxUsagePerUser int                 `json:"max_usage_per_user" validate:"required"`
	UsedCount       int                 `json:"used_count"`
	MinOrderValue   int                 `json:"min_order_value"`
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
	Start           string              `json:"start" validate:"required"`
	End             string              `json:"end" validate:"required"`
	TotalMaxUsage   int                 `json:"total_max_usage" validate:"required"`
	MaxUsagePerUser int                 `json:"max_usage_per_user" validate:"required"`
	UsedCount       int                 `json:"used_count"`
	MinOrderValue   int                 `json:"min_order_value"`
	Status          bool                `json:"status"`
}

type DeleteDiscountCouponRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}
