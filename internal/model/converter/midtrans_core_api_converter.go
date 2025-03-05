package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func MidtransCoreAPIToResponse(midtransCoreAPIOrder *entity.MidtransCoreAPIOrder) *model.MidtransCoreAPIOrderResponse {
	return &model.MidtransCoreAPIOrderResponse{
		ID:                midtransCoreAPIOrder.ID,
		StatusCode:        midtransCoreAPIOrder.StatusCode,
		StatusMessage:     midtransCoreAPIOrder.StatusMessage,
		TransactionId:     midtransCoreAPIOrder.TransactionId,
		OrderId:           midtransCoreAPIOrder.OrderId,
		MidtransOrderId:   midtransCoreAPIOrder.MidtransOrderId,
		GrossAmount:       midtransCoreAPIOrder.GrossAmount,
		Currency:          midtransCoreAPIOrder.Currency,
		PaymentType:       midtransCoreAPIOrder.PaymentType,
		ExpiryTime:        midtransCoreAPIOrder.ExpiryTime,
		TransactionTime:   midtransCoreAPIOrder.TransactionTime,
		TransactionStatus: midtransCoreAPIOrder.TransactionStatus,
		FraudStatus:       midtransCoreAPIOrder.FraudStatus,
		Actions:           MidtransActionsToResponse(&midtransCoreAPIOrder.Actions),
		CreatedAt:         midtransCoreAPIOrder.CreatedAt,
		UpdatedAt:         midtransCoreAPIOrder.UpdatedAt,
	}
}
