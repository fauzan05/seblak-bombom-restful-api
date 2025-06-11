package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/enum_state"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAdmin(t *testing.T) {
	ClearAll()
	DoCreateApplicationSettingByAdminToken(t)
	requestBody := model.RegisterUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "F3196813@gmail.com",
		Phone:     "08123456789",
		Password:  "JohnDoe123#",
		Role:      enum_state.ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Admin-Key", "rahasia-123#")

	response, err := app.Test(request, int(time.Second)*5)
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

func TestRegisterCustomer(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, tokenAdmin)
	requestBody := model.RegisterUserRequest{
		FirstName: "Customer",
		LastName:  "1",
		Email:     "fauzan.hidayat@binus.ac.id",
		Phone:     "0982131244",
		Password:  "Customer1#",
		Role:      enum_state.CUSTOMER,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users?lang=id", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
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
		Role:      enum_state.ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, int(time.Second)*5)
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
	DoRegisterAdmin(t)

	requestBody := model.RegisterUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		Phone:     "08123456789",
		Password:  "johndoe123",
		Role:      enum_state.ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, int(time.Second)*5)
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
	DoRegisterAdmin(t)

	requestBody := model.LoginUserRequest{
		Email:    "F3196813@gmail.com",
		Password: "JohnDoe123#",
	}

	DoVerificationEmail(t, requestBody.Email)

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, int(time.Second)*5)
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
	DoRegisterAdmin(t)

	requestBody := model.LoginUserRequest{
		Email:    "johndoe123@email.com",
		Password: "joe123",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "user not found : record not found", responseBody.Error)
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

	response, err := app.Test(request, int(time.Second)*5)
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

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "token isn't valid : token is expired!", responseBody.Error)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	request := httptest.NewRequest(http.MethodGet, "/api/users/current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second)*5)
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
	assert.Equal(t, enum_state.ADMIN, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserFailed(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	request := httptest.NewRequest(http.MethodGet, "/api/users/current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token+"adasd")

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "token isn't valid : token is expired!", responseBody.Error)
}

func TestChangePasswordFailedConfirmation(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.UpdateUserPasswordRequest{
		OldPassword:        "johndoe123",
		NewPassword:        "lala123",
		NewPasswordConfirm: "lala456",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/current/password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.False(t, responseBody.Data)
}

func TestChangePassword(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.UpdateUserPasswordRequest{
		OldPassword:        "johndoe123",
		NewPassword:        "testing123",
		NewPasswordConfirm: "testing123",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/current/password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestChangePasswordOldPasswordIsWrong(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.UpdateUserPasswordRequest{
		OldPassword:        "lala123",
		NewPassword:        "testing123",
		NewPasswordConfirm: "testing123",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/current/password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.False(t, responseBody.Data)
}

func TestChangeProfile(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.UpdateUserRequest{
		FirstName: "john-test",
		LastName:  "doe-test",
		Email:     "johndoe-test@mail.com",
		Phone:     "99999999999",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestChangeProfileEmailDuplicate(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	DoRegisterCustomer(t)
	token := DoLoginCustomer(t)

	requestBody := model.UpdateUserRequest{
		FirstName: "john-test",
		LastName:  "doe-test",
		Email:     "johndoe@email.com",
		Phone:     "99999999999",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, "/api/users/current", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusConflict, response.StatusCode)
	assert.Equal(t, "email has already exists!", responseBody.Error)
}

func TestForgotPasswordEmailNotFound(t *testing.T) {
	ClearAll()

	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "failed to find email address: record not found", responseBody.Error)
}

func TestForgotPasswordEmailFound(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)
}

func TestForgotPasswordResend(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	request = httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)
}

func TestForgotPasswordValidate(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	// validate
	requestBodyValidate := model.ValidateForgotPassword{
		VerificationCode: responseBody.Data.VerificationCode,
	}

	bodyJson, err = json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/validate", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBodyValidate.Data)
}

func TestForgotPasswordValidateExpired(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	// update tanggalnya ke 5 menit sebelum
	newPasswordReset := new(entity.PasswordReset)
	newPasswordReset.ID = responseBody.Data.ID
	db.Model(newPasswordReset).First(newPasswordReset)

	newPasswordReset.ExpiresAt = time.Now().Add((time.Minute * -5) + (time.Second * -1))
	db.Save(newPasswordReset)

	// validate
	requestBodyValidate := model.ValidateForgotPassword{
		VerificationCode: responseBody.Data.VerificationCode,
	}

	bodyJson, err = json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/validate", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "password reset was expired!", responseBodyValidate.Error)
}

func TestForgotPasswordValidateNotFound(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)

	// validate
	requestBodyValidate := model.ValidateForgotPassword{
		VerificationCode: 457236,
	}

	bodyJson, err := json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/validate", 1), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "password reset not found!", responseBodyValidate.Error)
}

func TestForgotPasswordValidateVerificationCodeNotMatch(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	// validate
	requestBodyValidate := model.ValidateForgotPassword{
		VerificationCode: 111111,
	}

	bodyJson, err = json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/validate", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "verification code is not match!", responseBodyValidate.Error)
}

func TestForgotPasswordResetPasswordSuccess(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	// validate
	requestBodyValidate := model.PasswordResetRequest{
		VerificationCode:   responseBody.Data.VerificationCode,
		NewPassword:        "Rahasia123#!",
		NewPasswordConfirm: "Rahasia123#!",
	}

	bodyJson, err = json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/reset-password", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.True(t, responseBodyValidate.Data)

	// coba login apakah bisa
	requestLogin := new(model.LoginUserRequest)
	requestLogin.Email = "F3196813@gmail.com"
	requestLogin.Password = "Rahasia123#!"

	bodyJson, err = json.Marshal(requestLogin)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyLogin := new(model.ApiResponse[model.UserTokenResponse])
	err = json.Unmarshal(bytes, responseBodyLogin)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBodyLogin.Data.Token)
	assert.True(t, responseBodyLogin.Data.ExpiryDate.After(time.Now()))
}

func TestForgotPasswordResetPasswordExpired(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	// ubah expired-nya ke -15 menit
	newPasswordReset := new(entity.PasswordReset)
	newPasswordReset.ID = responseBody.Data.ID
	db.Model(newPasswordReset).Update("expires_at", responseBody.Data.ExpiresAt.ToTime().Add((-20*time.Minute)+(-1*time.Second)))
	// validate
	requestBodyValidate := model.PasswordResetRequest{
		VerificationCode:   responseBody.Data.VerificationCode,
		NewPassword:        "Rahasia123#!",
		NewPasswordConfirm: "Rahasia123#!",
	}

	bodyJson, err = json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/reset-password", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "password reset was expired!", responseBodyValidate.Error)
}

func TestForgotPasswordResetPasswordPasswordConfirmNotSame(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	// validate
	requestBodyValidate := model.PasswordResetRequest{
		VerificationCode:   responseBody.Data.VerificationCode,
		NewPassword:        "Rahasia123#!",
		NewPasswordConfirm: "asdasdasd#!",
	}

	bodyJson, err = json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/reset-password", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "invalid request body : Key: 'PasswordResetRequest.NewPasswordConfirm' Error:Field validation for 'NewPasswordConfirm' failed on the 'eqfield' tag", responseBodyValidate.Error)
}

func TestForgotPasswordResetPasswordPasswordNotValid(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, token)

	DoRegisterCustomer(t)
	requestBody := model.CreateForgotPassword{
		Email: "F3196813@gmail.com",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/forgot-password", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.PasswordResetResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.VerificationCode)

	// validate
	requestBodyValidate := model.PasswordResetRequest{
		VerificationCode:   responseBody.Data.VerificationCode,
		NewPassword:        "rahasia123",
		NewPasswordConfirm: "rahasia123",
	}

	bodyJson, err = json.Marshal(requestBodyValidate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/users/forgot-password/%d/reset-password", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Host = "localhost"

	response, err = app.Test(request, int(time.Second)*5)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyValidate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBodyValidate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "Password must contain at least one uppercase letter;Password must contain at least one symbol (!@#~$%^&*()+|_);", responseBodyValidate.Error)
}
