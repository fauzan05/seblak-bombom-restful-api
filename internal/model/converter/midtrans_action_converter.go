package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func MidtransActionToResponse(midtransAction *entity.Action) *model.ActionResponse {
	return &model.ActionResponse{
		Name:   midtransAction.Name,
		Method: midtransAction.Method,
		URL:    midtransAction.URL,
	}
}

func MidtransActionsToResponse(midtransActions *[]entity.Action) *[]model.ActionResponse {
	getActions := make([]model.ActionResponse, len(*midtransActions))
	for i, action := range *midtransActions {
		getActions[i] = *MidtransActionToResponse(&action)
	}
	return &getActions
}