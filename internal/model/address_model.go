package model

type AddressResponse struct {
	ID              uint64    `json:"id,omitempty"`
	Regency         string    `json:"regency,omitempty"`
	Subdistrict     string    `json:"subdistrict,omitempty"`
	CompleteAddress string    `json:"complete_address,omitempty"`
	GoogleMapLink   string    `json:"google_map_link,omitempty"`
	IsMain          bool      `json:"is_main,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}

type AddressCreateRequest struct {
	Regency         string `json:"regency" validate:"required,max=100"`
	Subdistrict     string `json:"subdistrict" validate:"required,max=100"`
	CompleteAddress string `json:"complete_address" validate:"required"`
	GoogleMapLink   string `json:"google_map_link" validate:"required"`
	IsMain          bool   `json:"is_main" validate:"required"`
}

type AddressDeleteRequest struct {
	IdAddress uint64 `json:"-" validate:"required"`
}
