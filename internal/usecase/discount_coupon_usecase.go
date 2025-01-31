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

type DiscountCouponUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	DiscountCouponRepository *repository.DiscountCouponRepository
}

func NewDiscountCouponUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	DiscountCouponRepository *repository.DiscountCouponRepository) *DiscountCouponUseCase {
	return &DiscountCouponUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		DiscountCouponRepository: DiscountCouponRepository,
	}
}

var layoutTime string = "2006-01-02 15:04:05"

func (c *DiscountCouponUseCase) Add(ctx context.Context, request *model.CreateDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount := new(entity.DiscountCoupon)
	count, err := c.DiscountCouponRepository.CountDiscountByCode(tx, newDiscount, request.Code)
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
	if err := c.DiscountCouponRepository.Create(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to create a new discount : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) GetAll(ctx context.Context) (*[]model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDiscounts := new([]entity.DiscountCoupon)
	if err := c.DiscountCouponRepository.FindAll(tx, newDiscounts); err != nil {
		c.Log.Warnf("Failed to find all discounts : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountCouponsToResponse(newDiscounts), nil
}

func (c *DiscountCouponUseCase) GetById(ctx context.Context, request *model.GetDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDiscount := new(entity.DiscountCoupon)
	newDiscount.ID = request.ID
	if err := c.DiscountCouponRepository.FindById(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to find discount by id: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) Edit(ctx context.Context, request *model.UpdateDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount := new(entity.DiscountCoupon)
	newDiscount.ID = request.ID
	if err := c.DiscountCouponRepository.FindById(tx, newDiscount); err != nil {
		c.Log.Warnf("Can't find discount by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	count, err := c.DiscountCouponRepository.CountDiscountByCodeIsExist(tx, newDiscount, newDiscount.Code, request.Code)
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
	if err := c.DiscountCouponRepository.Update(tx, newDiscount); err != nil {
		c.Log.Warnf("Can't update discount by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) Remove(ctx context.Context, request *model.DeleteDiscountCouponRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newDiscount := new(entity.DiscountCoupon)
	newDiscount.ID = request.ID
	if err := c.DiscountCouponRepository.Delete(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to delete discount by id: %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}