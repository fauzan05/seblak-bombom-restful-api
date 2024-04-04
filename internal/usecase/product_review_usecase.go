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

type ProductReviewUseCase struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	Validate      *validator.Validate
	ProductReview *repository.ProductReviewRepository
}

func NewProductReviewUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	productReview *repository.ProductReviewRepository) *ProductReviewUseCase {
	return &ProductReviewUseCase{
		DB:            db,
		Log:           log,
		Validate:      validate,
		ProductReview: productReview,
	}
}

func (c *ProductReviewUseCase) Add(ctx context.Context, request *model.CreateProductReviewRequest) (*model.ProductReviewResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newProductReview := new(entity.ProductReview)
	newProductReview.ProductId = request.ProductId
	newProductReview.UserId = request.UserId
	newProductReview.Rate = request.Rate
	newProductReview.Comment = request.Comment
	if err := c.ProductReview.Create(tx, newProductReview); err != nil {
		c.Log.Warnf("Failed to create product review from database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ProductReviewToResponse(newProductReview), nil
}
