package converter

import (
	"fmt"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func OrderProductToResponse(orderProduct *entity.OrderProduct) *model.OrderProductResponse {
	return &model.OrderProductResponse{
		ID: orderProduct.ID,
		OrderId: orderProduct.OrderId,
		ProductId: orderProduct.ProductId,
		ProductName: orderProduct.ProductName,
		Price: fmt.Sprintf("%.2f", orderProduct.Price),
		Quantity: orderProduct.Quantity,
		CreatedAt: orderProduct.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: orderProduct.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func OrderProductsToResponse(orderProducts *[]entity.OrderProduct) *[]model.OrderProductResponse{
	getOrderProducts := make([]model.OrderProductResponse, len(*orderProducts))
	for i, orderProduct := range *orderProducts {
		getOrderProducts[i] = *OrderProductToResponse(&orderProduct)
	}
	return &getOrderProducts
}