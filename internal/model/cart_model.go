package model

type CartResponse struct {
	ID         uint64  `json:"id,omitempty"`
	UserID     uint64  `json:"user_id,omitempty"`
	ProductID  uint64  `json:"product_id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Quantity   int     `json:"quantity,omitempty"`
	Price      float32 `json:"price,omitempty"`
	TotalPrice float32 `json:"total_price,omitempty"`
	CreatedAt  string  `json:"created_at,omitempty"`
	UpdatedAt  string  `json:"updated_at,omitempty"`
}

type CreateCartRequest struct {
	UserID     uint64  `json:"user_id" validate:"required"`
	ProductID  uint64  `json:"product_id" validate:"required"`
	Quantity   int     `json:"quantity" validate:"required"`
}

type GetAllCartByCurrentUserRequest struct {
	UserID uint64 `json:"-" validate:"required"`
}

type UpdateCartRequest struct {
	ID       uint64 `json:"-" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type DeleteCartRequest struct {
	ID uint64 `json:"-" validate:"required"`
}
