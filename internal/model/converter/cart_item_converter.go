package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func CartItemToResponse(cartItem *entity.CartItem) *model.CartItemResponse {
	return &model.CartItemResponse{
		ID:        cartItem.ID,
		CartId:    cartItem.CartId,
		Product:   *ProductToResponse(cartItem.Product),
		Quantity:  cartItem.Quantity,
		CreatedAt: helper_others.TimeRFC3339(cartItem.CreatedAt),
		UpdatedAt: helper_others.TimeRFC3339(cartItem.UpdatedAt),
	}
}

func CartItemsToResponse(cartItems *[]entity.CartItem) *[]model.CartItemResponse {
	getCartItems := make([]model.CartItemResponse, len(*cartItems))
	for i, cartItem := range *cartItems {
		getCartItems[i] = *CartItemToResponse(&cartItem)
	}
	return &getCartItems
}
