package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	UserRepository   *repository.UserRepository
	WalletRepository *repository.WalletRepository
}

func NewWalletUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, walletRepository *repository.WalletRepository) *WalletUseCase {
	return &WalletUseCase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		UserRepository:   userRepository,
		WalletRepository: walletRepository,
	}
}

func (c *WalletUseCase) Withdraw(ctx context.Context, request *model.WithDrawWalletRequest) (*model.WalletResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newUser := new(entity.User)
	newUser.ID = request.UserId
	if err := c.UserRepository.FindWithPreloads(tx, newUser, "Wallet"); err != nil {
		c.Log.Warnf("failed to get user by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get user by id : %+v", err))
	}

	if request.Amount > newUser.Wallet.Balance {
		c.Log.Warnf("your balance is insufficient to perform this withdraw transaction!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "your balance is insufficient to perform this withdraw transaction!")
	}

	newWallet := new(entity.Wallet)
	newWallet.ID = newUser.Wallet.ID
	newWallet.UserId = newUser.ID
	newWallet.Balance = newUser.Wallet.Balance - request.Amount
	newWallet.Status = newUser.Wallet.Status
	if err := c.WalletRepository.Update(tx, newWallet); err != nil {
		c.Log.Warnf("failed to update wallet balance : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update wallet balance : %+v", err))
	}

	err := helper_others.SaveWalletTransaction(tx, newWallet.UserId, nil, request.Amount, enum_state.WITHDRAW, enum_state.CASH, request.Notes)
	if err != nil {
		c.Log.Warnf("failed to save wallet transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save wallet transaction : %+v", err))
	}

	if err := c.WalletRepository.FindFirst(tx, newWallet); err != nil {
		c.Log.Warnf("failed to get wallet by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to get wallet by id : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.WalletToResponse(newWallet), nil
}
