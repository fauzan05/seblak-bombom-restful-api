package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func CartToResponse(cart *entity.Cart) *model.CartResponse {
	response := &model.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: cart.Created_At,
		UpdatedAt: cart.Updated_At,
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
