package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type OrderProductRepository struct {
	Repository[entity.OrderProduct]
	Log *logrus.Logger
}

func NewOrderProductRepository(log *logrus.Logger) *OrderProductRepository {
	return &OrderProductRepository{
		Log: log,
	}
}