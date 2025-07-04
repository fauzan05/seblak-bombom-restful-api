package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func CartToResponse(cart *entity.Cart) *model.CartResponse {
	response := &model.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		CreatedAt: helper_others.TimeRFC3339(cart.CreatedAt),
		UpdatedAt: helper_others.TimeRFC3339(cart.UpdatedAt),
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
