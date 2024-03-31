package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	addresses := make([]model.AddressResponse, len(user.Addresses)) // Inisialisasi slice untuk menyimpan semua alamat
	// Konversi setiap alamat ke AddressResponse
	for i, address := range user.Addresses {
		addresses[i] = *AddressToResponse(&address)
	}
	return &model.UserResponse{
		ID:        user.ID,
		FirstName: user.Name.FirstName,
		LastName:  user.Name.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Addresses: addresses,
		Role:      user.Role,
		CreatedAt: user.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func UserTokenToResponse(token *entity.Token) *model.UserTokenResponse {
	return &model.UserTokenResponse{
		Token:      token.Token,
		ExpiryDate: token.ExpiryDate,
		CreatedAt:  token.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt:  token.Updated_At.Format("2006-01-02 15:04:05"),
	}
}
