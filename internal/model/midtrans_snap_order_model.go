package model

type MidtransSnapOrderResponse struct {
	ID          uint64 `json:"id,omitempty"`
	OrderId     uint64 `json:"order_id,omitempty"`
	Token       string `json:"token,omitempty"`
	RedirectUrl string `json:"redirect_url,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type CreateMidtransSnapOrderRequest struct {
	OrderId     uint64 `json:"order_id" validate:"required"`
}
