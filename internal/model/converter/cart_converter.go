package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func CartToResponse(cart *entity.Cart) *model.CartResponse {
	response := &model.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: cart.Updated_At.Format("2006-01-02 15:04:05"),
	}

	if cart.CartItems != nil {
		response.CartItems = *CartItemsToResponse(&cart.CartItems)
	}

	return response
}

func CartsToResponse(carts *[]entity.Cart) *[]model.CartResponse {
	getCarts := make([]model.CartResponse, len(*carts))
	for i, cart := range *carts {
		getCarts[i] = *CartToResponse(&cart)
	}
	return &getCarts
}
