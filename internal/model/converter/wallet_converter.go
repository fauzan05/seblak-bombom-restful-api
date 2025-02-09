package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func WalletToResponse(wallet *entity.Wallet) *model.WalletResponse {
	return &model.WalletResponse{
		ID:        wallet.ID,
		Balance:   wallet.Balance,
		Status:    wallet.Status,
		CreatedAt: wallet.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: wallet.Updated_At.Format("2006-01-02 15:04:05"),
	}
}
