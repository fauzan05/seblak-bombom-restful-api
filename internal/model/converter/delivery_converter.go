package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func DeliveryToResponse(delivery *entity.Delivery) *model.DeliveryResponse {
	return &model.DeliveryResponse{
		ID:        delivery.ID,
		Cost:      delivery.Cost,
		City:      delivery.City,
		District:  delivery.District,
		Village:   delivery.Village,
		Hamlet:    delivery.Hamlet,
		CreatedAt: delivery.CreatedAt,
		UpdatedAt: delivery.UpdatedAt,
	}
}

func DeliveriesToResponse(deliveries *[]entity.Delivery) *[]model.DeliveryResponse {
	getDeliveries := make([]model.DeliveryResponse, len(*deliveries))
	for i, delivery := range *deliveries {
		getDeliveries[i] = *DeliveryToResponse(&delivery)
	}
	return &getDeliveries
}