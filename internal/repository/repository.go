package repository

import (
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Pagination struct {
	Page     int
	PageSize int
	Column   string
	SortBy   string
}

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(&entity).Error
}

func (r *Repository[T]) CreateInBatch(db *gorm.DB, entity *[]T) error {
	return db.CreateInBatches(&entity, len(*entity)).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(&entity).Error
}

func (r *Repository[T]) UpdateCustomColumns(db *gorm.DB, entity *T, updateFields map[string]any) error {
	return db.Model(entity).Updates(updateFields).Error
}

func (r *Repository[T]) FindAndCountEntityByUserId(db *gorm.DB, entity *T, userId uint64) (int64, error) {
	var count int64
	err := db.Where("user_id = ?", userId).Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (r *Repository[T]) FindAndCountFirstWalletByUserId(db *gorm.DB, entity *T, userId uint64, status string) (int64, error) {
	var count int64

	// Menggunakan entity langsung tanpa & karena entity sudah merupakan pointer
	err := db.Model(entity).Where("user_id = ? AND status = ?", userId, status).Count(&count).Error
	if err != nil {
		return 0, err
	}

	// Mencari record pertama yang sesuai dengan kriteria
	result := db.Where("user_id = ? AND status = ?", userId, status).First(entity)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Jika data tidak ditemukan, kembalikan count dan nil error
			return count, nil
		}
		// Jika ada error lain, kembalikan error tersebut
		return 0, result.Error
	}

	// Kembalikan count dan nil error jika tidak ada masalah
	return count, nil
}

func (r *Repository[T]) FindFirstWalletByUserId(db *gorm.DB, entity *T, userId uint64, status string) error {
	err := db.Model(entity).Where("user_id = ? AND status = ?", userId, status).First(entity)
	if err != nil {
		return err.Error
	}
	return nil
}

func (r *Repository[T]) FindFirstPayoutByXenditPayoutId(db *gorm.DB, entity *T, xenditPayoutId string) error {
	return db.Where("xendit_payout_id = ?", xenditPayoutId).First(entity).Error
}

func (r *Repository[T]) FindAndUpdateAddressToNonPrimary(db *gorm.DB, entity *T) error {
	var totalAddress int64
	db.Model(entity).Where("is_main = ?", true).Count(&totalAddress)
	if totalAddress > 0 {
		return db.Model(entity).Where("is_main = ?", true).Update("is_main", false).Error
	} else {
		return nil
	}
}

func (r *Repository[T]) UpdateWalletBalance(db *gorm.DB, entity *T, userId uint64, balance float32) error {
	return db.Model(entity).Where("user_id = ?", userId).Update("balance", balance).Error
}

func (r *Repository[T]) FindTokenByUserId(db *gorm.DB, token *T, userId int) error {
	return db.Where("user_id = ?", userId).First(&token).Error
}

func (r *Repository[T]) FindFirst(db *gorm.DB, entity *T) error {
	result := db.First(&entity)

	if result.Error == gorm.ErrRecordNotFound {
		// Jika data tidak ditemukan, kamu bisa mengembalikan nil atau menangani sesuai kebutuhan
		return nil // Tidak ada error jika data tidak ditemukan
	}

	return result.Error // Kembalikan error jika ada kesalahan lain
}

func (r *Repository[T]) FindFirstAndCount(db *gorm.DB, entity *T) (int64, error) {
	var count int64
	err := db.Model(&entity).Count(&count)
	if err.Error != nil {
		return 0, err.Error
	}

	result := db.First(&entity)
	if result.Error == gorm.ErrRecordNotFound {
		// Jika data tidak ditemukan, kamu bisa mengembalikan nil atau menangani sesuai kebutuhan
		return 0, nil // Tidak ada error jika data tidak ditemukan
	}
	return count, result.Error // Kembalikan error jika ada kesalahan lain
}

