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
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}
}
