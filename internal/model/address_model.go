package model

import "seblak-bombom-restful-api/internal/helper/helper_others"

type AddressResponse struct {
	ID              uint64                    `json:"id"`
	CompleteAddress string                    `json:"complete_address"`
	GoogleMapsLink  string                    `json:"google_maps_link"`
	IsMain          bool                      `json:"is_main"`
	Delivery        DeliveryResponse          `json:"delivery"`
	CreatedAt       helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt       helper_others.TimeRFC3339 `json:"updated_at"`
}

type AddressCreateRequest struct {
	DeliveryId      uint64 `json:"delivery_id" validate:"required"`
	CompleteAddress string `json:"complete_address" validate:"required"`
	GoogleMapsLink  string `json:"google_maps_link"`
	IsMain          bool   `json:"is_main"`
}

type DeleteAddressRequest struct {
	IDs []uint64 `json:"-" validate:"required"`
}

type UpdateAddressRequest struct {
	ID              uint64 `json:"-" validate:"required"`
	UserId          uint64 `json:"-" validate:"required"`
	DeliveryId      uint64 `json:"delivery_id" validate:"required"`
	CompleteAddress string `json:"complete_address" validate:"required"`
	GoogleMapsLink  string `json:"google_maps_link"`
	IsMain          bool   `json:"is_main"`
}

type GetAddressRequest struct {
	ID uint64 `json:"-" validate:"required"`
}
