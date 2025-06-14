package tests

import (
	"encoding/json"
	"fmt"
	// "fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"

	// "time"

	"github.com/stretchr/testify/assert"
)

func TestWithdrawRequestByCustomerWallet(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)

	// set saldo wallet
	DoSetBalanceManually(tokenCust, float32(100000))

	customer := GetCurrentUserByToken(t, tokenCust)

	requestBody := model.WithdrawWalletRequest{
		UserId:            customer.ID,
		Method:            enum_state.WALLET_WITHDRAW_REQUEST_METHOD_CASH,
		BankName:          "",
		BankAccountNumber: "",
		BankAccountName:   "",
		Status:            enum_state.WALLET_WITHDRAW_REQUEST_STATUS_PENDING,
		Amount:            93000,
		Note:              "Saya mau narik duit ya",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/wallets/withdraw-cust", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.WithdrawWalletResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, float32(93000), responseBody.Data.Amount)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

// Approved
func TestWithdrawRequestByCustomerAdminApproval(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)

	// set saldo wallet
	DoSetBalanceManually(tokenCust, float32(100000))

	customer := GetCurrentUserByToken(t, tokenCust)

	requestBody := model.WithdrawWalletRequest{
		UserId:            customer.ID,
		Method:            enum_state.WALLET_WITHDRAW_REQUEST_METHOD_CASH,
		BankName:          "",
		BankAccountNumber: "",
		BankAccountName:   "",
		Status:            enum_state.WALLET_WITHDRAW_REQUEST_STATUS_PENDING,
		Amount:            93000,
		Note:              "Saya mau narik duit ya",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/wallets/withdraw-cust", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.WithdrawWalletResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, float32(93000), responseBody.Data.Amount)
	assert.Equal(t, enum_state.WALLET_WITHDRAW_REQUEST_STATUS_PENDING, responseBody.Data.Status)
	assert.Equal(t, requestBody.Note, responseBody.Data.Note)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	tokenAdmin := DoLoginAdmin(t)

	requestBodyApproval := model.WithdrawWalletApprovalRequest{
		Status:         enum_state.WALLET_WITHDRAW_REQUEST_STATUS_APPROVED,
		RejectionNotes: "",
	}

	bodyJson, err = json.Marshal(requestBodyApproval)
	assert.Nil(t, err)
	requestApproval := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/wallets/%d/withdraw-approval", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	requestApproval.Header.Set("Content-Type", "application/json")
	requestApproval.Header.Set("Accept", "application/json")
	requestApproval.Header.Set("Authorization", tokenAdmin)

	responseApproval, err := app.Test(requestApproval)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(responseApproval.Body)
	assert.Nil(t, err)

	responseBodyApproval := new(model.ApiResponse[model.WithdrawWalletResponse])
	err = json.Unmarshal(bytes, responseBodyApproval)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBodyApproval.Data.ID)
	assert.Equal(t, float32(93000), responseBodyApproval.Data.Amount)
	assert.Equal(t, requestBodyApproval.RejectionNotes, responseBodyApproval.Data.RejectionNotes)
	assert.Equal(t, requestBodyApproval.Status, responseBodyApproval.Data.Status)
	assert.NotNil(t, responseBodyApproval.Data.CreatedAt)
	assert.NotNil(t, responseBodyApproval.Data.UpdatedAt)

	customerBalance := GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, float32(7000), customerBalance.Wallet.Balance)
}

// Rejected
func TestWithdrawRequestByCustomerAdminRejected(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)

	// set saldo wallet
	DoSetBalanceManually(tokenCust, float32(100000))

	customer := GetCurrentUserByToken(t, tokenCust)

	requestBody := model.WithdrawWalletRequest{
		UserId:            customer.ID,
		Method:            enum_state.WALLET_WITHDRAW_REQUEST_METHOD_CASH,
		BankName:          "",
		BankAccountNumber: "",
		BankAccountName:   "",
		Status:            enum_state.WALLET_WITHDRAW_REQUEST_STATUS_PENDING,
		Amount:            93000,
		Note:              "Saya mau narik duit ya",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/wallets/withdraw-cust", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.WithdrawWalletResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, float32(93000), responseBody.Data.Amount)
	assert.Equal(t, enum_state.WALLET_WITHDRAW_REQUEST_STATUS_PENDING, responseBody.Data.Status)
	assert.Equal(t, requestBody.Note, responseBody.Data.Note)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	tokenAdmin := DoLoginAdmin(t)

	requestBodyApproval := model.WithdrawWalletApprovalRequest{
		Status:         enum_state.WALLET_WITHDRAW_REQUEST_STATUS_REJECTED,
		RejectionNotes: "Gapunya duit saya",
	}

	bodyJson, err = json.Marshal(requestBodyApproval)
	assert.Nil(t, err)
	requestApproval := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/wallets/%d/withdraw-approval", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	requestApproval.Header.Set("Content-Type", "application/json")
	requestApproval.Header.Set("Accept", "application/json")
	requestApproval.Header.Set("Authorization", tokenAdmin)

	responseApproval, err := app.Test(requestApproval)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(responseApproval.Body)
	assert.Nil(t, err)

	responseBodyApproval := new(model.ApiResponse[model.WithdrawWalletResponse])
	err = json.Unmarshal(bytes, responseBodyApproval)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, responseApproval.StatusCode)
	assert.NotNil(t, responseBodyApproval.Data.ID)
	assert.Equal(t, float32(93000), responseBodyApproval.Data.Amount)
	assert.Equal(t, requestBodyApproval.RejectionNotes, responseBodyApproval.Data.RejectionNotes)
	assert.Equal(t, requestBodyApproval.Status, responseBodyApproval.Data.Status)
	assert.NotNil(t, responseBodyApproval.Data.CreatedAt)
	assert.NotNil(t, responseBodyApproval.Data.UpdatedAt)

	customerBalance := GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, float32(100000), customerBalance.Wallet.Balance)
}
