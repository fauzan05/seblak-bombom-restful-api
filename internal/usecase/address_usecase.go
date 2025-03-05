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

type AddressUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	UserRepository     *repository.UserRepository
	AddressRepository  *repository.AddressRepository
	DeliveryRepository *repository.DeliveryRepository
	UserUseCase        *UserUseCase
}

func NewAddressUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, addressRepository *repository.AddressRepository,
	deliveryRepository *repository.DeliveryRepository, userUseCase *UserUseCase) *AddressUseCase {
	return &AddressUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		UserRepository:     userRepository,
		AddressRepository:  addressRepository,
		DeliveryRepository: deliveryRepository,
		UserUseCase:        userUseCase,
	}
}

func (c *AddressUseCase) Create(ctx context.Context, request *model.AddressCreateRequest, token *model.GetUserByTokenRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	currentUser, err := c.UserUseCase.GetUserByToken(ctx, token)
	if err != nil {
		c.Log.Warnf("Token isn't valid : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDelivery := new(entity.Delivery)
	newDelivery.ID = request.DeliveryId

	// cek apakah delivery id ada
	if err := c.DeliveryRepository.FindFirst(tx, newDelivery); err != nil {
		c.Log.Warnf("Failed create new address : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if newDelivery.City == "" {
		c.Log.Warnf("Failed to find delivery data : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Delivery data is not found!")
	}

	address := new(entity.Address)
	// update yang tadinya is_main = 1 menjadi 0
	if request.IsMain {
		if err := c.AddressRepository.FindAndUpdateAddressToNonPrimary(tx, address); err != nil {
			c.Log.Warnf("Failed to update address is main to non-primary : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	address.UserId = currentUser.ID
	address.DeliveryId = request.DeliveryId
	address.CompleteAddress = request.CompleteAddress
	address.GoogleMapsLink = request.GoogleMapsLink
	address.IsMain = request.IsMain
	if err := c.AddressRepository.Create(tx, address); err != nil {
		c.Log.Warnf("Failed create new address : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := c.AddressRepository.FindAddressById(tx, address); err != nil {
		c.Log.Warnf("Failed find updated address by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AddressToResponse(address), nil
}

func (c *AddressUseCase) GetAll(user *model.UserResponse) (*[]model.AddressResponse, error) {
	return converter.AddressesToResponse(&user.Addresses), nil
}

func (c *AddressUseCase) GetById(ctx context.Context, request *model.GetAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request query params : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newAddress := new(entity.Address)
	newAddress.ID = request.ID
	if err := c.AddressRepository.FindById(tx, newAddress); err != nil {
		c.Log.Warnf("Failed find address by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AddressToResponse(newAddress), nil
}

func (c *AddressUseCase) Edit(ctx context.Context, request *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newAddress := new(entity.Address)
	newAddress.ID = request.ID
	newAddress.UserId = request.UserId
	newAddress.CompleteAddress = request.CompleteAddress
	newAddress.IsMain = request.IsMain

	if err := c.AddressRepository.Update(tx, newAddress); err != nil {
		c.Log.Warnf("Failed edit address by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := c.AddressRepository.FindById(tx, newAddress); err != nil {
		c.Log.Warnf("Failed find updated address by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AddressToResponse(newAddress), nil
}

func (c *AddressUseCase) Delete(ctx context.Context, request *model.DeleteAddressRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := c.Validate.Struct(request); err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newAddresses := []entity.Address{}
	for _, idAddress := range request.IDs {
		newAddress := entity.Address{
			ID: idAddress,
		}
		newAddresses = append(newAddresses, newAddress)
	}

	if err := c.AddressRepository.DeleteInBatch(tx, &newAddresses); err != nil {
		c.Log.Warnf("Can't delete address by id : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
