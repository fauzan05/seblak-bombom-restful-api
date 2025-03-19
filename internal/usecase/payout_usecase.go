package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PayoutUseCase struct {
	DB               *gorm.DB
	Log              *logrus.Logger
	Validate         *validator.Validate
	PayoutRepository *repository.PayoutRepository
}

func NewPayoutUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	payoutRepository *repository.PayoutRepository) *PayoutUseCase {
	return &PayoutUseCase{
		DB:               db,
		Log:              log,
		Validate:         validate,
		PayoutRepository: payoutRepository,
	}
}

func (c *PayoutUseCase) Add(ctx context.Context, request *model.CreatePayoutRequest) (*model.PayoutResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	// xenditPayoutId := ""
	if request.Method == helper.PAYOUT_METHOD_OFFLINE {
		fmt.Println("DATANYA : ", request.XenditPayoutRequest)
	}

	newPayout := new(entity.Payout)
	newPayout.UserId = request.UserId
	if err := c.PayoutRepository.Create(tx, newPayout); err != nil {
		c.Log.Warnf("Failed to create payout request into database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to create payout request into database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.PayoutToResponse(newPayout), nil
}