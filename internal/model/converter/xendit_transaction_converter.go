package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func XenditTransactionToResponse(xenditTransaction entity.XenditTransactions) *model.XenditTransactionResponse {
	response := &model.XenditTransactionResponse{
		ID:              xenditTransaction.ID,
		ReferenceId:     xenditTransaction.ReferenceId,
		OrderId:         xenditTransaction.OrderId,
		Amount:          xenditTransaction.Amount,
		Currency:        xenditTransaction.Currency,
		PaymentMethod:   xenditTransaction.PaymentMethod,
		PaymentMethodId: xenditTransaction.PaymentMethodId,
		ChannelCode:     xenditTransaction.ChannelCode,
		QrString:        xenditTransaction.QrString,
		Status:          xenditTransaction.Status,
		Description:     xenditTransaction.Description,
		ExpiresAt:       xenditTransaction.ExpiresAt,
		CreatedAt:       xenditTransaction.Created_At,
		UpdatedAt:       xenditTransaction.Updated_At,
	}

	return response
}
