package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type DiscountUsageRepository struct {
	Repository[entity.DiscountUsage]
	Log *logrus.Logger
}

func NewDiscountUsageRepository(log *logrus.Logger) *DiscountUsageRepository {
	return &DiscountUsageRepository{
		Log: log,
	}
}
