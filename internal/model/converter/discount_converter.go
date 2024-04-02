package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func DiscountToResponse(discount *entity.Discount) *model.DiscountResponse {
	return &model.DiscountResponse{
		ID:          discount.ID,
		Name:        discount.Name,
		Description: discount.Description,
		Code:        discount.Code,
		Value:       discount.Value,
		Type:        discount.Type,
		Start:       discount.Start.Format("2006-01-02 15:04:05"),
		End:         discount.End.Format("2006-01-02 15:04:05"),
		Status:      discount.Status,
		CreatedAt:   discount.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt:   discount.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func DiscountsToResponse(discounts *[]entity.Discount) *[]model.DiscountResponse {
	getDiscounts := make([]model.DiscountResponse, len(*discounts))
	for i, discount := range *discounts {
		getDiscounts[i] = *DiscountToResponse(&discount)
	}
	return &getDiscounts
}