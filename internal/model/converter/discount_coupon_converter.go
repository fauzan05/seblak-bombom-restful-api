package converter

import (
	"seblak-bombom-restful-api/internal/entity"
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
		Start:           discount.Start,
		End:             discount.End,
		Status:          discount.Status,
		MaxUsagePerUser: discount.MaxUsagePerUser,
		TotalMaxUsage:   discount.TotalMaxUsage,
		UsedCount:       discount.UsedCount,
		MinOrderValue:   discount.MinOrderValue,
		CreatedAt:       discount.Created_At,
		UpdatedAt:       discount.Updated_At,
	}
}

func DiscountCouponsToResponse(discounts *[]entity.DiscountCoupon) *[]model.DiscountCouponResponse {
	getDiscounts := make([]model.DiscountCouponResponse, len(*discounts))
	for i, discount := range *discounts {
		getDiscounts[i] = *DiscountCouponToResponse(&discount)
	}
	return &getDiscounts
}
