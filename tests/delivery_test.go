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

func TestUpdateDeliveriesBadRequest(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	deliveryResponse := DoCreateDelivery(t, token)

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

func TestUpdateDeliveries(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	deliveryResponse := DoCreateDelivery(t, token)

	requestBody := model.UpdateDeliveryRequest{
		City:     "Kebumen-test",
		District: "Pejagoan-test",
		Village:  "Peniron-test",
		Hamlet:   "Jetis-test",
		Cost:     10000,
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, responseBody.Data.City, requestBody.City)
	assert.Equal(t, responseBody.Data.District, requestBody.District)
	assert.Equal(t, responseBody.Data.Village, requestBody.Village)
	assert.Equal(t, responseBody.Data.Hamlet, requestBody.Hamlet)
	assert.Equal(t, responseBody.Data.Cost, requestBody.Cost)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateDeliveriesNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.UpdateDeliveryRequest{
		City:     "Kebumen-test",
		District: "Pejagoan-test",
		Village:  "Peniron-test",
		Hamlet:   "Jetis-test",
		Cost:     10000,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/deliveries/%+v", -99), strings.NewReader(string(bodyJson)))
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

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestGetAllDeliveryPagination(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	DoCreateManyDelivery(t, 27)

	request := httptest.NewRequest(http.MethodGet, "/api/deliveries?search=kebumen&column=deliveries.id&sort_by=asc&per_page=5&page=3", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.DeliveryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(*responseBody.Data))
	assert.Equal(t, 3, responseBody.CurrentPages)
	assert.Equal(t, int64(27), responseBody.TotalDatas)
	assert.Equal(t, 6, responseBody.TotalPages)
}

func TestGetAllDeliveryPaginationSortingColumnDesc(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	DoCreateManyDelivery(t, 27)

	request := httptest.NewRequest(http.MethodGet, "/api/deliveries?search=kebumen&column=deliveries.id&sort_by=desc&per_page=5&page=3", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.DeliveryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(*responseBody.Data))
	assert.Equal(t, 3, responseBody.CurrentPages)
	assert.Equal(t, int64(27), responseBody.TotalDatas)
	assert.Equal(t, 6, responseBody.TotalPages)

	deliveries := *responseBody.Data
	for i := range len(deliveries) - 1 {
		assert.Greater(t, deliveries[i].ID, deliveries[i+1].ID)
	}
}

func TestGetAllDeliveryPaginationSortingColumnAsc(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	DoCreateManyDelivery(t, 27)

	request := httptest.NewRequest(http.MethodGet, "/api/deliveries?search=kebumen&column=deliveries.id&sort_by=asc&per_page=5&page=3", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.DeliveryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(*responseBody.Data))
	assert.Equal(t, 3, responseBody.CurrentPages)
	assert.Equal(t, int64(27), responseBody.TotalDatas)
	assert.Equal(t, 6, responseBody.TotalPages)

	deliveries := *responseBody.Data
	for i := range len(deliveries) - 1 {
		assert.Less(t, deliveries[i].ID, deliveries[i+1].ID)
	}
}

func TestGetAllDeliveryPaginationSortingColumnNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	DoCreateManyDelivery(t, 27)

	request := httptest.NewRequest(http.MethodGet, "/api/deliveries?search=kebumen&column=deliveries.mama&sort_by=desc&per_page=5&page=3", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "invalid sort column : deliveries.mama", responseBody.Error)
}

func TestDeleteDeliveries(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	requestBody := DoCreateManyDelivery(t, 5)
	token := DoLoginAdmin(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/deliveries?ids="+requestBody, nil)
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
}

func TestDeleteDeliveriesFailed(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	requestBody := "e,b,s,s"
	token := DoLoginAdmin(t)

	request := httptest.NewRequest(http.MethodDelete, "/api/deliveries?ids="+requestBody, nil)
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

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}
