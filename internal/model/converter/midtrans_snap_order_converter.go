package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func MidtransSnapOrderToResponse(midtransSnapOrder *entity.MidtransSnapOrder) *model.MidtransSnapOrderResponse {
	return &model.MidtransSnapOrderResponse{
		ID: midtransSnapOrder.ID,
		Token: midtransSnapOrder.Token,
		RedirectUrl: midtransSnapOrder.RedirectUrl,
		CreatedAt: midtransSnapOrder.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: midtransSnapOrder.Updated_At.Format("2006-01-02 15:04:05"),
	}
}