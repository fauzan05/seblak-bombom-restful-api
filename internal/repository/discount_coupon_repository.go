package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type DiscountCouponRepository struct {
	Repository[entity.DiscountCoupon]
	Log *logrus.Logger
}

func NewDiscountRepository(log *logrus.Logger) *DiscountCouponRepository {
	return &DiscountCouponRepository{
		Log: log,
	}
}