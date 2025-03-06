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

type CartUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CartRepository     *repository.CartRepository
	ProductRepository  *repository.ProductRepository
	CartItemRepository *repository.CartItemRepository
}

func NewCartUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	cartRepository *repository.CartRepository, productRepository *repository.ProductRepository,
	cartItemRepository *repository.CartItemRepository) *CartUseCase {
	return &CartUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		CartRepository:     cartRepository,
		ProductRepository:  productRepository,
		CartItemRepository: cartItemRepository,
	}
}

func (c *CartUseCase) Add(ctx context.Context, request *model.CreateCartRequest) (*model.CartResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// dicek apakah quantity yang dikirim itu negatif/nol
	if request.Quantity < 1 {
		c.Log.Warnf("Quantity must be more than 0 if first time adding product to cart!")
		return nil, fiber.NewError(fiber.StatusBadRequest, "Quantity must be more than 0 if first time adding product to cart!")
	}

	// dicek apakah produknya ada atau tidak
	newProduct := new(entity.Product)
	newProduct.ID = request.ProductID
	if err := c.ProductRepository.FindById(tx, newProduct); err != nil {
		c.Log.Warnf("Failed to find product by id into product table : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// cek apakah produk tersedia atau tidak
	if newProduct.Stock < 1 {
		c.Log.Warnf("Product was out of stock : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// cek apakah permintaan melebihi stok yang tersedia
	newProduct.Stock -= request.Quantity
	if newProduct.Stock < 0 {
		// jika jumlah kuantitasnya melebihi stok yang tersedia
		c.Log.Warnf("Quantity request out of stock from product : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// dicek terlebih dahulu apakah ada cart dengan user yang sama dan produk yang sama.
	newCart := new(entity.Cart)
	if err := c.CartRepository.FindCartByUserId(tx, newCart, request.UserID); err != nil {
		c.Log.Warnf("Failed to find cart by user id from cart table : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newCartItem := new(entity.CartItem)
	if newCart.ID == 0 {
		// jika tidak ada maka buat cart baru
		newCart.UserID = request.UserID
		if err := c.CartRepository.Create(tx, newCart); err != nil {
			c.Log.Warnf("Failed to create cart by user id into cart table : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		// setelah itu insert datanya ke tabel cart_items
		newCartItem.CartId = newCart.ID
		newCartItem.ProductID = request.ProductID
		newCartItem.Quantity = request.Quantity
		if err := c.CartItemRepository.Create(tx, newCartItem); err != nil {
			c.Log.Warnf("Failed to create cart item by user id into cart table : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		// kurangi stok produknya
		if err := c.ProductRepository.Update(tx, newProduct); err != nil {
			c.Log.Warnf("Failed to update product quantity into product table : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

	} else if newCart.ID > 0 {
		// jika tidak ada maka gunakan cart yang ada
		if err := c.CartRepository.FindWithPreloads(tx, newCart, "CartItems"); err != nil {
			c.Log.Warnf("Failed to find cart item by user id from cart from database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		// temukan apakah ada user id dan produk id yang sama di cart items
		totalCartItem, err := c.CartItemRepository.FindCartItemByUserIdAndProductId(tx, newCartItem, newCart.ID, request.ProductID)
		if err != nil {
			c.Log.Warnf("Failed to find cart item by user id and product id from cart from database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		if totalCartItem == 0 {
			// tambahkan baru
			newCartItem := new(entity.CartItem)
			newCartItem.CartId = newCart.ID
			newCartItem.ProductID = request.ProductID
			newCartItem.Quantity = request.Quantity

			if err := c.CartItemRepository.Create(tx, newCartItem); err != nil {
				c.Log.Warnf("Failed to create cart item by user id into database : %+v", err)
				return nil, fiber.ErrInternalServerError
			}
		} else {
			// update quantitynya saja
			newCartItem.Quantity = newCartItem.Quantity + request.Quantity
			updateQuantity := map[string]any{
				"quantity": newCartItem.Quantity,
			}

			if err := c.CartItemRepository.UpdateCustomColumns(tx, newCartItem, updateQuantity); err != nil {
				c.Log.Warnf("Failed to update quantity cart item into database : %+v", err)
				return nil, fiber.ErrInternalServerError
			}
		}

		// lalu kurangi stoknya
		newProductUpdate := new(entity.Product)
		newProductUpdate.ID = newProduct.ID
		updateProductQuantity := map[string]any{
			"stock": newProduct.Stock,
		}

		if err := c.ProductRepository.UpdateCustomColumns(tx, newProductUpdate, updateProductQuantity); err != nil {
			c.Log.Warnf("Failed to update stock product into database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}
	
	newCart.CartItems = nil

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CartToResponse(newCart), nil
}

func (c *CartUseCase) GetAllByCurrentUser(ctx context.Context, request *model.GetAllCartByCurrentUserRequest) (*model.CartResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newCart := new(entity.Cart)
	if err := c.CartRepository.FindCartItemByUserId(tx, newCart, request.UserID); err != nil {
		c.Log.Warnf("Failed to find all cart item by current user in the database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CartToResponse(newCart), nil
}

func (c *CartUseCase) UpdateQuantity(ctx context.Context, request *model.UpdateCartRequest) (*model.CartItemResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newCartItem := new(entity.CartItem)
	newCartItem.ID = request.CartItemID
	count, err := c.CartItemRepository.FindAndCountById(tx, newCartItem)
	if err != nil {
		c.Log.Warnf("Failed to find cart item by cart item id in the database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count == 0 {
		c.Log.Warnf("Cart item by id not found!")
		return nil, fiber.ErrBadRequest
	}

	findProduct := new(entity.Product)
	findProduct.ID = newCartItem.ProductID
	if err := c.ProductRepository.FindFirst(tx, findProduct); err != nil {
		c.Log.Warnf("Failed to find product by id in the database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newCartItem.Quantity = newCartItem.Quantity + request.Quantity
	if newCartItem.Quantity < 1 {
		newCartItem.Quantity = 0
		// hapus dari cart item
		deleteCartItem := new(entity.CartItem)
		deleteCartItem.ID = request.CartItemID
		if err := c.CartItemRepository.Delete(tx, deleteCartItem); err != nil {
			c.Log.Warnf("Failed to delete cart item by cart item id in the database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	} else {
		// tambah/kurang
		updateCartItem := new(entity.CartItem)
		updateCartItem.ID = request.CartItemID
		updateQuantity := map[string]any{
			"quantity": newCartItem.Quantity,
		}

		if err := c.CartItemRepository.UpdateCustomColumns(tx, updateCartItem, updateQuantity); err != nil {
			c.Log.Warnf("Failed to update quantity of cart item in the database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	// update stok produk
	findProduct.Stock = findProduct.Stock - request.Quantity
	updateProduct := new(entity.Product)
	updateProduct.ID = findProduct.ID
	updateStock := map[string]any{
		"stock": findProduct.Stock,
	}

	if err := c.ProductRepository.UpdateCustomColumns(tx, updateProduct, updateStock); err != nil {
		c.Log.Warnf("Failed to update stock of product in the database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.CartItemToResponse(newCartItem), nil
}