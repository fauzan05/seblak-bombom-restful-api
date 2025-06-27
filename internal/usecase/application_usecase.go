package usecase

import (
	"context"
	"fmt"
	// "os"
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

type ApplicationUseCase struct {
	DB                    *gorm.DB
	Log                   *logrus.Logger
	Validate              *validator.Validate
	ApplicationRepository *repository.ApplicationRepository
}

func NewApplicationUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, applicationRepository *repository.ApplicationRepository) *ApplicationUseCase {
	return &ApplicationUseCase{
		DB:                    db,
		Log:                   log,
		Validate:              validate,
		ApplicationRepository: applicationRepository,
	}
}

func (c *ApplicationUseCase) Add(ctx *fiber.Ctx, request *model.CreateApplicationRequest) (*model.ApplicationResponse, error) {
	tx := c.DB.WithContext(ctx.Context()).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newApplication := new(entity.Application)
	count, err := c.ApplicationRepository.FindAndCountById(tx, newApplication)
	if err != nil {
		c.Log.Warnf("failed to find application in database : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application in database : %+v", err))
	}

	var hashedFilename string
	if request.Logo != nil {
		err := helper.ValidateFile(1, request.Logo)
		if err != nil {
			c.Log.Warnf(err.Error())
			return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		hashedFilename = helper.HashFileName(request.Logo.Filename)
		err = ctx.SaveFile(request.Logo, fmt.Sprintf("uploads/images/application/%s", hashedFilename))
		if err != nil {
			c.Log.Warnf("failed to save uploaded file: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to save uploaded file: %+v", err))
		}

		// delete data gambar sebelumnya
		// if count > 0 {
		// 	if newApplication.LogoFilename != "" {
		// 		filePath := "uploads/images/application/"
		// 		err = os.Remove(filePath + newApplication.LogoFilename)
		// 		if err != nil {
		// 			c.Log.Warnf("failed to delete image file: %v\n", err)
		// 			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed delete image file: %v\n", err))
		// 		}
		// 	}
		// }
	}

	newApplication.AppName = request.AppName
	if request.Logo != nil {
		newApplication.LogoFilename = hashedFilename
	}
	newApplication.OpeningHours = request.OpeningHours

	newApplication.ClosingHours = request.ClosingHours

	newApplication.Address = request.Address
	newApplication.GoogleMapsLink = request.GoogleMapsLink
	newApplication.Description = request.Description
	newApplication.PhoneNumber = request.PhoneNumber
	newApplication.Email = request.Email
	newApplication.SocialMedia.InstagramName = request.InstagramName
	newApplication.SocialMedia.InstagramLink = request.InstagramLink
	newApplication.SocialMedia.TwitterName = request.TwitterName
	newApplication.SocialMedia.TwitterLink = request.TwitterLink
	newApplication.SocialMedia.FacebookName = request.FacebookName
	newApplication.SocialMedia.FacebookLink = request.FacebookLink
	newApplication.ServiceFee = request.ServiceFee
	// application settings harus berupa 1 baris data saja, tidak boleh lebih dari 2 karena akan membgingunkan nantinya saat pengambilan data mengenai pengaturan aplikasinya
	if count == 0 {
		// boleh dibuat
		if err := c.ApplicationRepository.Create(tx, newApplication); err != nil {
			c.Log.Warnf("failed to create new application request : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create new application request : %+v", err))
		}
	} else if count == 1 {
		// tidak boleh buat lagi, dan mengupdate yang sekarang
		if err := c.ApplicationRepository.Update(tx, newApplication); err != nil {
			c.Log.Warnf("failed to update new application request : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to update new application request : %+v", err))
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.ApplicationToResponse(newApplication), nil
}

func (c *ApplicationUseCase) Get(ctx context.Context) (*model.ApplicationResponse, error) {
	tx := c.DB.WithContext(ctx)

	newApplication := new(entity.Application)
	if err := c.ApplicationRepository.FindFirst(tx, newApplication); err != nil {
		c.Log.Warnf("failed to find application request : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find application request : %+v", err))
	}

	return converter.ApplicationToResponse(newApplication), nil
}
