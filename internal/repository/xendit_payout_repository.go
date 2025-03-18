package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type XenditPayoutRepository struct {
	Repository[entity.XenditPayout]
	Log *logrus.Logger
}

func NewXenditPayoutRepository(log *logrus.Logger) *XenditPayoutRepository {
	return &XenditPayoutRepository{
		Log: log,
	}
}
