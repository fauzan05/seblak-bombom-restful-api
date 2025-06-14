package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func WalletWithdrawToResponse(walletWithdrawRequest *entity.WalletWithdrawRequests) *model.WithdrawWalletResponse {
	response := &model.WithdrawWalletResponse{
		ID:                walletWithdrawRequest.ID,
		UserId:            walletWithdrawRequest.UserId,
		User:              *UserToResponse(walletWithdrawRequest.User),
		Amount:            walletWithdrawRequest.Amount,
		Method:            walletWithdrawRequest.Method,
		BankName:          walletWithdrawRequest.BankName,
		BankAccountNumber: walletWithdrawRequest.BankAcountNumber,
		BankAccountName:   walletWithdrawRequest.BankAcountName,
		Status:            walletWithdrawRequest.Status,
		Note:              walletWithdrawRequest.Note,
		RejectionNotes:    walletWithdrawRequest.RejectionNotes,
		CreatedAt:         helper_others.TimeRFC3339(walletWithdrawRequest.CreatedAt),
		UpdatedAt:         helper_others.TimeRFC3339(walletWithdrawRequest.UpdatedAt),
	}

	if walletWithdrawRequest.ProcessedBy != nil {
		response.ProcessedBy = *UserToResponse(walletWithdrawRequest.User)
	}

	if walletWithdrawRequest.ProcessedAt != nil {
		response.ProcessedAt = helper_others.TimeRFC3339(*walletWithdrawRequest.ProcessedAt)
	}

	return response
}

func WalletWithdrawsToResponse(walletWithdrawRequest *[]entity.WalletWithdrawRequests) *[]model.WithdrawWalletResponse {
	getWalletWithdrawRequest := make([]model.WithdrawWalletResponse, len(*walletWithdrawRequest))
	for i, walletWithdrawRequest := range *walletWithdrawRequest {
		getWalletWithdrawRequest[i] = *WalletWithdrawToResponse(&walletWithdrawRequest)
	}
	return &getWalletWithdrawRequest
}
