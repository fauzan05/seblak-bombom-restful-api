package model

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type CreatePayoutRequest struct {
	Amount              float32                 `json:"amount" validate:"required"`
	Currency            string                  `json:"currency"`
	Method              enum_state.PayoutMethod `json:"method"`
	Notes               string                  `json:"notes"`
	UserId              uint64                  `json:"user_id"`
	XenditPayoutRequest *CreateXenditPayout     `json:"xendit_payout_request" validate:"required_if=Method 1"`
}

type PayoutResponse struct {
	ID             uint64                    `json:"id"`
	XenditPayoutId string                    `json:"xendit_payout_id"`
	Amount         float32                   `json:"amount"`
	Currency       string                    `json:"currency"`
	Method         enum_state.PayoutMethod   `json:"method"`
	Status         enum_state.PayoutStatus   `json:"status"`
	Notes          string                    `json:"notes"`
	XenditPayout   *XenditPayoutResponse     `json:"xendit_payout"`
	CreatedAt      helper_others.TimeRFC3339 `json:"created_at,omitempty"`
	UpdatedAt      helper_others.TimeRFC3339 `json:"updated_at,omitempty"`
}
