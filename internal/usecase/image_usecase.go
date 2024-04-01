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

type ImageUseCase struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	ImageRepository *repository.ImageRepository
}

func NewImageUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	imageRepository *repository.ImageRepository) *ImageUseCase {
	return &ImageUseCase{
		DB:              db,
		Log:             log,
		Validate:        validate,
		ImageRepository: imageRepository,
	}
}


func (c *ImageUseCase) Add(ctx context.Context, request *model.AddImagesRequest) (*[]model.ImageResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newImages := make([]entity.Image, len(request.Images))
	for i, image := range request.Images {
		newImages[i].ProductId = image.ProductId
		newImages[i].FileName = image.FileName
		newImages[i].Type = image.Type
		newImages[i].Position = image.Position
	}

	if err := c.ImageRepository.CreateInBatch(tx, &newImages); err != nil {
		c.Log.Warnf("Failed to add images : %+v", err)
		return nil, fiber.ErrBadRequest
	}
	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}
	return converter.ImagesToResponse(&newImages), nil
}