package model

type ProductReviewResponse struct {
	ID        uint64 `json:"id,omitempty"`
	ProductId uint64 `json:"product_id,omitempty"`
	UserId    uint64 `json:"user_id,omitempty"`
	Rate      int    `json:"rate,omitempty"`
	Comment   string `json:"comment,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type CreateProductReviewRequest struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	UserId    uint64 `json:"-" validate:"required"` // via token
	Rate      int    `json:"rate" validate:"required"`
	Comment   string `json:"comment" validate:"required"`
}
