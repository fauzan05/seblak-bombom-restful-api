package repository

import (
	"fmt"
	"seblak-bombom-restful-api/internal/entity"

	"gorm.io/gorm"
)

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
	return db.Where("id = ?", token.UserId).Preload("Token").Preload("Addresses").Find(user).Error
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

func (r *Repository[T]) FindAndCountById(db *gorm.DB, entity *T) (int64, error) {
	var count int64
	err := db.Find(&entity).Count(&count).Error
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

func (r *Repository[T]) CheckEmailIsExists(db *gorm.DB, currentEmail string, requestEmail string) (int64, error) {
	var total int64
	err := db.Model(&entity.User{}).Where("email = ? AND email != ?", requestEmail, currentEmail).Count(&total).Error
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

func (r *Repository[T]) DeleteAllAddressByUserId(db *gorm.DB, entity *T, userId uint64) *gorm.DB {
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

func (r *Repository[T]) FindCurrentUserCartWithPreloads(db *gorm.DB, entity *T, preload string, userId uint64) error {
	return db.Where("user_id = ?", userId).Preload(preload).Find(&entity).Error
}

func (r *Repository[T]) FindWith2Preloads(db *gorm.DB, entity *T, preload1 string, preload2 string) error {
	return db.Preload(preload1).Preload(preload2).Find(&entity).Error
}

func (r *Repository[T]) FindAllWithJoins(db *gorm.DB, entity *[]T, join string) error {
	return db.Joins(join).Find(&entity).Error
}

func (r *Repository[T]) FindAllWith2Preloads(db *gorm.DB, entity *[]T, preload1 string, preload2 string, page int, pageSize int, search string, specificColumn string, value uint64, columnName string, sortBy string) error {
	offset := (page - 1) * pageSize
	query := db
	// Tambahkan klausa WHERE hanya jika specificColumn tidak kosong dan value tidak nil
	if specificColumn != "" && value != 0 {
		query = query.Where(specificColumn+" = ?", value)
	}

	// Tambahkan klausa WHERE untuk pencarian hanya jika search tidak kosong
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	// Tambahkan preload pertama jika diberikan
	if preload1 != "" {
		query = query.Preload(preload1)
	}

	// tambahkan query sort by
	if columnName != "" && sortBy != "" {
		fmt.Println("KOLOMNYA : ", columnName+" "+sortBy)
		query = query.Order(columnName + " " + sortBy)
	}

	// Tambahkan preload kedua dengan fungsi tambahan jika diberikan
	if preload2 != "" {
		query = query.Preload(preload2, func(db *gorm.DB) *gorm.DB {
			return db.Order("position ASC")
		})
	}

	// Tambahkan sorting, pagination, dan eksekusi query
	return query.Order("id desc").Offset(offset).Limit(pageSize).Find(entity).Error
}

func (r *Repository[T]) FindProductsPagination(db *gorm.DB, entity *[]map[string]interface{}, page int, pageSize int, search string, sortingColumn string, sortBy string, categoryId uint64) error {
	offset := (page - 1) * pageSize
	if sortingColumn == "" {
		sortingColumn = "products.id"
	}

	query := db.Table("products").
		Select(`
        products.id as product_id, 
        products.name as product_name, 
        products.description as product_description, 
        products.price as product_price, 
        products.stock as product_stock, 
        products.created_at as product_created_at, 
        products.updated_at as product_updated_at, 
        categories.id as category_id, 
        categories.name as category_name, 
        categories.description as category_desc, 
        categories.created_at as category_created_at, 
        categories.updated_at as category_updated_at
    `).
		Joins("LEFT JOIN categories ON categories.id = products.category_id").
		Where("products.name LIKE ?", "%"+search+"%").
		Order(fmt.Sprintf("%s %s", sortingColumn, sortBy)).
		Offset(offset).
		Limit(pageSize)

	// Hanya tambahkan klausa WHERE untuk categories.id jika categoryId != 0
	if categoryId > 0 {
		query = query.Where("categories.id = ?", categoryId)
	}

	rows, err := query.Rows()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var (
			productID         string
			productName       string
			productDesc       string
			productPrice      string
			productStock      string
			productCreatedAt  string
			productUpdatedAt  string
			categoryID        string
			categoryName      string
			categoryDesc      string
			categoryCreatedAt string
			categoryUpdatedAt string
		)

		// Scan data
		if err := rows.Scan(
			&productID,
			&productName,
			&productDesc,
			&productPrice,
			&productStock,
			&productCreatedAt,
			&productUpdatedAt,
			&categoryID,
			&categoryName,
			&categoryDesc,
			&categoryCreatedAt,
			&categoryUpdatedAt,
		); err != nil {
			return err
		}

		// Masukkan produk ke hasil
		product := map[string]interface{}{
			"product_id":          productID,
			"product_name":        productName,
			"product_desc":        productDesc,
			"product_price":       productPrice,
			"product_stock":       productStock,
			"product_created_at":  productCreatedAt,
			"product_updated_at":  productUpdatedAt,
			"category_id":         categoryID,
			"category_name":       categoryName,
			"category_desc":       categoryDesc,
			"category_created_at": categoryCreatedAt,
			"category_updated_at": categoryUpdatedAt,
			"images":              []map[string]interface{}{}, // Placeholder untuk images
		}
		results = append(results, product)
	}

	*entity = results
	return nil
}

func (r *Repository[T]) FetchImagesForProducts(db *gorm.DB, products *[]map[string]interface{}) error {
	// Ambil semua product_id
	productIDs := []string{}
	for _, product := range *products {
		productIDs = append(productIDs, product["product_id"].(string))
	}

	// Query semua images untuk product_id tersebut
	rows, err := db.Table("images").
		Select(`
            images.id as image_id,
            images.product_id as image_product_id,
            images.file_name as image_filename,
            images.type as image_type,
            images.position as image_position,
            images.created_at as image_created_at,
            images.updated_at as image_updated_at
        `).
		Where("images.product_id IN (?)", productIDs).
		Order(fmt.Sprintf("%s %s", "image_position", "asc")).
		Rows()

	if err != nil {
		return err
	}
	defer rows.Close()

	// Buat map untuk mengelompokkan images berdasarkan product_id
	imagesMap := make(map[string][]map[string]interface{})
	for rows.Next() {
		var (
			imageID        string
			imageProductId string
			imageFilename  string
			imageType      string
			imagePosition  string
			imageCreatedAt string
			imageUpdatedAt string
		)

		// Scan data
		if err := rows.Scan(
			&imageID,
			&imageProductId,
			&imageFilename,
			&imageType,
			&imagePosition,
			&imageCreatedAt,
			&imageUpdatedAt,
		); err != nil {
			return err
		}

		// Tambahkan ke map
		image := map[string]interface{}{
			"image_id":         imageID,
			"image_filename":   imageFilename,
			"image_position":   imagePosition,
			"image_type":       imageType,
			"image_created_at": imageCreatedAt,
			"image_updated_at": imageUpdatedAt,
		}
		imagesMap[imageProductId] = append(imagesMap[imageProductId], image)
	}

	// Gabungkan images dengan produk masing-masing
	for i, product := range *products {
		productID := product["product_id"].(string)
		if images, exists := imagesMap[productID]; exists {
			(*products)[i]["images"] = images
		}
	}

	return nil
}

func (r *Repository[T]) GetProductsWithPagination(db *gorm.DB, entity *[]map[string]interface{}, page int, pageSize int, search string, sortingColumn string, sortBy string, categoryId uint64) error {
	// Ambil data produk dengan paginasi
	err := r.FindProductsPagination(db, entity, page, pageSize, search, sortingColumn, sortBy, categoryId)
	if err != nil {
		return err
	}

	// Ambil images untuk produk-produk tersebut
	err = r.FetchImagesForProducts(db, entity)
	if err != nil {
		return err
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
	return db.Where("user_id = ?", userId).Joins("MidtransSnapOrder").Preload("OrderProducts").Find(&entity).Error
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

func (r *Repository[T]) GetCategoriesWithPagination(db *gorm.DB, entity *[]map[string]interface{}, page int, pageSize int, search string, sortingColumn string, sortBy string) error {
	// Ambil data produk dengan paginasi
	err := r.FindCategoriesPagination(db, entity, page, pageSize, search, sortingColumn, sortBy)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository[T]) FindCategoriesPagination(db *gorm.DB, entity *[]map[string]interface{}, page int, pageSize int, search string, sortingColumn string, sortBy string) error {
	offset := (page - 1) * pageSize
	if sortingColumn == "" {
		sortingColumn = "categories.id"
	}

	query := db.Table("categories").
		Select(` 
        categories.id as category_id, 
        categories.name as category_name, 
        categories.description as category_desc, 
        categories.created_at as category_created_at, 
        categories.updated_at as category_updated_at
    `).
		Where("categories.name LIKE ?", "%"+search+"%").
		Order(fmt.Sprintf("%s %s", sortingColumn, sortBy)).
		Offset(offset).
		Limit(pageSize)

	rows, err := query.Rows()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var (
			categoryID        string
			categoryName      string
			categoryDesc      string
			categoryCreatedAt string
			categoryUpdatedAt string
		)

		// Scan data
		if err := rows.Scan(
			&categoryID,
			&categoryName,
			&categoryDesc,
			&categoryCreatedAt,
			&categoryUpdatedAt,
		); err != nil {
			return err
		}

		// Masukkan kategori ke hasil
		product := map[string]interface{}{
			"category_id":         categoryID,
			"category_name":       categoryName,
			"category_desc":       categoryDesc,
			"category_created_at": categoryCreatedAt,
			"category_updated_at": categoryUpdatedAt,
			"images":              []map[string]interface{}{}, // Placeholder untuk images
		}
		results = append(results, product)
	}

	*entity = results
	return nil
}