package usecase

import (
	"context"
	"fmt"
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
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newDelivery := new(entity.Delivery)
	newDelivery.District = request.District
	newDelivery.City = request.City
	newDelivery.Village = request.Village
	newDelivery.Hamlet = request.Hamlet
	newDelivery.Cost = request.Cost
	if err := c.DeliveryRepository.Create(tx, newDelivery); err != nil {
		c.Log.Warnf("can't create delivery settings : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't create delivery settings : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.DeliveryToResponse(newDelivery), nil
}

func (c *DeliveryUseCase) GetAll(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string) (*[]model.DeliveryResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx)

	if page <= 0 {
		page = 1
	}

	if sortingColumn == "" {
		sortingColumn = "deliveries.id"
	}

	newPagination := new(repository.Pagination)
	newPagination.Page = page
	newPagination.PageSize = perPage
	newPagination.Column = sortingColumn
	newPagination.SortBy = sortBy
	allowedColumns := map[string]bool{
		"deliveries.id":         true,
		"deliveries.cost":       true,
		"deliveries.city":       true,
		"deliveries.district":   true,
		"deliveries.village":    true,
		"deliveries.hamlet":     true,
		"deliveries.created_at": true,
		"deliveries.updated_at": true,
	}

	if !allowedColumns[newPagination.Column] {
		c.Log.Warnf("invalid sort column : %s", newPagination.Column)
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid sort column : %s", newPagination.Column))
	}
	
	deliveries, totalDelivery, err := repository.Paginate(tx, &entity.Delivery{}, newPagination, func(d *gorm.DB) *gorm.DB {
		return d.Where("deliveries.cost LIKE ?", "%"+search+"%").
			Or("deliveries.city LIKE ?", "%"+search+"%").
			Or("deliveries.district LIKE ?", "%"+search+"%").
			Or("deliveries.village LIKE ?", "%"+search+"%").
			Or("deliveries.hamlet LIKE ?", "%"+search+"%")
	})

	if err != nil {
		c.Log.Warnf("failed to paginate deliveries : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to paginate deliveries : %+v", err))
	}

	// Hitung total halaman
	var totalPages int = 0
	totalPages = int(totalDelivery / int64(perPage))
	if totalDelivery%int64(perPage) > 0 {
		totalPages++
	}

	return converter.DeliveriesToResponse(&deliveries), totalDelivery, totalPages, nil
}

func (c *DeliveryUseCase) Edit(ctx context.Context, request *model.UpdateDeliveryRequest) (*model.DeliveryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newDelivery := new(entity.Delivery)
	newDelivery.ID = request.ID
	count, err := c.DeliveryRepository.FindAndCountById(tx, newDelivery)
	if err != nil {
		c.Log.Warnf("can't find delivery settings by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't find delivery settings by id : %+v", err))
	}

	if count < 1 {
		c.Log.Warnf("delivery settings not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "delivery settings not found!")
	}

	newDelivery.District = request.District
	newDelivery.City = request.City
	newDelivery.Village = request.Village
	newDelivery.Hamlet = request.Hamlet
	newDelivery.Cost = request.Cost

	fmt.Println(newDelivery)
	if err := c.DeliveryRepository.Update(tx, newDelivery); err != nil {
		c.Log.Warnf("can't update delivery settings by : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't update delivery settings by : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.DeliveryToResponse(newDelivery), nil
}

func (c *DeliveryUseCase) Delete(ctx context.Context, request *model.DeleteDeliveryRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newDeliveries := []entity.Delivery{}
	for _, idDelivery := range request.IDs {
		newDelivery := entity.Delivery{
			ID: idDelivery,
		}
		newDeliveries = append(newDeliveries, newDelivery)
	}

	if err := c.DeliveryRepository.DeleteInBatch(tx, &newDeliveries); err != nil {
		c.Log.Warnf("can't delete delivery from database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't delete delivery from database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}
