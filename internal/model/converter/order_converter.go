package converter

import (
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
	return &model.OrderResponse{
		ID:                 order.ID,
		Invoice:            order.Invoice,
		ProductId:          order.ProductId,
		ProductName:        order.ProductName,
		ProductDescription: order.ProductDescription,
		Price:              fmt.Sprintf("%.2f", order.Price),
		Quantity:           order.Quantity,
		Amount:             fmt.Sprintf("%.2f", order.Amount),
		DiscountType:       order.DiscountType,
		DiscountValue:      order.DiscountValue,
		UserId:             order.UserId,
		FirstName:          order.FirstName,
		LastName:           order.LastName,
		Email:              order.Email,
		Phone:              order.Phone,
		PaymentMethod:      order.PaymentMethod,
		PaymentStatus:      order.PaymentStatus,
		DeliveryStatus:     order.DeliveryStatus,
		IsDelivery:         order.IsDelivery,
		DeliveryCost:       fmt.Sprintf("%.2f", order.DeliveryCost),
		CategoryName:       order.CategoryName,
		CompleteAddress:    order.CompleteAddress,
		GoogleMapLink:      order.GoogleMapLink,
		Distance:           order.Distance,
		CreatedAt:          order.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt:          order.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func OrdersToResponse(orders *[]entity.Order) *[]model.OrderResponse {
	getOrders := make([]model.OrderResponse, len(*orders))
	for i, order := range *orders {
		getOrders[i] = *OrderToResponse(&order)
	}
	return &getOrders
}
