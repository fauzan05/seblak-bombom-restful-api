package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type CreateXenditTransaction struct {
	OrderId  uint64           `json:"order_id" validate:"required"`
	Lang     helper.Languange `json:"-"`
	TimeZone time.Location    `json:"-"`
}

type CreateXenditQRCode struct {
	ReferenceId string             `json:"reference_id" validate:"required"`
	Type        string             `json:"type" validate:"required"`
	Currency    string             `json:"currency" validate:"required"`
	Amount      float64            `json:"amount" validate:"required"`
	ExpiresAt   helper.TimeRFC3339 `json:"expires_at" validate:"required"`
}

type GetXenditQRCodeTransaction struct {
	OrderId uint64 `json:"-" validate:"required"`
}

type XenditTransactionResponse struct {
	ID              string             `json:"id"`
	ReferenceId     string             `json:"reference_id"`
	OrderId         uint64             `json:"order_id"`
	Amount          float64            `json:"amount"`
	Currency        string             `json:"currency"`
	PaymentMethod   string             `json:"payment_method"`
	PaymentMethodId string             `json:"payment_method_id"`
	ChannelCode     string             `json:"channel_code"`
	QrString        string             `json:"qr_string,omitempty"`
	Status          string             `json:"status"`
	Description     string             `json:"description"`
	FailureCode     string             `json:"failure_code"`
	Metadata        []byte             `json:"metadata"`
	ExpiresAt       helper.TimeRFC3339 `json:"expires_at"`
	CreatedAt       helper.TimeRFC3339 `json:"created_at"`
	UpdatedAt       helper.TimeRFC3339 `json:"updated_at"`
}

type XenditGetPaymentRequestCallbackStatus struct {
	Data struct {
		PaymentMethod struct {
			ID string `json:"id"`
		} `json:"payment_method" validate:"required"`
		Status    string             `json:"status" validate:"required"`
		Metadata  map[string]any     `json:"metadata"`
		UpdatedAt helper.TimeRFC3339 `json:"updated" validate:"required"`
	} `json:"data" validate:"required"`
	Lang            helper.Languange `json:"-"`
	TimeZone        time.Location    `json:"-"`
	BaseFrontEndURL string           `json:"-"`
}

type XenditGetPayoutRequestCallbackStatus struct {
	Data struct {
		PayoutId  string             `json:"id"`
		Status    string             `json:"status" validate:"required"`
		Amount    float32            `json:"amount" validate:"required"`
		UpdatedAt helper.TimeRFC3339 `json:"updated" validate:"required"`
	} `json:"data" validate:"required"`
}
