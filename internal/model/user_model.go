package model

import "time"

type UserResponse struct {
	ID        uint64          `json:"id,omitempty"`
	FirstName string          `json:"first_name,omitempty"`
	LastName  string          `json:"last_name,omitempty"`
	Email     string          `json:"email,omitempty"`
	Phone     string          `json:"phone,omitempty"`
	Addresses []AddressResponse `json:"addresses,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
}

type RegisterUserRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,max=100"`
	Phone     string `json:"phone" validate:"required,max=50"`
	Password  string `json:"password" validate:"required,max=100"`
}

type VerifyUserRequest struct {
	Token string `json:"token" validate:"required"`
}

type LoginUserRequst struct {
	Email    string `json:"email" validate:"required,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type UserTokenResponse struct {
	Token      string    `json:"token,omitempty"`
	ExpiryDate time.Time `json:"expiry_date,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

type GetUserByTokenRequest struct {
	Token string ` validate:"required"`
}
