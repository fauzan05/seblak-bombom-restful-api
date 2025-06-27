package model

import (
	"mime/multipart"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"time"
)

type UserResponse struct {
	ID          uint64                    `json:"id"`
	FirstName   string                    `json:"first_name"`
	LastName    string                    `json:"last_name"`
	Email       string                    `json:"email"`
	Phone       string                    `json:"phone"`
	Addresses   []AddressResponse         `json:"addresses,omitempty"`
	Role        enum_state.Role           `json:"role"`
	Wallet      WalletResponse            `json:"wallet"`
	Cart        CartResponse              `json:"cart"`
	UserProfile string                    `json:"user_profile"`
	CreatedAt   helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt   helper_others.TimeRFC3339 `json:"updated_at"`
}

type RegisterUserRequest struct {
	FirstName string               `json:"first_name" validate:"required,max=100"`
	LastName  string               `json:"last_name"`
	Email     string               `json:"email" validate:"required,max=100"`
	Phone     string               `json:"phone" validate:"required,max=50"`
	Password  string               `json:"password" validate:"required"`
	Role      enum_state.Role      `json:"role" validate:"required"`
	Lang      enum_state.Languange `json:"-"`
	TimeZone  time.Location        `json:"-"`
}

type VerifyEmailRegisterRequest struct {
	VerificationToken string               `json:"verification_token" validate:"required,max=100"`
	Lang              enum_state.Languange `json:"-"`
	TimeZone          time.Location        `json:"-"`
	BaseFrontEndURL   string               `json:"-"`
}

type VerifyUserRequest struct {
	Token string `json:"token" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName   string                `json:"first_name" validate:"required,max=100"`
	LastName    string                `json:"last_name" validate:"required,max=100"`
	Email       string                `json:"email" validate:"required,max=100"`
	Phone       string                `json:"phone" validate:"required,max=50"`
	UserProfile *multipart.FileHeader `json:"user_profile"`
}

type UpdateUserPasswordRequest struct {
	OldPassword        string `json:"old_password" validate:"required,max=100"`
	NewPassword        string `json:"new_password" validate:"required,min=8,max=100"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,eqfield=NewPassword"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Remember bool   `json:"remember"`
}

type UserTokenResponse struct {
	Token      string                    `json:"token"`
	ExpiryDate time.Time                 `json:"expiry_date"`
	CreatedAt  helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt  helper_others.TimeRFC3339 `json:"updated_at"`
}

type GetUserByTokenRequest struct {
	Token string `validate:"required"`
}

type DeleteCurrentUserRequest struct {
	OldPassword string               `json:"old_password" validate:"required"`
	Lang        enum_state.Languange `json:"-"`
	TimeZone    time.Location        `json:"-"`
}

type CreateForgotPassword struct {
	Email string               `json:"email" validate:"required"`
	Lang  enum_state.Languange `json:"-"`
}

type ValidateForgotPassword struct {
	ID               uint64 `json:"-" validate:"required"`
	VerificationCode int    `json:"verification_code" validate:"required"`
}

type PasswordResetRequest struct {
	ID                 uint64               `json:"-" validate:"required"`
	VerificationCode   int                  `json:"verification_code" validate:"required"`
	NewPassword        string               `json:"new_password" validate:"required,min=8,max=100"`
	NewPasswordConfirm string               `json:"new_password_confirm" validate:"required,eqfield=NewPassword"`
	Lang               enum_state.Languange `json:"-"`
	TimeZone           time.Location        `json:"-"`
}

type PasswordResetResponse struct {
	ID               uint64                    `json:"id"`
	UserId           uint64                    `json:"user_id"`
	VerificationCode int                       `json:"verification_code"`
	ExpiresAt        helper_others.TimeRFC3339 `json:"expires_at"`
	CreatedAt        helper_others.TimeRFC3339 `json:"created_at"`
}
