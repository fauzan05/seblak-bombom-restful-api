package repository

import (
	"seblak-bombom-restful-api/internal/entity"

	"github.com/sirupsen/logrus"
)

type TokenRepository struct {
	Repository[entity.Token]
	Log       *logrus.Logger
}

func NewTokenRepository(log *logrus.Logger) *TokenRepository {
	return &TokenRepository{
		Log: log,
	}
}