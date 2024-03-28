package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		FirstName: user.Name.FirstName,
		LastName:  user.Name.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.Created_At,
		UpdatedAt: user.Updated_At,
	}
}

func UserTokenToResponse(user *entity.User) *model.UserTokenResponse {
	return &model.UserTokenResponse{
		Token:      user.Token.Token,
		ExpiryDate: user.Token.ExpiryDate,
		CreatedAt:  user.Token.Created_At,
		UpdatedAt:  user.Token.Updated_At,
	}
}
