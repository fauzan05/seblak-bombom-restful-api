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
	UserUseCase 	*UserUseCase
}

func NewAddressUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository, addressRepository *repository.AddressRepository, userUseCase *UserUseCase) *AddressUseCase {
	return &AddressUseCase{
		DB:                db,
		Log:               log,
		Validate:          validate,
		UserRepository:    userRepository,
		AddressRepository: addressRepository,
		UserUseCase: userUseCase,
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
	address.GoogleMapLink = request.GoogleMapLink
	address.IsMain = request.IsMain
	if err := c.AddressRepository.Create(tx, address); err != nil {
		c.Log.Warnf("Failed create new address : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.AddressToResponse(address), nil
}