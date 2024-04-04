package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type ProductReviewRepository struct {
	Repository[entity.ProductReview]
	Log *logrus.Logger
}

func NewProductReviewRepository(log *logrus.Logger) *ProductReviewRepository {
	return &ProductReviewRepository{
		Log: log,
	}
}
