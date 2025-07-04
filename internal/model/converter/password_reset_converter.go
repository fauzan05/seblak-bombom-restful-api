package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func PasswordResetToResponse(passwordReset *entity.PasswordReset) *model.PasswordResetResponse {
	return &model.PasswordResetResponse{
		ID:               passwordReset.ID,
		UserId:           passwordReset.UserId,
		VerificationCode: passwordReset.VerificationCode,
		ExpiresAt:        helper_others.TimeRFC3339(passwordReset.ExpiresAt),
		CreatedAt:        helper_others.TimeRFC3339(passwordReset.CreatedAt),
	}
}
