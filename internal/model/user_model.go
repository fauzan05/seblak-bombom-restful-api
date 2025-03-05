package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type UserResponse struct {
	ID        uint64            `json:"id"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Email     string            `json:"email"`
	Phone     string            `json:"phone"`
	Addresses []AddressResponse `json:"addresses,omitempty"`
	Role      helper.Role       `json:"role"`
	Wallet    WalletResponse    `json:"wallet"`
	Cart      CartResponse      `json:"cart"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type RegisterUserRequest struct {
	FirstName string      `json:"first_name" validate:"required,max=100"`
	LastName  string      `json:"last_name" validate:"required,max=100"`
	Email     string      `json:"email" validate:"required,max=100"`
	Phone     string      `json:"phone" validate:"required,max=50"`
	Password  string      `json:"password" validate:"required,max=100"`
	Role      helper.Role `json:"role" validate:"required"`
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
	NewPassword        string `json:"new_password" validate:"required,max=100"`
	NewPasswordConfirm string `json:"new_password_confirm" validate:"required,max=100,eqfield=NewPassword"`
}

type LoginUserRequst struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type UserTokenResponse struct {
	Token      string    `json:"token"`
	ExpiryDate time.Time `json:"expiry_date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetUserByTokenRequest struct {
	Token string `validate:"required"`
}

type DeleteCurrentUserRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
}
