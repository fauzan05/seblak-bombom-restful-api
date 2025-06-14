package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func XenditPayoutToResponse(xenditPayout *entity.XenditPayout) *model.XenditPayoutResponse {
	response := &model.XenditPayoutResponse{
		ID:                xenditPayout.ID,
		UserId:            xenditPayout.UserID,
		BusinessId:        xenditPayout.BusinessID,
		ReferenceId:       xenditPayout.ReferenceID,
		Amount:            xenditPayout.Amount,
		Currency:          xenditPayout.Currency,
		Description:       xenditPayout.Description,
		ChannelCode:       xenditPayout.ChannelCode,
		AccountNumber:     xenditPayout.AccountNumber,
		AccountHolderName: xenditPayout.AccountHolderName,
		Status:            xenditPayout.Status,
		CreatedAt:         helper_others.TimeRFC3339(xenditPayout.CreatedAt),
		UpdatedAt:         helper_others.TimeRFC3339(xenditPayout.UpdatedAt),
		EstimatedArrival:  helper_others.TimeRFC3339(xenditPayout.EstimatedArrival),
	}

	if xenditPayout.User != nil {
		response.User = UserToResponse(xenditPayout.User)
	}

	return response
}

func XenditPayoutsToResponse(xenditPayouts *[]entity.XenditPayout) *[]model.XenditPayoutResponse {
	getXenditPayouts := make([]model.XenditPayoutResponse, len(*xenditPayouts))
	for i, xenditPayout := range *xenditPayouts {
		getXenditPayouts[i] = *XenditPayoutToResponse(&xenditPayout)
	}
	return &getXenditPayouts
}
