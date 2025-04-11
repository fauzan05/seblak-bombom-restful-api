package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		Phone:     "08123456789",
		Password:  "johndoe123",
		Role:      helper.ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, requestBody.LastName, responseBody.Data.LastName)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.Equal(t, requestBody.Role, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		FirstName: "",
		LastName:  "",
		Email:     "",
		Phone:     "",
		Password:  "",
		Role:      helper.ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestRegisterEmailDuplicate(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.RegisterUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		Phone:     "08123456789",
		Password:  "johndoe123",
		Role:      helper.ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
}

func TestLogin(t *testing.T) {
	ClearAll()
	TestRegister(t)

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
	assert.NotNil(t, responseBody.Data.ExpiryDate)
}

func TestLoginFailed(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		Email:    "johndoe123@email.com",
		Password: "joe123",
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

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestLogout(t *testing.T) {
	ClearAll()
	TestLogin(t)

	user := new(entity.User)
	err := db.Preload("Token").Where("email = ?", "johndoe@email.com").First(&user).Error
	assert.Nil(t, err)
	assert.NotNil(t, user.Token.Token)

	request := httptest.NewRequest(http.MethodDelete, "/api/users/logout", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", user.Token.Token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestLogoutWrongAuthorization(t *testing.T) {
	ClearAll()

	fakeToken := "adasd2123asdasd"
	request := httptest.NewRequest(http.MethodDelete, "/api/users/logout", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", fakeToken)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	TestRegister(t)
	token := DoLogin(t)

	request := httptest.NewRequest(http.MethodGet, "/api/users/current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "John", responseBody.Data.FirstName)
	assert.Equal(t, "Doe", responseBody.Data.LastName)
	assert.Equal(t, "johndoe@email.com", responseBody.Data.Email)
	assert.Equal(t, "08123456789", responseBody.Data.Phone)
	assert.Equal(t, helper.ADMIN, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserFailed(t *testing.T) {
	ClearAll()
	TestRegister(t)
	token := DoLogin(t)

	request := httptest.NewRequest(http.MethodGet, "/api/users/current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token + "adasd")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestChangePassword(t *testing.T) {
	ClearAll()
	TestRegister(t)
	token := DoLogin(t)

	requestBody := model.UpdateUserPasswordRequest{
		OldPassword: "johndoe123",
		NewPassword: "lala123",
		NewPasswordConfirm: "lala456",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/current/password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}