package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type CreatePayoutRequest struct {
	Amount              float32             `json:"amount" validate:"required"`
	Currency            string              `json:"currency"`
	Method              helper.PayoutMethod `json:"method"`
	Notes               string              `json:"notes"`
	UserId              uint64              `json:"user_id"`
	XenditPayoutRequest *CreateXenditPayout `json:"xendit_payout_request" validate:"required_if=Method 1"`
}

type PayoutResponse struct {
	ID                uint64                `json:"id"`
	XenditPayoutId    string                `json:"xendit_payout_id"`
	Amount            float32               `json:"amount"`
	Currency          string                `json:"currency"`
	Method            helper.PayoutMethod   `json:"method"`
	Status            helper.PayoutStatus   `json:"status"`
	Notes             string                `json:"notes"`
	CancellationNotes string                `json:"cancellation_notes"`
	XenditPayout      *XenditPayoutResponse `json:"xendit_payout"`
	CreatedAt         time.Time             `json:"created_at,omitempty"`
	UpdatedAt         time.Time             `json:"updated_at,omitempty"`
}
