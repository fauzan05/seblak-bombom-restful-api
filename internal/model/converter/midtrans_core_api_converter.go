package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func MidtransCoreAPIToResponse(midtransCoreAPIOrder *entity.MidtransCoreAPIOrder) *model.MidtransCoreAPIOrderResponse {
	return &model.MidtransCoreAPIOrderResponse{
		ID: midtransCoreAPIOrder.ID,
		StatusCode: midtransCoreAPIOrder.StatusCode,
		StatusMessage: midtransCoreAPIOrder.StatusMessage,
		TransactionId: midtransCoreAPIOrder.TransactionId,
		OrderId: midtransCoreAPIOrder.OrderId,
		MidtransOrderId: midtransCoreAPIOrder.MidtransOrderId,
		GrossAmount: midtransCoreAPIOrder.GrossAmount,
		Currency: midtransCoreAPIOrder.Currency,
		PaymentType: midtransCoreAPIOrder.PaymentType,
		ExpiryTime: midtransCoreAPIOrder.ExpiryTime.Format("2006-01-02 15:04:05"),
		TransactionTime: midtransCoreAPIOrder.TransactionTime.Format("2006-01-02 15:04:05"),
		TransactionStatus: midtransCoreAPIOrder.TransactionStatus,
		FraudStatus: midtransCoreAPIOrder.FraudStatus,
		Actions: MidtransActionsToResponse(&midtransCoreAPIOrder.Actions),
		CreatedAt: midtransCoreAPIOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: midtransCoreAPIOrder.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
