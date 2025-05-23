package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type UserResponse struct {
	ID        uint64             `json:"id"`
	FirstName string             `json:"first_name"`
	LastName  string             `json:"last_name"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Addresses []AddressResponse  `json:"addresses,omitempty"`
	Role      helper.Role        `json:"role"`
	Wallet    WalletResponse     `json:"wallet"`
	Cart      CartResponse       `json:"cart"`
	CreatedAt helper.TimeRFC3339 `json:"created_at"`
	UpdatedAt helper.TimeRFC3339 `json:"updated_at"`
}

type RegisterUserRequest struct {
	FirstName string           `json:"first_name" validate:"required,max=100"`
	LastName  string           `json:"last_name" validate:"required,max=100"`
	Email     string           `json:"email" validate:"required,max=100"`
	Phone     string           `json:"phone" validate:"required,max=50"`
	Password  string           `json:"password" validate:"required,min=8,max=100"`
	Role      helper.Role      `json:"role" validate:"required"`
	Lang      helper.Languange `json:"-"`
	TimeZone  time.Location    `json:"-"`
}

type VerifyEmailRegisterRequest struct {
	VerificationToken string           `json:"verification_token" validate:"required,max=100"`
	Lang              helper.Languange `json:"-"`
	TimeZone          time.Location    `json:"-"`
	BaseFrontEndURL   string           `json:"-"`
}

type VerifyUserRequest struct {
	Token string `json:"token" validate:"required"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,max=100"`
	Phone     string `json:"phone" validate:"required,max=50"`
}

type UpdateUserPasswordRequest struct {
	OldPassword        string `json:"old_password" validate:"required,max=100"`
	NewPassword        string `json:"new_password" validate:"required,min=8,max=100"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,eqfield=NewPassword"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type UserTokenResponse struct {
	Token      string             `json:"token"`
	ExpiryDate time.Time          `json:"expiry_date"`
	CreatedAt  helper.TimeRFC3339 `json:"created_at"`
	UpdatedAt  helper.TimeRFC3339 `json:"updated_at"`
}

type GetUserByTokenRequest struct {
	Token string `validate:"required"`
}

type DeleteCurrentUserRequest struct {
	OldPassword string           `json:"old_password" validate:"required"`
	Lang        helper.Languange `json:"-"`
	TimeZone    time.Location    `json:"-"`
}

type CreateForgotPassword struct {
	Email string           `json:"email" validate:"required"`
	Lang  helper.Languange `json:"-"`
}

type ValidateForgotPassword struct {
	ID               uint64 `json:"-" validate:"required"`
	VerificationCode int    `json:"verification_code" validate:"required"`
}

type PasswordResetRequest struct {
	ID                 uint64           `json:"-" validate:"required"`
	VerificationCode   int              `json:"verification_code" validate:"required"`
	NewPassword        string           `json:"new_password" validate:"required,min=8,max=100"`
	NewPasswordConfirm string           `json:"new_password_confirm" validate:"required,eqfield=NewPassword"`
	Lang               helper.Languange `json:"-"`
	TimeZone           time.Location    `json:"-"`
}

type PasswordResetResponse struct {
	ID               uint64             `json:"id"`
	UserId           uint64             `json:"user_id"`
	VerificationCode int                `json:"verification_code"`
	ExpiresAt        helper.TimeRFC3339 `json:"expires_at"`
	CreatedAt        helper.TimeRFC3339 `json:"created_at"`
}