func (r *Repository[T]) FindXenditTransaction(db *gorm.DB, entity *T, payment_method_id string) (int64, error) {
	var count int64
	err := db.Model(&entity).Where("payment_method_id = ?", payment_method_id).Count(&count)
	if err.Error != nil {
		return 0, err.Error
	}

	result := db.Where("payment_method_id = ?", payment_method_id).First(&entity)
	if result.Error == gorm.ErrRecordNotFound {
		// Jika data tidak ditemukan, kamu bisa mengembalikan nil atau menangani sesuai kebutuhan
		return 0, nil // Tidak ada error jika data tidak ditemukan
	}
	return count, result.Error // Kembalikan error jika ada kesalahan lain
}

func (r *Repository[T]) FindXenditTransactionByPaymentMethodId(db *gorm.DB, entity *T, paymentMethodId string) (int64, error) {
	var count int64
	err := db.Model(&entity).Where("payment_method_id = ?", paymentMethodId).Count(&count)
	if err.Error != nil {
		return 0, err.Error
	}

	result := db.Where("payment_method_id = ?", paymentMethodId).First(&entity)
	if result.Error == gorm.ErrRecordNotFound {
		// Jika data tidak ditemukan, kamu bisa mengembalikan nil atau menangani sesuai kebutuhan
		return 0, nil // Tidak ada error jika data tidak ditemukan
	}
	return count, result.Error // Kembalikan error jika ada kesalahan lain
}

func (r *Repository[T]) FindMidtransCoreAPIOrderByOrderId(db *gorm.DB, entity *T, orderId uint64) error {
	result := db.Where("order_id = ?", orderId).Preload("Actions").First(&entity)

	if result.Error == gorm.ErrRecordNotFound {
		// Jika data tidak ditemukan, kamu bisa mengembalikan nil atau menangani sesuai kebutuhan
		return nil // Tidak ada error jika data tidak ditemukan
	}

	return result.Error // Kembalikan error jika ada kesalahan lain
}

func (r *Repository[T]) FindAllActiveBalance(db *gorm.DB, entity *T) (*float32, error) {
	var totalBalance float32
	result := db.Model(entity).Select("COALESCE(SUM(balance), 0)").Where("status = ?", 1).Scan(&totalBalance)

	return &totalBalance, result.Error
}

func (r *Repository[T]) FindCount(db *gorm.DB, entity *T) (int64, error) {
	var count int64
	err := db.Model(&entity).Count(&count).Error
	return count, err
}

func (r *Repository[T]) FindUserByToken(db *gorm.DB, user *T, token_code string) error {
	token := new(entity.Token)
	// temukan data user_id
	tokenWithUser := db.Where("token = ?", token_code).Joins("User").Find(&token).Error
	if tokenWithUser != nil {
		return tokenWithUser //return errornya
	}
	return db.Where("id = ?", token.UserId).Preload("Token").Preload("Addresses").Preload("Addresses.Delivery").Preload("Wallet").Preload("Cart").Preload("Cart.CartItems").Find(user).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(&entity).Error
}

func (r *Repository[T]) DeleteInBatch(db *gorm.DB, entity *[]T) error {
	return db.Delete(&entity).Error
}

func (r *Repository[T]) DeleteByProductId(db *gorm.DB, entity *T, productId uint64) error {
	return db.Where("product_id = ?", productId).Delete(&entity).Error
}

func (r *Repository[T]) FindById(db *gorm.DB, entity *T) error {
	return db.First(&entity).Error
}

func (r *Repository[T]) FindAddressById(db *gorm.DB, entity *T) error {
	return db.Preload("Delivery").First(&entity).Error
}

