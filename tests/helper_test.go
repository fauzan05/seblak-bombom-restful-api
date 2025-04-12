package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ClearAll() {
	ClearTokens()
	ClearWallets()
	ClearDeliveries()
	ClearAddresses()
	ClearCarts()
	ClearUsers()
}

func ClearTokens() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Token{}).Error
	if err != nil {
		log.Fatalf("Failed clear token data : %+v", err)
	}
}
func ClearCarts() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Cart{}).Error
	if err != nil {
		log.Fatalf("Failed clear cart data : %+v", err)
	}
}

func ClearDeliveries() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Delivery{}).Error
	if err != nil {
		log.Fatalf("Failed clear delivery data : %+v", err)
	}
}

func ClearAddresses() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Address{}).Error
	if err != nil {
		log.Fatalf("Failed clear address data : %+v", err)
	}
}

func ClearWallets() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Wallet{}).Error
	if err != nil {
		log.Fatalf("Failed clear wallet data : %+v", err)
	}
}

func ClearUsers() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func DoLoginAdmin(t *testing.T) string {
	requestBody := model.LoginUserRequest{
		Email:    "johndoe@email.com",
		Password: "johndoe123",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserTokenResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)

	return responseBody.Data.Token
}

func DoLoginCustomer(t *testing.T) string {
	requestBody := model.LoginUserRequest{
		Email:    "customer1@email.com",
		Password: "customer1",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserTokenResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)

	return responseBody.Data.Token
}