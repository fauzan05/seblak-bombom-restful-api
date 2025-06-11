package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDrawWallet(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)

	// set saldo wallet
	DoSetBalanceManually(tokenCust, float32(100000))

	customer := GetCurrentUserByToken(t, tokenCust)

	requestBody := model.WithDrawWalletRequest{
		UserId: customer.ID,
		Amount: 93000,
		Notes: "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPut, "/api/wallets/withdraw", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenAdmin)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.WalletResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, float32(7000), responseBody.Data.Balance)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}
