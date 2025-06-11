package model

import (
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"time"
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

type WithdrawWalletRequest struct {
	UserId            uint64                           `json:"user_id" validate:"required"`
	Amount            float32                          `json:"amount" validate:"required"`
	Method            enum_state.WalletWithdrawRequest `json:"method" validate:"required"`
	BankName          string                           `json:"bank_name"`
	BankAccountNumber string                           `json:"bank_account_number"`
	BankAccountName   string                           `json:"bank_account_name"`
	Status            enum_state.WalletWithdrawRequest `json:"status" validate:"required"`
	RejectionNotes    string                           `json:"rejection_notes"`
	ProcessedBy       *uint64                          `json:"processed_by"`
	ProcessedAt       time.Time                        `json:"processed_at"`
	Notes             string                           `json:"notes"`
	AdminNotes        string                           `json:"admin_notes"`
}

type WithdrawWalletApprovalRequest struct {
	UserId            uint64                           `json:"user_id" validate:"required"`
	Amount            float32                          `json:"amount" validate:"required"`
	Method            enum_state.WalletWithdrawRequest `json:"method" validate:"required"`
	BankName          string                           `json:"bank_name"`
	BankAccountNumber string                           `json:"bank_account_number"`
	BankAccountName   string                           `json:"bank_account_name"`
	Status            enum_state.WalletWithdrawRequest `json:"status" validate:"required"`
	RejectionNotes    string                           `json:"rejection_notes"`
	ProcessedBy       *uint64                          `json:"processed_by"`
	ProcessedAt       time.Time                        `json:"processed_at"`
	Notes             string                           `json:"notes"`
	AdminNotes        string                           `json:"admin_notes"`
}

type UpdateWalletBalance struct {
	ID      uint64  `json:"-" validate:"required"`
	Balance float32 `json:"balance" validate:"required"`
}

type SuspendWallet struct {
	IDs []uint64 `json:"-" validate:"required"`
}
