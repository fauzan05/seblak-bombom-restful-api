package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func CartToResponse(cart *entity.Cart) *model.CartResponse {
	return &model.CartResponse{
		ID:         cart.ID,
		UserID:     cart.UserID,
		ProductID:  cart.ProductID,
		Name:       cart.Name,
		Quantity:   cart.Quantity,
		Price:      cart.Price,
		TotalPrice: cart.TotalPrice,
		Stock:      cart.Stock,
		CreatedAt:  cart.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt:  cart.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func CartsToResponse(carts *[]entity.Cart) *[]model.CartResponse {
	getCarts := make([]model.CartResponse, len(*carts))
	for i, cart := range *carts {
		getCarts[i] = *CartToResponse(&cart)
	}
	return &getCarts
}
