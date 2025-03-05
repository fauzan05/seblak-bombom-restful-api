package converter

import (
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
	response := &model.OrderResponse{
		ID:              order.ID,
		Invoice:         order.Invoice,
		Amount:          fmt.Sprintf("%.2f", order.Amount),
		DiscountType:    order.DiscountType,
		DiscountValue:   order.DiscountValue,
		TotalDiscount:   order.TotalDiscount,
		UserId:          order.UserId,
		FirstName:       order.FirstName,
		LastName:        order.LastName,
		Email:           order.Email,
		Phone:           order.Phone,
		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   order.PaymentStatus,
		OrderStatus:     order.OrderStatus,
		IsDelivery:      order.IsDelivery,
		DeliveryCost:    fmt.Sprintf("%.2f", order.DeliveryCost),
		CompleteAddress: order.CompleteAddress,
		Note:            order.Note,
		OrderProducts:   *OrderProductsToResponse(&order.OrderProducts),
		CreatedAt:       order.Created_At,
		UpdatedAt:       order.Updated_At,
	}

	if order.XenditTransaction != nil {
		response.XenditTransaction = XenditTransactionToResponse(*order.XenditTransaction)
	}

	return response
}

func OrdersToResponse(orders *[]entity.Order) *[]model.OrderResponse {
	getOrders := make([]model.OrderResponse, len(*orders))
	for i, order := range *orders {
		getOrders[i] = *OrderToResponse(&order)
	}
	return &getOrders
}
