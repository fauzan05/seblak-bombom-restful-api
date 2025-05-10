package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type PasswordResetRepository struct {
	Repository[entity.PasswordReset]
	Log *logrus.Logger
}

func NewPasswordResetRepository(log *logrus.Logger) *PasswordResetRepository {
	return &PasswordResetRepository{
		Log: log,
	}
}
