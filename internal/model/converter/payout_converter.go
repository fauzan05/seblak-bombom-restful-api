package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func PayoutToResponse(payout *entity.Payout) *model.PayoutResponse {
	response := &model.PayoutResponse{
		ID:                payout.ID,
		XenditPayoutId:    payout.XenditPayoutId,
		Amount:            payout.Amount,
		Currency:          payout.Currency,
		Method:            payout.Method,
		Status:            payout.Status,
		Notes:             payout.Notes,
		CancellationNotes: payout.CancellationNotes,
		CreatedAt:         payout.Created_At,
		UpdatedAt:         payout.Updated_At,
	}

	if payout.XenditPayout != nil {
		response.XenditPayout = XenditPayoutToResponse(payout.XenditPayout)
	}

	return response
}

func PayoutsToResponse(payouts *[]entity.Payout) *[]model.PayoutResponse {
	getPayouts := make([]model.PayoutResponse, len(*payouts))
	for i, payout := range *payouts {
		getPayouts[i] = *PayoutToResponse(&payout)
	}
	return &getPayouts
}
