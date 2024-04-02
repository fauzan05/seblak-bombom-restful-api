package converter

import (
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
		Price:              order.Price,
		Quantity:           order.Quantity,
		Amount:             order.Amount,
		DiscountValue:      order.DiscountValue,
		DiscountType:       order.DiscountType,
		UserId:             order.UserId,
		FirstName:          order.FirstName,
		LastName:           order.LastName,
		Email:              order.Email,
		Phone:              order.Phone,
		PaymentMethod:      order.PaymentMethod,
		PaymentStatus:      order.PaymentStatus,
		DeliveryStatus:     order.DeliveryStatus,
		IsDelivery:         order.IsDelivery,
		DeliveryCost:       order.DeliveryCost,
		CategoryName:       order.CategoryName,
		CompleteAddress:    order.CompleteAddress,
		GoogleMapLink:      order.GoogleMapLink,
		CreatedAt:          order.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt:          order.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

// func OrdersToResponse(categories *[]entity.Category) *[]model.CategoryResponse {
// 	getCategories := make([]model.CategoryResponse, len(*categories))
// 	for i , category := range *categories {
// 		getCategories[i] = *CategoryToResponse(&category)
// 	}
// 	return &getCategories
// }
