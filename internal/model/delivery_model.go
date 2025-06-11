package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type DeliveryResponse struct {
	ID        uint64                    `json:"id"`
	City      string                    `json:"city"`
	District  string                    `json:"district"`
	Village   string                    `json:"village"`
	Hamlet    string                    `json:"hamlet"`
	Cost      float32                   `json:"cost"`
	CreatedAt helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt helper_others.TimeRFC3339 `json:"updated_at"`
}

type CreateDeliveryRequest struct {
	City     string  `json:"city" validate:"required"`
	District string  `json:"district" validate:"required"`
	Village  string  `json:"village" validate:"required"`
	Hamlet   string  `json:"hamlet" validate:"required"`
	Cost     float32 `json:"cost" validate:"required"`
}

type UpdateDeliveryRequest struct {
	ID       uint64  `json:"-" validate:"required"`
	City     string  `json:"city" validate:"required"`
	District string  `json:"district" validate:"required"`
	Village  string  `json:"village" validate:"required"`
	Hamlet   string  `json:"hamlet" validate:"required"`
	Cost     float32 `json:"cost" validate:"required"`
}

type DeleteDeliveryRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}
