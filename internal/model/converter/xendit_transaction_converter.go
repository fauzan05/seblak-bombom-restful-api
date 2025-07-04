package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
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
		FailureCode:     xenditTransaction.FailureCode,
		Metadata:        xenditTransaction.Metadata,
		ExpiresAt:       helper_others.TimeRFC3339(xenditTransaction.ExpiresAt),
		CreatedAt:       helper_others.TimeRFC3339(xenditTransaction.CreatedAt),
		UpdatedAt:       helper_others.TimeRFC3339(xenditTransaction.UpdatedAt),
	}

	return response
}
