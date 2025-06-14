package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type WalletWithdrawRequestRepository struct {
	Repository[entity.WalletWithdrawRequests]
	Log *logrus.Logger
}

func NewWalletWithdrawRequestRepository(log *logrus.Logger) *WalletWithdrawRequestRepository {
	return &WalletWithdrawRequestRepository{
		Log: log,
	}
}