func (r *Repository[T]) FindAndCountById(db *gorm.DB, entity *T) (int64, error) {
	var count int64
	err := db.Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (r *Repository[T]) FindAndCountProductById(db *gorm.DB, entity *T) (int64, error) {
	var count int64
	err := db.Preload("Category").Preload("Images").Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (c *Repository[T]) DeleteToken(db *gorm.DB, entity *T, token string) *gorm.DB {
	result := db.Where("token = ?", token).Delete(&entity)
	return result
}

func (r *Repository[T]) FindByEmail(db *gorm.DB, entity *T, email string) error {
	return db.Where("email = ?", email).First(&entity).Error
}

func (r *Repository[T]) FindVerifyToken(db *gorm.DB, entity *T, verifyToken string) error {
	return db.Where("verification_token = ?", verifyToken).First(&entity).Error
}

func (r *Repository[T]) CheckEmailIsExists(db *gorm.DB, currentEmail string, requestEmail string) (int64, error) {
	var total int64
	err := db.Model(&entity.User{}).Where("email = ?", requestEmail).Where("email != ?", currentEmail).Count(&total).Error
	return total, err
}

func (r *Repository[T]) FindUserById(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("id = ?", userId).Preload("Token").Preload("Addresses").Find(&entity).Error
}

func (r *Repository[T]) FindImagesByProductIds(db *gorm.DB, entities *[]T, productIds []uint64) error {
	return db.Where("product_id IN ?", productIds).Find(&entities).Error
}

func (r *Repository[T]) FindUserByIdWithAddress(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("id = ?", userId).Preload("Addresses").Find(&entity).Error
}

func (r *Repository[T]) UserCountByEmail(db *gorm.DB, entity *T, email string) (int64, error) {
	var total int64
	err := db.Model(new(T)).Where("email = ?", email).Count(&total).Error
	return total, err
}

func (r *Repository[T]) DeleteAllByUserId(db *gorm.DB, entity *T, userId uint64) *gorm.DB {
	result := db.Where("user_id = ?", userId).Delete(&entity)
	return result
}

func (r *Repository[T]) FindAll(db *gorm.DB, entities *[]T) error {
	return db.Find(&entities).Error
}

func (r *Repository[T]) FindWithJoins(db *gorm.DB, entity *T, join string) error {
	return db.Joins(join).Find(&entity).Error
}

func (r *Repository[T]) FindWithPreloads(db *gorm.DB, entity *T, preload string) error {
	return db.Preload(preload).Find(&entity).Error
}

func (r *Repository[T]) FindOrderByInvoiceId(db *gorm.DB, entity *T, invoiceId string) error {
	return db.Where("invoice = ?", invoiceId).Preload("OrderProducts").Find(&entity).Error
}

func (r *Repository[T]) FindCurrentUserCartWithPreloads(db *gorm.DB, entity *T, preload string, userId uint64) error {
	return db.Where("user_id = ?", userId).Preload(preload).Find(&entity).Error
}

func (r *Repository[T]) FindWith2Preloads(db *gorm.DB, entity *T, preload1 string, preload2 string) error {
	return db.Preload(preload1).Preload(preload2).Find(&entity).Error
}

func (r *Repository[T]) FindWith3Preloads(db *gorm.DB, entity *T, preload1 string, preload2 string, preload3 string) error {
	return db.Preload(preload1).Preload(preload2).Preload(preload3).Find(&entity).Error
}

func (r *Repository[T]) FindCartItemByUserId(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("user_id = ?", userId).Preload("CartItems").Preload("CartItems.Product").Preload("CartItems.Product.Category").Find(&entity).Error
}

func (r *Repository[T]) FindCartItemByUserIdAndProductId(db *gorm.DB, entity *T, cartId uint64, productId uint64) (int64, error) {
	var count int64
	err := db.Model(&entity).Where("cart_id = ?", cartId).Where("product_id = ?", productId).Count(&count)
	if err.Error != nil {
		return 0, err.Error
	}

	result := db.Where("cart_id = ?", cartId).Where("product_id = ?", productId).First(&entity)
	if result.Error == gorm.ErrRecordNotFound {
		// Jika data tidak ditemukan, kamu bisa mengembalikan nil atau menangani sesuai kebutuhan
		return 0, nil // Tidak ada error jika data tidak ditemukan
	}
	return count, result.Error // Kembalikan error jika ada kesalahan lain
}

func (r *Repository[T]) FirstXenditTransactionByOrderId(db *gorm.DB, entity *T, orderId uint64, preload1 string, preload2 string) error {
	return db.Where("order_id = ?", orderId).Preload(preload1).Preload(preload2).First(&entity).Error
}

func (r *Repository[T]) FindEntityByOrderId(db *gorm.DB, entity *T, orderId uint64) error {
	return db.Where("order_id = ?", orderId).Find(&entity).Error
}

func (r *Repository[T]) FindAllWithJoins(db *gorm.DB, entity *[]T, join string) error {
	return db.Joins(join).Find(&entity).Error
}

func (r *Repository[T]) FindDiscountUsage(db *gorm.DB, entity *T, couponId uint64, userId uint64) error {
	err := db.Where("coupon_id = ?", couponId).Where("user_id = ?", userId).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (r *Repository[T]) CountDiscountByCode(db *gorm.DB, entity *T, discountCode string) (int64, error) {
	var count int64
	err := db.Where("code = ?", discountCode).Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (r *Repository[T]) CountDiscountByCodeIsExist(db *gorm.DB, entity *T, currentCode string, requestCode string) (int64, error) {
	var count int64
	err := db.Where("code = ? AND code != ?", requestCode, currentCode).Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (r *Repository[T]) CountProductItems(db *gorm.DB, entity *T, search string, categoryId uint64) (int64, error) {
	var count int64
	query := db.Where("products.name LIKE ?", "%"+search+"%")
	if categoryId > 0 {
		query.Where("category_id = ?", categoryId)
	}
	err := query.Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (r *Repository[T]) CountCategoryItems(db *gorm.DB, entity *T, search string) (int64, error) {
	var count int64
	err := db.Where("categories.name LIKE ?", "%"+search+"%").Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (r *Repository[T]) FindAllOrdersByUserId(db *gorm.DB, entity *[]T, userId uint64) error {
	return db.Where("user_id = ?", userId).Joins("XenditTransaction").Preload("OrderProducts").Find(&entity).Error
}

func (r *Repository[T]) FindMidtransSnapOrderByOrderId(db *gorm.DB, entity *T, orderId uint64) error {
	return db.Where("order_id = ?", orderId).Find(&entity).Error
}

func (r *Repository[T]) FindCartByUserId(db *gorm.DB, entity *T, userId uint64) error {
	return db.Where("user_id = ?", userId).Find(entity).Error
}

func (r *Repository[T]) FindAllProductByCartItem(db *gorm.DB, entity *[]T, userId uint64) error {
	return db.Where("user_id = ?", userId).Preload("Product").Find(entity).Error
}

func (r *Repository[T]) CountDiscountCouponItems(db *gorm.DB, entity *T, search string) (int64, error) {
	var count int64
	err := db.Where("discount_coupons.name LIKE ?", "%"+search+"%").Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func (r *Repository[T]) CountXenditPayouts(db *gorm.DB, entity *T, search string) (int64, error) {
	var count int64
	err := db.Where("xendit_payouts.amount LIKE ?", "%"+search+"%").Or("xendit_payouts.description LIKE ?", "%"+search+"%").Find(&entity).Count(&count).Error
	if err != nil {
		return int64(0), err
	}
	return count, nil
}

func Paginate[T any](db *gorm.DB, model *T, pagination *Pagination, queryFn func(*gorm.DB) *gorm.DB) ([]T, int64, error) {
	var results []T
	var total int64

	offset := (pagination.Page - 1) * pagination.PageSize

	// Validasi arah sort
	sortBy := strings.ToLower(pagination.SortBy)
	if sortBy != "asc" && sortBy != "desc" {
		sortBy = "asc" // fallback
	}

	// Apply Unscoped before passing to queryFn
	baseQuery := db.Unscoped().Model(model)
	q := queryFn(baseQuery).Order(fmt.Sprintf("%s %s", pagination.Column, pagination.SortBy))

	// Hitung total
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Ambil data
	if err := q.Limit(pagination.PageSize).
		Offset(offset).
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
