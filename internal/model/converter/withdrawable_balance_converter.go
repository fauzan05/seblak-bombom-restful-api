package converter

import "seblak-bombom-restful-api/internal/model"

func WithdrawableBalanceResponse(withdrawableBalance *float32, totalWalletBalance *float32) *model.GetWithdrawableBalanceResponse {
	return &model.GetWithdrawableBalanceResponse{
		WithdrawableBalance: *withdrawableBalance,
		TotalWalletBalance:  *totalWalletBalance,
	}
}
