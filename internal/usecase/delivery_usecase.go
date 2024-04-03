package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeliveryUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	DeliveryRepository *repository.DeliveryRepository
}

func NewDeliveryUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	deliveryRepository *repository.DeliveryRepository) *DeliveryUseCase {
	return &DeliveryUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		DeliveryRepository: deliveryRepository,
	}
}

func (c *DeliveryUseCase) Add(ctx context.Context, request *model.CreateDeliveryRequest) (*model.DeliveryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDelivery := new(entity.Delivery)

	err = c.DeliveryRepository.FindFirst(tx, newDelivery)
	if err == nil {
		c.Log.Warnf("Delivery settings has been exist/created : %+v", err)
		return nil, fiber.ErrConflict
	}

	newDelivery.Cost = request.Cost
	newDelivery.Distance = request.Distance
	if err := c.DeliveryRepository.Create(tx, newDelivery); err != nil {
		c.Log.Warnf("Can't create delivery settings : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DeliveryToResponse(newDelivery), nil
}

func (c *DeliveryUseCase) Get(ctx context.Context) (*model.DeliveryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDelivery := new(entity.Delivery)
	c.DeliveryRepository.FindFirst(tx, newDelivery)

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DeliveryToResponse(newDelivery), nil
}

func (c *DeliveryUseCase) Edit(ctx context.Context, request *model.UpdateDeliveryRequest) (*model.DeliveryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDelivery := new(entity.Delivery)
	newDelivery.ID = request.ID
	count, err := c.DeliveryRepository.FindAndCountById(tx, newDelivery)
	if err != nil {
		c.Log.Warnf("Can't find delivery settings by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count < 1 {
		c.Log.Warnf("Delivery settings by id not found : %+v", err)
		return nil, fiber.ErrNotFound
	}
	newDelivery.Cost = request.Cost
	newDelivery.Distance = request.Distance
	if err := c.DeliveryRepository.Update(tx, newDelivery); err != nil {
		c.Log.Warnf("Can't update delivery settings by : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DeliveryToResponse(newDelivery), nil
}

func (c *DeliveryUseCase) Delete(ctx context.Context, request *model.DeleteDeliveryRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newDelivery := new(entity.Delivery)
	newDelivery.ID = request.ID
	if err := c.DeliveryRepository.Delete(tx, newDelivery); err != nil {
		c.Log.Warnf("Can't delete delivery setting from database : %+v", err)
		return false, fiber.ErrBadRequest
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
