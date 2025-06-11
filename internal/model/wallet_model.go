package model

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type WalletResponse struct {
	ID        uint64                    `json:"id"`
	Balance   float32                   `json:"balance"`
	Status    enum_state.WalletStatus   `json:"status"`
	CreatedAt helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt helper_others.TimeRFC3339 `json:"updated_at"`
}

type GetWalletBalance struct {
	ID uint64 `json:"-" validate:"required"`
}

type WithDrawWalletRequest struct {
	UserId uint64  `json:"user_id" validate:"required"`
	Amount float32 `json:"amount" validate:"required"`
	Notes  string  `json:"notes"`
}

type UpdateWalletBalance struct {
	ID      uint64  `json:"-" validate:"required"`
	Balance float32 `json:"balance" validate:"required"`
}

type SuspendWallet struct {
	IDs []uint64 `json:"-" validate:"required"`
}
