package usecase

import (
	"context"
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
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

type DiscountCouponUseCase struct {
	DB                       *gorm.DB
	Log                      *logrus.Logger
	Validate                 *validator.Validate
	DiscountCouponRepository *repository.DiscountCouponRepository
}

func NewDiscountCouponUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	DiscountCouponRepository *repository.DiscountCouponRepository) *DiscountCouponUseCase {
	return &DiscountCouponUseCase{
		DB:                       db,
		Log:                      log,
		Validate:                 validate,
		DiscountCouponRepository: DiscountCouponRepository,
	}
}

func (c *DiscountCouponUseCase) Add(ctx context.Context, request *model.CreateDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request format. Please check your input!")
	}

	newDiscount := new(entity.DiscountCoupon)
	count, err := c.DiscountCouponRepository.CountDiscountByCode(tx, newDiscount, request.Code)
	if err != nil {
		c.Log.Warnf("Failed to count discount by code : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Oops! Something went wrong. Error Code: 500.")
	}

	if count > 0 {
		c.Log.Warnf("Discount code has been used : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Discount code has already exist, please use another discount code!")
	}

	newDiscount.Name = request.Name
	newDiscount.Description = request.Description
	newDiscount.Code = request.Code
	newDiscount.Value = request.Value
	newDiscount.Type = request.Type
	newDiscount.Start, err = time.Parse(time.RFC3339, request.Start)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request format. Please check your input!")
	}

	newDiscount.End, err = time.Parse(time.RFC3339, request.End)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request format. Please check your input!")
	}

	newDiscount.Status = request.Status
	newDiscount.TotalMaxUsage = request.TotalMaxUsage
	newDiscount.MaxUsagePerUser = request.MaxUsagePerUser
	newDiscount.UsedCount = request.UsedCount
	newDiscount.MinOrderValue = request.MinOrderValue

	if err := c.DiscountCouponRepository.Create(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to create a new discount : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Oops! Something went wrong. Error Code: 500.")
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Oops! Something went wrong. Error Code: 500.")
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) GetAll(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string) (*[]model.DiscountCouponResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if page <= 0 {
		page = 1
	}

	var result []map[string]interface{} // entity kosong yang akan diisi
	if err := c.DiscountCouponRepository.FindDiscountCouponsPagination(tx, &result, page, perPage, search, sortingColumn, sortBy); err != nil {
		c.Log.Warnf("Failed to find all discounts : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	newDiscountCoupons := new([]entity.DiscountCoupon)
	err := MapDiscountCoupon(result, newDiscountCoupons)
	if err != nil {
		c.Log.Warnf("Failed map discount coupons : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	var totalPages int = 0
	newDiscountCoupon := new(entity.DiscountCoupon)
	totalDiscountCoupons, err := c.DiscountCouponRepository.CountDiscountCouponItems(tx, newDiscountCoupon, search)
	if err != nil {
		c.Log.Warnf("Failed to count discount coupon : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	// Hitung total halaman
	totalPages = int(totalDiscountCoupons / int64(perPage))
	if totalDiscountCoupons%int64(perPage) > 0 {
		totalPages++
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, 0, 0, fiber.ErrInternalServerError
	}

	return converter.DiscountCouponsToResponse(newDiscountCoupons), totalDiscountCoupons, totalPages, nil
}

func (c *DiscountCouponUseCase) GetById(ctx context.Context, request *model.GetDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newDiscount := new(entity.DiscountCoupon)
	newDiscount.ID = request.ID
	if err := c.DiscountCouponRepository.FindById(tx, newDiscount); err != nil {
		c.Log.Warnf("Failed to find discount by id: %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) Edit(ctx context.Context, request *model.UpdateDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount := new(entity.DiscountCoupon)
	newDiscount.ID = request.ID
	if err := c.DiscountCouponRepository.FindById(tx, newDiscount); err != nil {
		c.Log.Warnf("Can't find discount by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	count, err := c.DiscountCouponRepository.CountDiscountByCodeIsExist(tx, newDiscount, newDiscount.Code, request.Code)
	if err != nil {
		c.Log.Warnf("Can't find discount by code : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		c.Log.Warnf("Discount code has been used : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount.ID = request.ID
	newDiscount.Name = request.Name
	newDiscount.Description = request.Description
	newDiscount.Code = request.Code
	newDiscount.Value = request.Value
	newDiscount.Type = request.Type
	newDiscount.Start, err = time.Parse(time.RFC3339, request.Start)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount.End, err = time.Parse(time.RFC3339, request.End)
	if err != nil {
		c.Log.Warnf("Can't parse to time : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newDiscount.Status = request.Status
	newDiscount.TotalMaxUsage = request.TotalMaxUsage
	newDiscount.MaxUsagePerUser = request.MaxUsagePerUser
	newDiscount.UsedCount = request.UsedCount
	newDiscount.MinOrderValue = request.MinOrderValue

	if err := c.DiscountCouponRepository.Update(tx, newDiscount); err != nil {
		c.Log.Warnf("Can't update discount by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) Remove(ctx context.Context, request *model.DeleteDiscountCouponRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return false, fiber.ErrBadRequest
	}

	newDiscountCoupons := []entity.DiscountCoupon{}

	for _, idCoupon := range request.IDs {
		newDiscountCoupon := entity.DiscountCoupon{
			ID: idCoupon,
		}

		newDiscountCoupons = append(newDiscountCoupons, newDiscountCoupon)
	}

	if err := c.DiscountCouponRepository.DeleteInBatch(tx, &newDiscountCoupons); err != nil {
		c.Log.Warnf("Failed to delete discount by id: %+v", err)
		return false, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return false, fiber.ErrInternalServerError
	}

	return true, nil
}

func MapDiscountCoupon(rows []map[string]interface{}, results *[]entity.DiscountCoupon) error {

	for _, row := range rows {
		discountCouponIdStr, ok := row["discount_coupon_id"].(string)
		if !ok || discountCouponIdStr == "" {
			return fmt.Errorf("missing or invalid discount_coupon_id")
		}

		discountCouponId, err := strconv.ParseUint(discountCouponIdStr, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse category_id: %v", err)
		}

		discountCouponName, _ := row["discount_coupon_name"].(string)
		discountCouponDesc, _ := row["discount_coupon_desc"].(string)
		discountCouponCode, _ := row["discount_coupon_code"].(string)
		discountCouponValue, _ := strconv.ParseFloat(row["discount_coupon_value"].(string), 32)

		discountCouponType, _ := strconv.Atoi(row["discount_coupon_type"].(string))
		discountCouponStartStr, _ := row["discount_coupon_start"].(string)
		discountCouponStart, err := time.Parse(time.RFC3339, discountCouponStartStr)
		if err != nil {
			return fmt.Errorf("failed to parse discount_coupon_start: %v", err)
		}
		discountCouponEndStr, _ := row["discount_coupon_end"].(string)
		discountCouponEnd, err := time.Parse(time.RFC3339, discountCouponEndStr)
		if err != nil {
			return fmt.Errorf("failed to parse discount_coupon_end: %v", err)
		}
		discountCouponTotalMaxUsage, _ := strconv.Atoi(row["discount_coupon_total_max_usage"].(string))
		discountCouponMaxUsagePerUser, _ := strconv.Atoi(row["discount_coupon_max_usage_per_user"].(string))
		discountCouponUsedCount, _ := strconv.Atoi(row["discount_coupon_used_count"].(string))
		discountCouponMinOrderValue, _ := strconv.Atoi(row["discount_coupon_min_order_value"].(string))
		discountCouponStatus, _ := strconv.ParseBool(row["discount_coupon_status"].(string))

		discountCouponCreatedAtStr, _ := row["discount_coupon_created_at"].(string)
		discountCouponCreatedAt, err := time.Parse(time.RFC3339, discountCouponCreatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse discount_coupon_created_at: %v", err)
		}

		discountCouponUpdatedAtStr, _ := row["discount_coupon_updated_at"].(string)
		discountCouponUpdatedAt, err := time.Parse(time.RFC3339, discountCouponUpdatedAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse discount_coupon_updated_at: %v", err)
		}

		newDiscountCoupon := entity.DiscountCoupon{
			ID:              discountCouponId,
			Name:            discountCouponName,
			Description:     discountCouponDesc,
			Code:            discountCouponCode,
			Value:           float32(discountCouponValue),
			Type:            helper.DiscountType(discountCouponType),
			Start:           discountCouponStart,
			End:             discountCouponEnd,
			TotalMaxUsage:   discountCouponTotalMaxUsage,
			MaxUsagePerUser: discountCouponMaxUsagePerUser,
			UsedCount:       discountCouponUsedCount,
			MinOrderValue:   discountCouponMinOrderValue,
			Status:          discountCouponStatus,
			Created_At:      discountCouponCreatedAt,
			Updated_At:      discountCouponUpdatedAt,
		}

		// Tambahkan ke hasil
		*results = append(*results, newDiscountCoupon)
	}

	return nil
}
