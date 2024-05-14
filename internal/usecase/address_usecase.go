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
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	UserRepository    *repository.UserRepository
	AddressRepository *repository.AddressRepository
	UserUseCase       *UserUseCase
}

func NewAddressUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, addressRepository *repository.AddressRepository, userUseCase *UserUseCase) *AddressUseCase {
	return &AddressUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		UserRepository:    userRepository,
		AddressRepository: addressRepository,
		UserUseCase:       userUseCase,
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

	address := new(entity.Address)
	address.UserId = currentUser.ID
	address.Regency = request.Regency
	address.SubDistrict = request.Subdistrict
	address.CompleteAddress = request.CompleteAddress
	address.Longitude = request.Longitude
	address.Latitude = request.Latitude
	address.IsMain = request.IsMain

	if err := c.AddressRepository.Create(tx, address); err != nil {
		c.Log.Warnf("Failed create new address : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := c.AddressRepository.FindById(tx, address); err != nil {
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
	newAddress.Regency = request.Regency
	newAddress.SubDistrict = request.Subdistrict
	newAddress.CompleteAddress = request.CompleteAddress
	newAddress.Longitude = request.Longitude
	newAddress.Latitude = request.Latitude
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

	newAddress := new(entity.Address)
	newAddress.ID = request.ID
	if err := c.AddressRepository.Delete(tx, newAddress); err != nil {
		c.Log.Warnf("Can't delete address by id : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}
