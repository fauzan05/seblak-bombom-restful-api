package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type CartItemRepository struct {
	Repository[entity.CartItem]
	Log *logrus.Logger
}

func NewCartItemRepository(log *logrus.Logger) *CartItemRepository {
	return &CartItemRepository{
		Log: log,
	}
}