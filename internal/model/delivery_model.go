package model

type DeliveryResponse struct {
	ID        uint64  `json:"id"`
	City      string  `json:"city"`
	District  string  `json:"district"`
	Village   string  `json:"village"`
	Hamlet    string  `json:"hamlet"`
	Cost      float32 `json:"cost"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
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
	ID uint64 `json:"-" validate:"required"`
}
