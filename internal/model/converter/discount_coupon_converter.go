package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
)

func DiscountCouponToResponse(discount *entity.DiscountCoupon) *model.DiscountCouponResponse {
	return &model.DiscountCouponResponse{
		ID:              discount.ID,
		Name:            discount.Name,
		Description:     discount.Description,
		Code:            discount.Code,
		Value:           discount.Value,
		Type:            discount.Type,
		Start:           helper.TimeRFC3339(discount.Start),
		End:             helper.TimeRFC3339(discount.End),
		Status:          discount.Status,
		MaxUsagePerUser: discount.MaxUsagePerUser,
		UsedCount:       discount.UsedCount,
		MinOrderValue:   discount.MinOrderValue,
		CreatedAt:       helper.TimeRFC3339(discount.CreatedAt),
		UpdatedAt:       helper.TimeRFC3339(discount.UpdatedAt),
	}
}

func DiscountCouponsToResponse(discounts *[]entity.DiscountCoupon) *[]model.DiscountCouponResponse {
	getDiscounts := make([]model.DiscountCouponResponse, len(*discounts))
	for i, discount := range *discounts {
		getDiscounts[i] = *DiscountCouponToResponse(&discount)
	}
	return &getDiscounts
}
