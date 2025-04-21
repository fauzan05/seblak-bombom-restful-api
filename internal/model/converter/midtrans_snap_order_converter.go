package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
)

func MidtransSnapOrderToResponse(midtransSnapOrder *entity.MidtransSnapOrder) *model.MidtransSnapOrderResponse {
	return &model.MidtransSnapOrderResponse{
		ID:          midtransSnapOrder.ID,
		Token:       midtransSnapOrder.Token,
		RedirectUrl: midtransSnapOrder.RedirectUrl,
		CreatedAt:   helper.TimeRFC3339(midtransSnapOrder.CreatedAt),
		UpdatedAt:   helper.TimeRFC3339(midtransSnapOrder.UpdatedAt),
	}
}
