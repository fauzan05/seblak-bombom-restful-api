package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"

	// "seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MidtransSnapOrderUseCase struct {
	Log                         *logrus.Logger
	DB                          *gorm.DB
	Validate                    *validator.Validate
	SnapClient                  *snap.Client
	CoreAPIClient               *coreapi.Client
	OrderRepository             *repository.OrderRepository
	ProductRepository           *repository.ProductRepository
	MidtransSnapOrderRepository *repository.MidtransSnapOrderRepository
}

func NewMidtransSnapOrderUseCase(log *logrus.Logger, validate *validator.Validate, orderRepository *repository.OrderRepository,
	snapClient *snap.Client, coreAPIClient *coreapi.Client, db *gorm.DB, midtransSnapOrderRepository *repository.MidtransSnapOrderRepository,
	productRepository *repository.ProductRepository) *MidtransSnapOrderUseCase {
	return &MidtransSnapOrderUseCase{
		Log:                         log,
		Validate:                    validate,
		OrderRepository:             orderRepository,
		ProductRepository:           productRepository,
		SnapClient:                  snapClient,
		CoreAPIClient:               coreAPIClient,
		DB:                          db,
		MidtransSnapOrderRepository: midtransSnapOrderRepository,
	}
}

func (c *MidtransSnapOrderUseCase) Add(ctx context.Context, request *model.CreateMidtransSnapOrderRequest) (*model.MidtransSnapOrderResponse, error) {
	tx := c.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newSnapClient := snap.Client{}
	newSnapClient.New(c.SnapClient.ServerKey, c.SnapClient.Env)

	// temukan data order berdasarkan invoice dari request
	selectedOrder := new(entity.Order)
	selectedOrder.ID = request.OrderId
	if err := c.OrderRepository.FindWithPreloads(tx, selectedOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	var midtransItemDetails []midtrans.ItemDetails
	for _, product := range selectedOrder.OrderProducts {
		newMidtransItemDetail := midtrans.ItemDetails{
			ID:    strconv.Itoa(int(product.ID)),
			Qty:   int32(product.Quantity),
			Price: int64(product.Price),
			Name:  product.ProductName,
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemDetail)
	}
	// cek apakah ada biaya pengiriman, jika ada maka tambahkan ke midtransItemDetails agar value-nya sama
	if selectedOrder.DeliveryCost > 0 {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID:    "1",
			Qty:   1,
			Price: int64(selectedOrder.DeliveryCost),
			Name:  "Delivery Cost",
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}
	// cek juga apakah terdapat diskon, jika ada maka tambahkan ke ItemDetails
	if selectedOrder.TotalDiscount > 0 {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID:    "2",
			Qty:   1,
			Price: -int64(selectedOrder.TotalDiscount),
			Name:  "Discount",
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}

	orderId := int(selectedOrder.ID)
	midtransRequest := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(orderId),
			GrossAmt: int64(selectedOrder.Amount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: selectedOrder.FirstName,
			LName: selectedOrder.LastName,
			Email: selectedOrder.Email,
			Phone: selectedOrder.Phone,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items:           &midtransItemDetails,
	}

	snapResponse, snapErr := newSnapClient.CreateTransaction(midtransRequest)
	if snapErr != nil {
		c.Log.Warnf("Failed to create new midtrans transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	newMidtransSnapOrder := new(entity.MidtransSnapOrder)
	newMidtransSnapOrder.OrderId = request.OrderId
	newMidtransSnapOrder.Token = snapResponse.Token
	newMidtransSnapOrder.RedirectUrl = snapResponse.RedirectURL
	if err := c.MidtransSnapOrderRepository.Create(tx, newMidtransSnapOrder); err != nil {
		c.Log.Warnf("Failed to create new midtrans snap order : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		c.Log.Warnf("Failed to commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.MidtransSnapOrderToResponse(newMidtransSnapOrder), nil
}

func (c *MidtransSnapOrderUseCase) Get(ctx context.Context, request *model.GetMidtransSnapOrderRequest) (*model.OrderResponse, error) {
	tx := c.DB.WithContext(ctx)

	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	selectedOrder := new(entity.Order)
	selectedOrder.ID = request.OrderId
	if err := c.OrderRepository.FindWithPreloads(tx, selectedOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// untuk mengecek notifikasi harus menggunakan core api
	newCoreAPIClient := coreapi.Client{}
	newCoreAPIClient.New(c.CoreAPIClient.ServerKey, c.CoreAPIClient.Env)
	orderIdIntConversion := int(request.OrderId)
	orderIdStringConversion := strconv.Itoa(orderIdIntConversion)
	transactionStatusResponse, e := newCoreAPIClient.CheckTransaction(orderIdStringConversion)
	if e != nil {
		c.Log.Warnf("Failed to check the transaction by order id : %+v", err)
		return nil, fiber.ErrBadRequest
	} else {
		if transactionStatusResponse != nil {
			// 5. Do set transaction status based on response from check transaction status
			// if transactionStatusResponse.TransactionStatus == "capture" {
			// 	if transactionStatusResponse.FraudStatus == "challenge" {
			// 		// TODO set transaction status on your database to 'challenge'
			// 		// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			// 		selectedOrder.PaymentStatus = helper.PENDING_PAYMENT
			// 		if err := c.OrderRepository.Update(tx, selectedOrder); err != nil {
			// 			c.Log.Warnf("Failed to update status success order by id : %+v", err)
			// 			return nil, fiber.ErrInternalServerError
			// 		}
			// 	} else if transactionStatusResponse.FraudStatus == "accept" {
			// 		// TODO set transaction status on your database to 'success'
			// 		selectedOrder.PaymentStatus = helper.PAID_PAYMENT
			// 		if err := c.OrderRepository.Update(tx, selectedOrder); err != nil {
			// 			c.Log.Warnf("Failed to update status success order by id : %+v", err)
			// 			return nil, fiber.ErrInternalServerError
			// 		}
			// 	}
			// } else if transactionStatusResponse.TransactionStatus == "settlement" {
			// 	// TODO set transaction status on your databaase to 'success'
			// 	selectedOrder.PaymentStatus = helper.PAID_PAYMENT
			// 	if err := c.OrderRepository.Update(tx, selectedOrder); err != nil {
			// 		c.Log.Warnf("Failed to update status success order by id : %+v", err)
			// 		return nil, fiber.ErrInternalServerError
			// 	}
			// } else if transactionStatusResponse.TransactionStatus == "deny" {
			// 	// TODO you can ignore 'deny', because most of the time it allows payment retries
			// 	// and later can become success
			// 	selectedOrder.PaymentStatus = helper.PENDING_PAYMENT
			// 	if err := c.OrderRepository.Update(tx, selectedOrder); err != nil {
			// 		c.Log.Warnf("Failed to update status pending order by id : %+v", err)
			// 		return nil, fiber.ErrInternalServerError
			// 	}
			// } else if transactionStatusResponse.TransactionStatus == "cancel" || transactionStatusResponse.TransactionStatus == "expire" {
			// 	// TODO set transaction status on your databaase to 'failure'
			// 	for _, orderProduct := range selectedOrder.OrderProducts {
			// 		newProduct := new(entity.Product)
			// 		newProduct.ID = orderProduct.ProductId
			// 		// mencari data terkini dari produk dengan id
			// 		if err := c.ProductRepository.FindById(tx, newProduct); err != nil {
			// 			c.Log.Warnf("Failed to find product by id : %+v", err)
			// 			return nil, fiber.ErrInternalServerError
			// 		}
			// 		// tambahkan/kembalikan stok produk karena transaksinya gagal
			// 		newProduct.Stock += orderProduct.Quantity
			// 		// perbarui stok barang sekarang
			// 		if err := c.ProductRepository.Update(tx, newProduct); err != nil {
			// 			c.Log.Warnf("Failed to update product stock : %+v", err)
			// 			return nil, fiber.ErrInternalServerError
			// 		}
			// 	}
			// 	selectedOrder.PaymentStatus = helper.FAILED_PAYMENT
			// 	if err := c.OrderRepository.Update(tx, selectedOrder); err != nil {
			// 		c.Log.Warnf("Failed to update status failed order by id : %+v", err)
			// 		return nil, fiber.ErrInternalServerError
			// 	}
			// } else if transactionStatusResponse.TransactionStatus == "pending" {
			// 	// TODO set transaction status on your databaase to 'pending' / waiting payment
			// 	selectedOrder.PaymentStatus = helper.PENDING_PAYMENT
			// 	if err := c.OrderRepository.Update(tx, selectedOrder); err != nil {
			// 		c.Log.Warnf("Failed to update status pending order by id : %+v", err)
			// 		return nil, fiber.ErrInternalServerError
			// 	}
			// }
		}
	}

	if err := c.OrderRepository.FindWithJoins(tx, selectedOrder, "MidtransSnapOrder"); err != nil {
		c.Log.Warnf("Failed to find order by id with joins midtrans snap order : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.OrderToResponse(selectedOrder), nil
}
