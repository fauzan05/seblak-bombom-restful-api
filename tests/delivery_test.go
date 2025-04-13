package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDelivery(t *testing.T) {
	// ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.CreateDeliveryRequest{
		City:     "Kebumen",
		District: "Pejagoan",
		Village:  "Peniron",
		Hamlet:   "Jetis",
		Cost:     5000,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/deliveries", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.DeliveryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.City, responseBody.Data.City)
	assert.Equal(t, requestBody.District, responseBody.Data.District)
	assert.Equal(t, requestBody.Village, responseBody.Data.Village)
	assert.Equal(t, requestBody.Hamlet, responseBody.Data.Hamlet)
	assert.Equal(t, requestBody.Cost, responseBody.Data.Cost)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestCreateDeliveryFailed(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.CreateDeliveryRequest{
		City:     "",
		District: "",
		Village:  "",
		Hamlet:   "",
		Cost:     0,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/deliveries", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.DeliveryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestUpdateDeliveries(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	deliveryResponse := DoCreateDelivery(t)
	token := DoLoginAdmin(t)

	requestBody := model.UpdateDeliveryRequest{
		City:     "",
		District: "",
		Village:  "",
		Hamlet:   "",
		Cost:     0,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/deliveries/%+v", deliveryResponse.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.DeliveryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}