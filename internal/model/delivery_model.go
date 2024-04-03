package model

type DeliveryResponse struct {
	ID        uint64 `json:"id,omitempty"`
	Cost      int    `json:"cost,omitempty"`
	Distance  int    `json:"distance,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateDeliveryRequest struct {
	Cost     int    `json:"cost" validate:"required"`
	Distance int    `json:"distance" validate:"required"`
}

type UpdateDeliveryRequest struct {
	ID       uint64 `json:"-" validate:"required"`
	Cost     int    `json:"cost" validate:"required"`
	Distance int    `json:"distance" validate:"required"`
}

type DeleteDeliveryRequest struct {
	ID uint64 `json:"-" validate:"required"`
}
