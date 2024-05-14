package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type CartRepository struct {
	Repository[entity.Cart]
	Log *logrus.Logger
}

func NewCartRepository(log *logrus.Logger) *CartRepository {
	return &CartRepository{
		Log: log,
	}
}