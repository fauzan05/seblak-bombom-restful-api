package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type WalletResponse struct {
	ID        uint64              `json:"id"`
	Balance   float32             `json:"balance"`
	Status    helper.WalletStatus `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
}

type TopUpWalletBalance struct {
	Amount float32 `json:"balance" validate:"required"`
}

type GetWalletBalance struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateWalletBalance struct {
	ID      uint64  `json:"-" validate:"required"`
	Balance float32 `json:"balance" validate:"required"`
}

type SuspendWallet struct {
	IDs []uint64 `json:"-" validate:"required"`
}
