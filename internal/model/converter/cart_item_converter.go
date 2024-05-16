package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func CartItemToResponse(cartItem *entity.CartItem) *model.CartItemResponse {
	return &model.CartItemResponse{
		ID:        cartItem.ID,
		CartId:    cartItem.CartId,
		Product:   *ProductToResponse(cartItem.Product),
		Quantity:  cartItem.Quantity,
		CreatedAt: cartItem.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: cartItem.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func CartItemsToResponse(cartItems *[]entity.CartItem) *[]model.CartItemResponse {
	getCartItems := make([]model.CartItemResponse, len(*cartItems))
	for i, cartItem := range *cartItems {
		getCartItems[i] = *CartItemToResponse(&cartItem)
	}
	return &getCartItems
}
