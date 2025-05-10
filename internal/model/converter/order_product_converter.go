package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
)

func OrderProductToResponse(orderProduct *entity.OrderProduct) *model.OrderProductResponse {
	return &model.OrderProductResponse{
		ID:          orderProduct.ID,
		OrderId:     orderProduct.OrderId,
		ProductId:   orderProduct.ProductId,
		ProductName: orderProduct.ProductName,
		Category:    orderProduct.Category,
		Price:       orderProduct.Price,
		Quantity:    orderProduct.Quantity,
		CreatedAt:   helper.TimeRFC3339(orderProduct.CreatedAt),
		UpdatedAt:   helper.TimeRFC3339(orderProduct.UpdatedAt),
	}
}

func OrderProductsToResponse(orderProducts *[]entity.OrderProduct) *[]model.OrderProductResponse {
	getOrderProducts := make([]model.OrderProductResponse, len(*orderProducts))
	for i, orderProduct := range *orderProducts {
		getOrderProducts[i] = *OrderProductToResponse(&orderProduct)
	}
	return &getOrderProducts
}
