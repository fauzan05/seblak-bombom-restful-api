package model

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type DiscountCouponResponse struct {
	ID              uint64                    `json:"id"`
	Name            string                    `json:"name"`
	Description     string                    `json:"description"`
	Code            string                    `json:"code"`
	Value           float32                   `json:"value"`
	Type            enum_state.DiscountType   `json:"type"`
	Start           helper_others.TimeRFC3339 `json:"start"`
	End             helper_others.TimeRFC3339 `json:"end"`
	Status          bool                      `json:"status"`
	MaxUsagePerUser int                       `json:"max_usage_per_user"`
	UsedCount       int                       `json:"used_count"`
	MinOrderValue   float32                   `json:"min_order_value"`
	CreatedAt       helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt       helper_others.TimeRFC3339 `json:"updated_at"`
}

type CreateDiscountCouponRequest struct {
	Name            string                    `json:"name" validate:"required,max=100"`
	Description     string                    `json:"description" validate:"required"`
	Code            string                    `json:"code" validate:"required,max=100"`
	Value           float32                   `json:"value" validate:"required"`
	Type            enum_state.DiscountType   `json:"type" validate:"required"`
	Start           helper_others.TimeRFC3339 `json:"start" validate:"required"`
	End             helper_others.TimeRFC3339 `json:"end" validate:"required"`
	MaxUsagePerUser int                       `json:"max_usage_per_user" validate:"required"`
	UsedCount       int                       `json:"used_count"`
	MinOrderValue   float32                   `json:"min_order_value"`
	Status          bool                      `json:"status"`
}

type GetDiscountCouponRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateDiscountCouponRequest struct {
	ID              uint64                    `json:"-" validate:"required"`
	Name            string                    `json:"name" validate:"required,max=100"`
	Description     string                    `json:"description" validate:"required"`
	Code            string                    `json:"code" validate:"required,max=100"`
	Value           float32                   `json:"value" validate:"required"`
	Type            enum_state.DiscountType   `json:"type" validate:"required"`
	Start           helper_others.TimeRFC3339 `json:"start" validate:"required"`
	End             helper_others.TimeRFC3339 `json:"end" validate:"required"`
	MaxUsagePerUser int                       `json:"max_usage_per_user" validate:"required"`
	UsedCount       int                       `json:"used_count"`
	MinOrderValue   float32                   `json:"min_order_value"`
	Status          bool                      `json:"status"`
}

type DeleteDiscountCouponRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}
