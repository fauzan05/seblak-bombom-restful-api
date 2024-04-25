package usecase

import (
	"context"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"seblak-bombom-restful-api/internal/model/converter"
	"seblak-bombom-restful-api/internal/repository"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MidtransUseCase struct {
	Log *logrus.Logger
	DB *gorm.DB
	Validate *validator.Validate
	OrderRepository *repository.OrderRepository
	SnapClient *snap.Client
}

func NewMidtransUseCase(log *logrus.Logger, validate *validator.Validate, orderRepository *repository.OrderRepository,
	snapClient *snap.Client, db *gorm.DB) *MidtransUseCase {
	return &MidtransUseCase{
		Log: log,
		Validate: validate,
		OrderRepository: orderRepository,
		SnapClient: snapClient,
		DB: db,
	}
}

func (c *MidtransUseCase) Add(ctx context.Context, request *model.CreateSnapRequest) (*model.SnapResponse, error) {
	err := c.Validate.Struct(request)
	if err != nil {
		c.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	newSnapClient := snap.Client{}
	newSnapClient.New(c.SnapClient.ServerKey, midtrans.Sandbox)

	// temukan data order berdasarkan invoice dari request
	selectedOrder := new(entity.Order)
	selectedOrder.ID = request.OrderId
	if err := c.OrderRepository.FindWithPreloads(c.DB, selectedOrder, "OrderProducts"); err != nil {
		c.Log.Warnf("Failed to find order by id : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if selectedOrder.PaymentMethod != helper.ONLINE {
		c.Log.Warnf("This order has an onsite payment method : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	var midtransItemDetails []midtrans.ItemDetails
	for _, product := range selectedOrder.OrderProducts {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID: strconv.Itoa(int(product.ID)),
			Qty: int32(product.Quantity),
			Price: int64(product.Price),
			Name: product.ProductName,
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}
	// cek apakah ada biaya pengiriman, jika ada maka tambahkan ke midtransItemDetails agar value-nya sama
	if selectedOrder.DeliveryCost != 0 {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID: "1",
			Qty: 1,
			Price: int64(selectedOrder.DeliveryCost),
			Name: "Delivery Cost",
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}

	if selectedOrder.TotalDiscount > 0 {
		newMidtransItemsDetails := midtrans.ItemDetails{
			ID: "2",
			Qty: 1,
			Price: -int64(selectedOrder.TotalDiscount),
			Name: "Discount",
		}
		midtransItemDetails = append(midtransItemDetails, newMidtransItemsDetails)
	}
	
	orderId := int(selectedOrder.ID)
	midtransRequest := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID: strconv.Itoa(orderId),
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
		Items: &midtransItemDetails,
	}

	snapResponse, snapErr := newSnapClient.CreateTransaction(midtransRequest)
	if snapErr != nil {
		c.Log.Warnf("Failed to create new midtrans transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.MidtransToResponse(snapResponse), nil
}