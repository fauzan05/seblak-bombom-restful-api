package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
)

func PayoutToResponse(payout *entity.Payout) *model.PayoutResponse {
	response := &model.PayoutResponse{
		ID:             payout.ID,
		XenditPayoutId: helper.NullStringToString(payout.XenditPayoutId),
		Amount:         payout.Amount,
		Currency:       payout.Currency,
		Method:         payout.Method,
		Status:         payout.Status,
		Notes:          payout.Notes,
		CreatedAt:      payout.CreatedAt,
		UpdatedAt:      payout.UpdatedAt,
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
