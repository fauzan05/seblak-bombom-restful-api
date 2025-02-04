package usecase

import (
	"context"
	"fmt"
	"os"
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
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newApplication := new(entity.Application)
	if request.ID > 0 {
		newApplication.ID = request.ID
		if err := c.ApplicationRepository.FindById(tx, newApplication); err != nil {
			c.Log.Warnf("Failed to find current application data in database : %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to find current application data in database")
		}
	}

	var hashedFilename string
	if request.Logo != nil {
		hashedFilename = hashFileName(request.Logo.Filename)
		err = ctx.SaveFile(request.Logo, fmt.Sprintf("../uploads/images/application/%s", hashedFilename))
		if err != nil {
			c.Log.Warnf("Failed to save file: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to save uploaded file!")
		}

		// delete data gambar sebelumnya
		if newApplication.ID > 0 {
			if newApplication.LogoFilename != "" {
				filePath := "../uploads/images/application/"
				err = os.Remove(filePath + newApplication.LogoFilename)
				if err != nil {
					fmt.Printf("Error deleting file: %v\n", err)
				return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to delete logo images in directory!")
				}
			}
		}
	}

	count, err := c.ApplicationRepository.FindCount(tx, newApplication)
	if err != nil {
		c.Log.Warnf("Failed to find application request : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to find application request!")
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

	// application settings harus berupa 1 baris data saja, tidak boleh lebih dari 2 karena akan membgingunkan nantinya saat pengambilan data mengenai pengaturan aplikasinya
	if count < 1 {
		// boleh dibuat
		if err := c.ApplicationRepository.Create(tx, newApplication); err != nil {
			c.Log.Warnf("Failed to create new application request : %+v", err)
			return nil, fiber.ErrBadRequest
		}
	} else if count > 0 {
		// tidak boleh buat lagi, dan mengupdate yang sekarang
		if err := c.ApplicationRepository.Update(tx, newApplication); err != nil {
			c.Log.Warnf("Failed to create new application request : %+v", err)
			return nil, fiber.ErrBadRequest
		}
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
