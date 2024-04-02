package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DiscountUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	DiscountRepository *repository.DiscountRepository
}

func NewDiscountUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	discountRepository *repository.DiscountRepository) *DiscountUseCase {
	return &DiscountUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		DiscountRepository: discountRepository,
	}
}

var layoutTime string = "2006-01-02 15:04:05"

func (c *DiscountUseCase) Add(ctx context.Context, request *model.CreateDiscountRequest) (*model.DiscountResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount := new(entity.Discount)
	count, err := c.DiscountRepository.CountDiscountByCode(tx, newDiscount, request.Code)
	if err != nil {
		c.Log.Warnf("Failed to count discount by code : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		c.Log.Warnf("Discount code has been used : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount.Name = request.Name
	newDiscount.Description = request.Description
	newDiscount.Code = request.Code
	newDiscount.Value = request.Value
	newDiscount.Type = request.Type
	newDiscount.Start, err = time.Parse(layoutTime, request.Start)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	newDiscount.End, err = time.Parse(layoutTime, request.End)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	newDiscount.Status = request.Status
	if err := c.DiscountRepository.Create(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to create a new discount : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountToResponse(newDiscount), nil
}

func (c *DiscountUseCase) GetAll(ctx context.Context) (*[]model.DiscountResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDiscounts := new([]entity.Discount)
	if err := c.DiscountRepository.FindAll(tx, newDiscounts); err != nil {
		c.Log.Warnf("Failed to find all discounts : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountsToResponse(newDiscounts), nil
}

func (c *DiscountUseCase) GetById(ctx context.Context, request *model.GetDiscountRequest) (*model.DiscountResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDiscount := new(entity.Discount)
	newDiscount.ID = request.ID
	if err := c.DiscountRepository.FindById(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to find discount by id: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountToResponse(newDiscount), nil
}

func (c *DiscountUseCase) Edit(ctx context.Context, request *model.UpdateDiscountRequest) (*model.DiscountResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount := new(entity.Discount)
	newDiscount.ID = request.ID
	if err := c.DiscountRepository.FindById(tx, newDiscount); err != nil {
		c.Log.Warnf("Can't find discount by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	count, err := c.DiscountRepository.CountDiscountByCodeIsExist(tx, newDiscount, newDiscount.Code, request.Code)
	if err != nil {
		c.Log.Warnf("Can't find discount by code : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		c.Log.Warnf("Discount code has been used : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount.ID = request.ID
	newDiscount.Name = request.Name
	newDiscount.Description = request.Description
	newDiscount.Code = request.Code
	newDiscount.Value = request.Value
	newDiscount.Type = request.Type
	newDiscount.Start, err = time.Parse(layoutTime, request.Start)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	newDiscount.End, err = time.Parse(layoutTime, request.End)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	newDiscount.Status = request.Status
	if err := c.DiscountRepository.Update(tx, newDiscount); err != nil {
		c.Log.Warnf("Can't update discount by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountToResponse(newDiscount), nil
}

func (c *DiscountUseCase) Remove(ctx context.Context, request *model.DeleteDiscountRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newDiscount := new(entity.Discount)
	newDiscount.ID = request.ID
	if err := c.DiscountRepository.Delete(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to delete discount by id: %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}