package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	addresses := make([]model.AddressResponse, len(user.Addresses)) // Inisialisasi slice untuk menyimpan semua alamat
	// Konversi setiap alamat ke AddressResponse
	for i, address := range user.Addresses {
		addresses[i] = *AddressToResponse(&address)
	}
	response := &model.UserResponse{
		ID:          user.ID,
		FirstName:   user.Name.FirstName,
		LastName:    user.Name.LastName,
		Email:       user.Email,
		Phone:       user.Phone,
		Addresses:   addresses,
		Role:        user.Role,
		UserProfile: user.UserProfile,
		CreatedAt:   helper_others.TimeRFC3339(user.CreatedAt),
		UpdatedAt:   helper_others.TimeRFC3339(user.UpdatedAt),
	}

	if user.Wallet != nil {
		response.Wallet = *WalletToResponse(user.Wallet)
	}

	if user.Cart != nil {
		response.Cart = *CartToResponse(user.Cart)
	}

	return response
}

func UserTokenToResponse(token *entity.Token) *model.UserTokenResponse {
	return &model.UserTokenResponse{
		Token:      token.Token,
		ExpiryDate: token.ExpiryDate,
		CreatedAt:  helper_others.TimeRFC3339(token.CreatedAt),
		UpdatedAt:  helper_others.TimeRFC3339(token.UpdatedAt),
	}
}
