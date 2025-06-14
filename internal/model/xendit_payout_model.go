package model

import (
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type CreateXenditPayout struct {
	ChannelCode       string  `json:"channel_code" validate:"required"`
	UserId            uint64  `json:"-" validate:"required"`
	AccountNumber     string  `json:"account_number" validate:"required"`
	AccountHolderName string  `json:"account_holder_name" validate:"required"`
	Amount            float32 `json:"amount" validate:"required"`
	Description       string  `json:"description" validate:"max=100"`
	Currency          string  `json:"currency" validate:"required,max=10"`
}

type GetWithdrawableBalanceResponse struct {
	WithdrawableBalance float32 `json:"withdrawable_balance"`
	TotalWalletBalance  float32 `json:"total_wallet_balance"`
}

type XenditPayoutResponse struct {
	ID                string                    `json:"id"`
	UserId            uint64                    `json:"user_id"`
	BusinessId        string                    `json:"business_id"`
	ReferenceId       string                    `json:"reference_id"`
	Amount            float32                   `json:"amount"`
	Currency          string                    `json:"currency"`
	Description       string                    `json:"description"`
	ChannelCode       string                    `json:"channel_code"`
	AccountNumber     string                    `json:"account_number"`
	AccountHolderName string                    `json:"account_holder_name"`
	Status            string                    `json:"status"`
	CreatedAt         helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt         helper_others.TimeRFC3339 `json:"updated_at"`
	EstimatedArrival  helper_others.TimeRFC3339 `json:"estimated_arrival"`
	User              *UserResponse             `json:"user,omitempty"`
}

type CancelXenditPayout struct {
	PayoutId string `json:"-" validate:"required"`
}

type GetPayoutById struct {
	PayoutId string `json:"-" validate:"required"`
}
