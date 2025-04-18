package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
	response := &model.OrderResponse{
		ID:                order.ID,
		Invoice:           order.Invoice,
		DiscountType:      order.DiscountType,
		DiscountValue:     order.DiscountValue,
		TotalDiscount:     order.TotalDiscount,
		UserId:            order.UserId,
		FirstName:         order.FirstName,
		LastName:          order.LastName,
		Email:             order.Email,
		Phone:             order.Phone,
		PaymentGateway:    order.PaymentGateway,
		PaymentMethod:     order.PaymentMethod,
		PaymentStatus:     order.PaymentStatus,
		ChannelCode:       order.ChannelCode,
		OrderStatus:       order.OrderStatus,
		IsDelivery:        order.IsDelivery,
		DeliveryCost:      order.DeliveryCost,
		CompleteAddress:   order.CompleteAddress,
		Note:              order.Note,
		TotalProductPrice: order.TotalProductPrice,
		TotalFinalPrice:   order.TotalFinalPrice,
		CreatedAt:         helper.TimeRFC3339(order.CreatedAt),
		UpdatedAt:         helper.TimeRFC3339(order.UpdatedAt),
		OrderProducts:     *OrderProductsToResponse(&order.OrderProducts),
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
