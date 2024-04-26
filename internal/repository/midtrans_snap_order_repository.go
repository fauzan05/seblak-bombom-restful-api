package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type MidtransSnapOrderRepository struct {
	Repository[entity.MidtransSnapOrder]
	Log *logrus.Logger
}

func NewMidtransSnapOrderRepository(log *logrus.Logger) *MidtransSnapOrderRepository {
	return &MidtransSnapOrderRepository{
		Log: log,
	}
}