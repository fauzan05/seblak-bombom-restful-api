package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID: user.ID,
		FirstName: user.Name.FirstName,
		LastName: user.Name.LastName,
		Email: user.Email,
		Phone: user.Phone,
		CreatedAt: user.Created_At,
		UpdatedAt: user.Updated_At,
	}
}