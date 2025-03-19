package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type PayoutRepository struct {
	Repository[entity.Payout]
	Log *logrus.Logger
}

func NewPayoutRepository(log *logrus.Logger) *PayoutRepository {
	return &PayoutRepository{
		Log: log,
	}
}
