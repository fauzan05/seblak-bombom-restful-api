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
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newDiscount := new(entity.DiscountCoupon)
	count, err := c.DiscountCouponRepository.CountDiscountByCode(tx, newDiscount, request.Code)
	if err != nil {
		c.Log.Warnf("failed to count discount by code : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to count discount by code : %+v", err))
	}

	if count > 0 {
		c.Log.Warnf("discount code has already exist, please use another discount code!")
		return nil, fiber.NewError(fiber.StatusConflict, "discount code has already exist, please use another discount code!")
	}

	newDiscount.Name = request.Name
	newDiscount.Description = request.Description
	newDiscount.Code = request.Code
	newDiscount.Value = request.Value
	newDiscount.Type = request.Type
	newDiscount.Start = request.Start.ToTime()
	newDiscount.End = request.End.ToTime()
	newDiscount.Status = request.Status
	newDiscount.TotalMaxUsage = request.TotalMaxUsage
	newDiscount.MaxUsagePerUser = request.MaxUsagePerUser
	newDiscount.UsedCount = request.UsedCount
	newDiscount.MinOrderValue = request.MinOrderValue

	if err := c.DiscountCouponRepository.Create(tx, newDiscount); err != nil {
		c.Log.Warnf("failed to create a new discount : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to create a new discount : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) GetAll(ctx context.Context, page int, perPage int, search string, sortingColumn string, sortBy string, status bool) (*[]model.DiscountCouponResponse, int64, int, error) {
	tx := c.DB.WithContext(ctx)

	if page <= 0 {
		page = 1
	}

	if sortingColumn == "" {
		sortingColumn = "discount_coupons.id"
	}

	newPagination := new(repository.Pagination)
	newPagination.Page = page
	newPagination.PageSize = perPage
	newPagination.Column = sortingColumn
	newPagination.SortBy = sortBy
	allowedColumns := map[string]bool{
		"discount_coupons.id":          true,
		"discount_coupons.name":        true,
		"discount_coupons.description": true,
		"discount_coupons.created_at":  true,
		"discount_coupons.updated_at":  true,
	}

	if !allowedColumns[newPagination.Column] {
		c.Log.Warnf("invalid sort column : %s", newPagination.Column)
		return nil, 0, 0, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid sort column : %s", newPagination.Column))
	}

	discountCoupons, totalDiscountCoupon, err := repository.Paginate(tx, &entity.DiscountCoupon{}, newPagination, func(d *gorm.DB) *gorm.DB {
		return d.Where(
			d.Where("discount_coupons.name LIKE ?", "%"+search+"%").
				Or("discount_coupons.code LIKE ?", "%"+search+"%").
				Or("discount_coupons.value LIKE ?", "%"+search+"%"),
		).Where("discount_coupons.status = ?", status)
	})

	if err != nil {
		c.Log.Warnf("failed to paginate discount coupons : %+v", err)
		return nil, 0, 0, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to paginate discount coupons : %+v", err))
	}

	// Hitung total halaman
	var totalPages int = 0
	totalPages = int(totalDiscountCoupon / int64(perPage))
	if totalDiscountCoupon%int64(perPage) > 0 {
		totalPages++
	}

	return converter.DiscountCouponsToResponse(&discountCoupons), totalDiscountCoupon, totalPages, nil
}

func (c *DiscountCouponUseCase) GetById(ctx context.Context, request *model.GetDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx)

	newDiscount := new(entity.DiscountCoupon)
	newDiscount.ID = request.ID
	if err := c.DiscountCouponRepository.FindById(tx, newDiscount); err != nil {
		c.Log.Warnf("failed to find discount by id: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to find discount by id: %+v", err))
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) Edit(ctx context.Context, request *model.UpdateDiscountCouponRequest) (*model.DiscountCouponResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newDiscount := new(entity.DiscountCoupon)
	newDiscount.ID = request.ID
	count, err := c.DiscountCouponRepository.FindAndCountById(tx, newDiscount)
	if err != nil {
		c.Log.Warnf("can't find discount by code : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't find discount by code : %+v", err))
	}

	if count == 0 {
		c.Log.Warnf("discount coupon not found!")
		return nil, fiber.NewError(fiber.StatusNotFound, "discount coupon not found!")
	}

	count, err = c.DiscountCouponRepository.CountDiscountByCodeIsExist(tx, newDiscount, newDiscount.Code, request.Code)
	if err != nil {
		c.Log.Warnf("can't find discount by code : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't find discount by code : %+v", err))
	}

	if count > 0 {
		c.Log.Warnf("discount code has been used : %+v", err)
		return nil, fiber.NewError(fiber.StatusConflict, fmt.Sprintf("discount code has been used : %+v", err))
	}

	newDiscount.ID = request.ID
	newDiscount.Name = request.Name
	newDiscount.Description = request.Description
	newDiscount.Code = request.Code
	newDiscount.Value = request.Value
	newDiscount.Type = request.Type
	newDiscount.Start = request.Start.ToTime()
	newDiscount.End = request.End.ToTime()
	newDiscount.Status = request.Status
	newDiscount.TotalMaxUsage = request.TotalMaxUsage
	newDiscount.MaxUsagePerUser = request.MaxUsagePerUser
	newDiscount.UsedCount = request.UsedCount
	newDiscount.MinOrderValue = request.MinOrderValue

	if err := c.DiscountCouponRepository.Update(tx, newDiscount); err != nil {
		c.Log.Warnf("can't update discount by id : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("can't update discount by id : %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return converter.DiscountCouponToResponse(newDiscount), nil
}

func (c *DiscountCouponUseCase) Remove(ctx context.Context, request *model.DeleteDiscountCouponRequest) (bool, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("invalid request body : %+v", err)
		return false, fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("invalid request body : %+v", err))
	}

	newDiscountCoupons := []entity.DiscountCoupon{}

	for _, idCoupon := range request.IDs {
		newDiscountCoupon := entity.DiscountCoupon{
			ID: idCoupon,
		}

		newDiscountCoupons = append(newDiscountCoupons, newDiscountCoupon)
	}

	if err := c.DiscountCouponRepository.DeleteInBatch(tx, &newDiscountCoupons); err != nil {
		c.Log.Warnf("failed to delete discount by id: %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to delete discount by id: %+v", err))
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("failed to commit transaction : %+v", err)
		return false, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("failed to commit transaction : %+v", err))
	}

	return true, nil
}
