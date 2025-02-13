package model

import "seblak-bombom-restful-api/internal/helper"

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
	ExpiryTime        string                   `json:"expiry_time"`
	TransactionTime   string                   `json:"transaction_time"`
	TransactionStatus helper.TransactionStatus `json:"transaction_status"`
	FraudStatus       string                   `json:"fraud_status"`
	Actions           *[]ActionResponse         `json:"actions"`
	CreatedAt         string                   `json:"created_at"`
	UpdatedAt         string                   `json:"updated_at"`
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
