package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type XenditTransctionRepository struct {
	Repository[entity.XenditTransactions]
	Log *logrus.Logger
}

func NewXenditTransactionRepository(log *logrus.Logger) *XenditTransctionRepository {
	return &XenditTransctionRepository{
		Log: log,
	}
}
