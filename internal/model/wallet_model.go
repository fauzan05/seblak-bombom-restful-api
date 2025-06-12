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

type WithdrawWalletRequest struct {
	UserId            uint64                           `json:"user_id" validate:"required"`
	Amount            float32                          `json:"amount" validate:"required"`
	Method            enum_state.WalletWithdrawRequest `json:"method" validate:"required"`
	BankName          string                           `json:"bank_name"`
	BankAccountNumber string                           `json:"bank_account_number"`
	BankAccountName   string                           `json:"bank_account_name"`
	Status            enum_state.WalletWithdrawRequest `json:"status" validate:"required"`
	Note              string                           `json:"note"`
}

type WithdrawWalletApprovalRequest struct {
	ID             uint64                           `json:"-" validate:"required"`
	Status         enum_state.WalletWithdrawRequest `json:"status" validate:"required"`
	RejectionNotes string                           `json:"rejection_notes"`
	CurrentAdminId uint64                           `json:"-" validate:"required"`
}

type WithdrawWalletResponse struct {
	ID                uint64                           `json:"id"`
	UserId            uint64                           `json:"user_id"`
	User              UserResponse                     `json:"user"`
	Amount            float32                          `json:"amount"`
	Method            enum_state.WalletWithdrawRequest `json:"method"`
	BankName          string                           `json:"bank_name"`
	BankAccountNumber string                           `json:"bank_account_number"`
	BankAccountName   string                           `json:"bank_account_name"`
	Status            enum_state.WalletWithdrawRequest `json:"status"`
	Note              string                           `json:"note"`
	RejectionNotes    string                           `json:"rejection_notes"`
	ProcessedBy       UserResponse                     `json:"processed_by"`
	ProcessedAt       helper_others.TimeRFC3339        `json:"processed_at"`
	CreatedAt         helper_others.TimeRFC3339        `json:"created_at"`
	UpdatedAt         helper_others.TimeRFC3339        `json:"updated_at"`
}

type UpdateWalletBalance struct {
	ID      uint64  `json:"-" validate:"required"`
	Balance float32 `json:"balance" validate:"required"`
}

type SuspendWallet struct {
	IDs []uint64 `json:"-" validate:"required"`
}
