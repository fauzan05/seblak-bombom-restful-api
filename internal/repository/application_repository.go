package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type ApplicationRepository struct {
	Repository[entity.Application]
	Log *logrus.Logger
}

func NewApplicationRepository(log *logrus.Logger) *ApplicationRepository {
	return &ApplicationRepository{
		Log: log,
	}
}