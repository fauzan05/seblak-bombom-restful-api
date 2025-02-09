package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type WalletRepository struct {
	Repository[entity.Wallet]
	Log *logrus.Logger
}

func NewWalletRepository(log *logrus.Logger) *WalletRepository {
	return &WalletRepository{
		Log: log,
	}
}
