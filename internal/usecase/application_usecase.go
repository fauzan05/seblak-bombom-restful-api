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

type ApplicationUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	ApplicationRepository *repository.ApplicationRepository
}

func NewApplicationUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, applicationRepository *repository.ApplicationRepository) *ApplicationUseCase {
	return &ApplicationUseCase{
		DB: db,
		Log: log,
		Validate: validate,
		ApplicationRepository: applicationRepository,
	}
}

func (c *ApplicationUseCase) Add(ctx context.Context, request *model.CreateApplicationRequest) (*model.ApplicationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newApplication := new(entity.Application)
	count, err := c.ApplicationRepository.FindCount(tx, newApplication)
	if err != nil {
		c.Log.Warnf("Failed to find application request : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	if count < 1 {
		// boleh dibuat
		newApplication.AppName = request.AppName
		// validasi format jam
		// parsedTime := time.Date(0, 1, 1, 7, 0, 0, 0, time.UTC)

		newApplication.OpeningHours = request.OpeningHours

		newApplication.ClosingHours = request.ClosingHours
		
		newApplication.Address = request.Address
		newApplication.Longitude = request.Longitude
		newApplication.Latitude = request.Latitude
		newApplication.GoogleMapLink = request.GoogleMapLink
		newApplication.Description = request.Description
		newApplication.PhoneNumber = request.PhoneNumber
		newApplication.Email = request.Email
		newApplication.SocialMedia.InstagramName = request.InstagramName
		newApplication.SocialMedia.InstagramLink = request.InstagramLink
		newApplication.SocialMedia.TwitterName = request.TwitterName
		newApplication.SocialMedia.TwitterLink = request.TwitterLink
		newApplication.SocialMedia.FacebookName = request.FacebookName
		newApplication.SocialMedia.FacebookLink = request.FacebookLink

		if err := c.ApplicationRepository.Create(tx, newApplication); err != nil {
			c.Log.Warnf("Failed to create new application request : %+v", err)
			return nil, fiber.ErrBadRequest
		}

	} else if count > 0{
		// tidak boleh buat lagi
		c.Log.Warnf("Failed to create new application request : %+v", err)
		return nil, fiber.ErrBadRequest
	} 

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ApplicationToResponse(newApplication), nil
}

func (c *ApplicationUseCase) Edit(ctx context.Context, request *model.UpdateApplicationRequest) (*model.ApplicationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newApplication := new(entity.Application)
	count, err := c.ApplicationRepository.FindCount(tx, newApplication)
	if err != nil {
		c.Log.Warnf("Failed to find application request : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	if count == 1 {
		// boleh dibuat
		newApplication.AppName = request.AppName
		// validasi format jam
		// parsedTime := time.Date(0, 1, 1, 7, 0, 0, 0, time.UTC)

		newApplication.OpeningHours = request.OpeningHours

		newApplication.ClosingHours = request.ClosingHours
		
		newApplication.Address = request.Address
		newApplication.Longitude = request.Longitude
		newApplication.Latitude = request.Latitude
		newApplication.GoogleMapLink = request.GoogleMapLink
		newApplication.Description = request.Description
		newApplication.PhoneNumber = request.PhoneNumber
		newApplication.Email = request.Email
		newApplication.SocialMedia.InstagramName = request.InstagramName
		newApplication.SocialMedia.InstagramLink = request.InstagramLink
		newApplication.SocialMedia.TwitterName = request.TwitterName
		newApplication.SocialMedia.TwitterLink = request.TwitterLink
		newApplication.SocialMedia.FacebookName = request.FacebookName
		newApplication.SocialMedia.FacebookLink = request.FacebookLink

		if err := c.ApplicationRepository.Update(tx, newApplication); err != nil {
			c.Log.Warnf("Failed to update new application request : %+v", err)
			return nil, fiber.ErrBadRequest
		}

	} else if count < 1{
		// tidak boleh buat lagi
		c.Log.Warnf("Failed to update new application request : %+v", err)
		return nil, fiber.ErrBadRequest
	} 

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ApplicationToResponse(newApplication), nil
}

func (c *ApplicationUseCase) Get(ctx context.Context) (*model.ApplicationResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newApplication := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApplication); err != nil {
		c.Log.Warnf("Failed to find application request : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ApplicationToResponse(newApplication), nil
}
