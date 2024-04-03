package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type DeliveryRepository struct {
	Repository[entity.Delivery]
	Log *logrus.Logger
}

func NewDeliveryRepository(log *logrus.Logger) *DeliveryRepository {
	return &DeliveryRepository{
		Log: log,
	}
}
