package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type DiscountRepository struct {
	Repository[entity.Discount]
	Log *logrus.Logger
}

func NewDiscountRepository(log *logrus.Logger) *DiscountRepository {
	return &DiscountRepository{
		Log: log,
	}
}