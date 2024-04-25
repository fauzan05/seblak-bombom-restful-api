package model

type SnapResponse struct {
	Token       string `json:"token,omitempty"`
	RedirectUrl string `json:"redirect_url,omitempty"`
}

type CreateSnapRequest struct {
	OrderId uint64 `json:"order_id" validate:"required"`
}
