package model

type DeliveryResponse struct {
	ID        uint64  `json:"id,omitempty"`
	Cost      float32 `json:"cost,omitempty"`
	Distance  float32 `json:"distance,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
}

type CreateDeliveryRequest struct {
	Cost     float32 `json:"cost" validate:"required"`
	Distance float32 `json:"distance" validate:"required"`
}

type UpdateDeliveryRequest struct {
	ID       uint64  `json:"-" validate:"required"`
	Cost     float32 `json:"cost" validate:"required"`
	Distance float32 `json:"distance" validate:"required"`
}

type DeleteDeliveryRequest struct {
	ID uint64 `json:"-" validate:"required"`
}
