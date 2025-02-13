package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type MidtransCoreAPIOrderRepository struct {
	Repository[entity.MidtransCoreAPIOrder]
	Log *logrus.Logger
}

func NewMidtransCoreAPIOrderRepository(log *logrus.Logger) *MidtransCoreAPIOrderRepository {
	return &MidtransCoreAPIOrderRepository{
		Log: log,
	}
}