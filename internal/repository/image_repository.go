package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type ImageRepository struct {
	Repository[entity.Image]
	Log *logrus.Logger
}

func NewImageRepository(log *logrus.Logger) *ImageRepository {
	return &ImageRepository{
		Log: log,
	}
}
