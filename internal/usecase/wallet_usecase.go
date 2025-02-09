package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletUseCase struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	UserRepository    *repository.UserRepository
}

func NewWalletUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository) *WalletUseCase {
	return &WalletUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		UserRepository:    userRepository,
	}
}

func (c *WalletUseCase) AddBalance(ctx context.Context, request *model.TopUpWalletBalance) (model.WalletResponse, error) {
	newWallet := new(entity.Wallet)
	return *converter.WalletToResponse(newWallet), nil
}