package model

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type MidtransCoreAPIOrderResponse struct {
	ID                uint64                   `json:"id"`
	StatusCode        string                   `json:"status_code"`
	StatusMessage     string                   `json:"status_message"`
	TransactionId     string                   `json:"transaction_id"`
	OrderId           uint64                   `json:"order_id"`
	MidtransOrderId   string                   `json:"midtrans_order_id"`
	GrossAmount       float32                  `json:"gross_amount"`
	Currency          string                   `json:"currency"`
	PaymentType       string                   `json:"payment_type"`
	ExpiryTime        time.Time                `json:"expiry_time"`
	TransactionTime   time.Time                `json:"transaction_time"`
	TransactionStatus helper.TransactionStatus `json:"transaction_status"`
	FraudStatus       string                   `json:"fraud_status"`
	Actions           *[]ActionResponse        `json:"actions"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}

type ActionResponse struct {
	Name   string               `json:"name"`
	Method helper.RequestMethod `json:"method"`
	URL    string               `json:"url"`
}

type CreateMidtransCoreAPIOrderRequest struct {
	OrderId uint64 `json:"order_id" validate:"required"`
}

type GetMidtransCoreAPIOrderRequest struct {
	OrderId uint64 `json:"order_id" validate:"required"`
}

type GetMidtransNotification struct {
	OrderId         string `json:"order_id"`
	StatusCode      string `json:"status_code"`
	GrossAmount     string `json:"gross_amount"`
	TransactionTime string `json:"transaction_time"`
	SignatureKey    string `json:"signature_key"`
	// Tambahkan field lain sesuai kebutuhan
}
