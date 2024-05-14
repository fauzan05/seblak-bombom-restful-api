package model

type AddressResponse struct {
	ID              uint64  `json:"id,omitempty"`
	Regency         string  `json:"regency,omitempty"`
	Subdistrict     string  `json:"subdistrict,omitempty"`
	CompleteAddress string  `json:"complete_address,omitempty"`
	Longitude       float64 `json:"longitude,omitempty"`
	Latitude        float64 `json:"latitude,omitempty"`
	IsMain          bool    `json:"is_main"`
	CreatedAt       string  `json:"created_at,omitempty"`
	UpdatedAt       string  `json:"updated_at,omitempty"`
}

type AddressCreateRequest struct {
	Regency         string  `json:"regency" validate:"required,max=100"`
	Subdistrict     string  `json:"subdistrict" validate:"required,max=100"`
	CompleteAddress string  `json:"complete_address" validate:"required"`
	Longitude       float64 `json:"longitude" validate:"required"`
	Latitude        float64 `json:"latitude" validate:"required"`
	IsMain          bool    `json:"is_main"`
}

type DeleteAddressRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateAddressRequest struct {
	ID              uint64  `json:"-" validate:"required"`
	UserId          uint64  `json:"-" validate:"required"`
	Regency         string  `json:"regency" validate:"required,max=100"`
	Subdistrict     string  `json:"subdistrict" validate:"required,max=100"`
	CompleteAddress string  `json:"complete_address" validate:"required"`
	Longitude       float64 `json:"longitude" validate:"required"`
	Latitude        float64 `json:"latitude" validate:"required"`
	IsMain          bool    `json:"is_main"`
}

type GetAddressRequest struct {
	ID uint64 `json:"-" validate:"required"`
}
