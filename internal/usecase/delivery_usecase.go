package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"strconv"
	"time"

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
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newDelivery := new(entity.Delivery)
	newDelivery.District = request.District
	newDelivery.City = request.City
	newDelivery.Village = request.Village
	newDelivery.Hamlet = request.Hamlet
	newDelivery.Cost = request.Cost
	if err := c.DeliveryRepository.Create(tx, newDelivery); err != nil {
		c.Log.Warnf("Can't create delivery settings : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Can't create delivery settings : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.DeliveryToResponse(newDelivery), nil
}

func (c *DeliveryUseCase) GetAll(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string) (*[]model.DeliveryResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx)

	if page <= 0 {
		page = 1
	}

	var result []map[string]any // entity kosong yang akan diisi
	if err := c.DeliveryRepository.FindDeliveriesPagination(tx, &result, page, perPage, search, sortingColumn, sortBy); err != nil {
		c.Log.Warnf("Failed to find all deliveries : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to find all deliveries : %+v", err))
	}

	newDelivery := new([]entity.Delivery)
	err := MapDeliveries(result, newDelivery)
	if err != nil {
		c.Log.Warnf("Failed to map delivery : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed map delivery : %+v", err))
	}

	var totalPages int = 0
	getAllDelivery := new(entity.Delivery)
	totalDeliveries, err := c.DeliveryRepository.CountDeliveryItems(tx, getAllDelivery, search)
	if err != nil {
		c.Log.Warnf("Failed to count products: %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to count products: %+v", err))
	}

	// Hitung total halaman
	totalPages = int(totalDeliveries / int64(perPage))
	if totalDeliveries%int64(perPage) > 0 {
		totalPages++
	}

	return converter.DeliveriesToResponse(newDelivery), totalDeliveries, totalPages, nil
}

func (c *DeliveryUseCase) Edit(ctx context.Context, request *model.UpdateDeliveryRequest) (*model.DeliveryResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newDelivery := new(entity.Delivery)
	newDelivery.ID = request.ID
	count, err := c.DeliveryRepository.FindAndCountById(tx, newDelivery)
	if err != nil {
		c.Log.Warnf("Can't find delivery settings by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Can't find delivery settings by id : %+v", err))
	}

	if count < 1 {
		c.Log.Warnf("Delivery settings not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "Delivery settings not found!")
	}

	newDelivery.District = request.District
	newDelivery.City = request.City
	newDelivery.Village = request.Village
	newDelivery.Hamlet = request.Hamlet
	newDelivery.Cost = request.Cost

	fmt.Println(newDelivery);
	if err := c.DeliveryRepository.Update(tx, newDelivery); err != nil {
		c.Log.Warnf("Can't update delivery settings by : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Can't update delivery settings by : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return converter.DeliveryToResponse(newDelivery), nil
}

func (c *DeliveryUseCase) Delete(ctx context.Context, request *model.DeleteDeliveryRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Invalid request body : %+v", err))
	}

	newDeliveries := []entity.Delivery{}
	for _, idDelivery := range request.IDs {
		newDelivery := entity.Delivery{
			ID: idDelivery,
		}
		newDeliveries = append(newDeliveries, newDelivery)
	}

	if err := c.DeliveryRepository.DeleteInBatch(tx, &newDeliveries); err != nil {
		c.Log.Warnf("Can't delete delivery from database : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Can't delete delivery from database : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to commit transaction : %+v", err))
	}

	return true, nil
}

func MapDeliveries(rows []map[string]interface{}, results *[]entity.Delivery) error {

	for _, row := range rows {
		// Ambil dan validasi delivery_id
		deliveryIdStr, ok := row["delivery_id"].(string)
		if !ok || deliveryIdStr == "" {
			return fmt.Errorf("missing or invalid delivery_id")
		}

		deliveryId, err := strconv.ParseUint(deliveryIdStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse delivery_id: %v", err)
		}

		// Ambil field kategori
		deliveryCity, _ := row["delivery_city"].(string)
		deliveryDistrict, _ := row["delivery_district"].(string)
		deliveryVillage, _ := row["delivery_village"].(string)
		deliveryHamlet, _ := row["delivery_hamlet"].(string)
		deliveryCostStr, _ := row["delivery_cost"].(string)
		deliveryCost, err := strconv.ParseFloat(deliveryCostStr, 32)
		if err != nil {
			return fmt.Errorf("failed to parse delivery_cost : %v", err)
		}

		// Parse created_at dan updated_at kategori
		deliveryCreatedAtStr, _ := row["delivery_created_at"].(string)
		categoryCreatedAt, err := time.Parse(time.RFC3339, deliveryCreatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse delivery_created_at: %v", err)
		}

		deliveryUpdatedAtStr, _ := row["delivery_updated_at"].(string)
		categoryUpdatedAt, err := time.Parse(time.RFC3339, deliveryUpdatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse delivery_updated_at: %v", err)
		}

		// Buat objek kategori
		newDelivery := entity.Delivery{
			ID:         deliveryId,
			City:       deliveryCity,
			District:   deliveryDistrict,
			Village:    deliveryVillage,
			Hamlet:     deliveryHamlet,
			Cost:       float32(deliveryCost),
			CreatedAt: categoryCreatedAt,
			UpdatedAt: categoryUpdatedAt,
		}

		// Tambahkan ke hasil
		*results = append(*results, newDelivery)
	}

	return nil
}
