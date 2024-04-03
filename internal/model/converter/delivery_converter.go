package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func DeliveryToResponse(delivery *entity.Delivery) *model.DeliveryResponse {
	return &model.DeliveryResponse{
		ID:        delivery.ID,
		Cost:      delivery.Cost,
		Distance:  delivery.Distance,
		CreatedAt: delivery.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: delivery.Updated_At.Format("2006-01-02 15:04:05"),
	}
}