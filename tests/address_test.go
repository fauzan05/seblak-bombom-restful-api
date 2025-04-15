package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	createDeliveryResponse := DoCreateDelivery(t, token)
	for i := 1; i <= 3; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      createDeliveryResponse.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6",
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.AddressResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.DeliveryId, responseBody.Data.Delivery.ID)
		assert.Equal(t, requestBody.CompleteAddress, responseBody.Data.CompleteAddress)
		assert.Equal(t, requestBody.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
		assert.Equal(t, requestBody.IsMain, responseBody.Data.IsMain)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)
	}
}

func TestCreateAddressFailed(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.AddressCreateRequest{
		CompleteAddress: "Complete Address",
		GoogleMapsLink:  "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6",
		IsMain:          true,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestUpdateAddress(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	var firstAddress model.AddressResponse
	createDeliveryResponse := DoCreateDelivery(t, token)
	for i := 1; i <= 3; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      createDeliveryResponse.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6",
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.AddressResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.DeliveryId, responseBody.Data.Delivery.ID)
		assert.Equal(t, requestBody.CompleteAddress, responseBody.Data.CompleteAddress)
		assert.Equal(t, requestBody.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
		assert.Equal(t, requestBody.IsMain, responseBody.Data.IsMain)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		// ambil data pertama karena nantinya tidak akan menjadi is_main = true
		if i == 1 {
			firstAddress = responseBody.Data
		}
	}

	updateAddress := new(model.UpdateAddressRequest)
	updateAddress.CompleteAddress = "Complete Address Update"
	updateAddress.DeliveryId = createDeliveryResponse.ID
	updateAddress.GoogleMapsLink = "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 Update"
	updateAddress.IsMain = true

	bodyJson, err := json.Marshal(updateAddress)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/current/addresses/%+v", firstAddress.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, firstAddress.ID, responseBody.Data.ID)
	assert.Equal(t, updateAddress.CompleteAddress, responseBody.Data.CompleteAddress)
	assert.Equal(t, updateAddress.DeliveryId, responseBody.Data.Delivery.ID)
	assert.Equal(t, updateAddress.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
	assert.Equal(t, updateAddress.IsMain, responseBody.Data.IsMain)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateAddressFailed(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	var firstAddress model.AddressResponse
	createDeliveryResponse := DoCreateDelivery(t, token)
	for i := 1; i <= 3; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      createDeliveryResponse.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6",
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.AddressResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.DeliveryId, responseBody.Data.Delivery.ID)
		assert.Equal(t, requestBody.CompleteAddress, responseBody.Data.CompleteAddress)
		assert.Equal(t, requestBody.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
		assert.Equal(t, requestBody.IsMain, responseBody.Data.IsMain)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		// ambil data pertama karena nantinya tidak akan menjadi is_main = true
		if i == 1 {
			firstAddress = responseBody.Data
		}
	}

	updateAddress := new(model.UpdateAddressRequest)
	updateAddress.CompleteAddress = ""
	updateAddress.DeliveryId = 0
	updateAddress.GoogleMapsLink = ""
	updateAddress.IsMain = true

	bodyJson, err := json.Marshal(updateAddress)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/users/current/addresses/%+v", firstAddress.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestGetAllAddressByCurrentUser(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	createDeliveryResponse := DoCreateDelivery(t, token)
	for i := 1; i <= 3; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      createDeliveryResponse.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  fmt.Sprintf("https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 %+v", i),
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.AddressResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.DeliveryId, responseBody.Data.Delivery.ID)
		assert.Equal(t, requestBody.CompleteAddress, responseBody.Data.CompleteAddress)
		assert.Equal(t, requestBody.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
		assert.Equal(t, requestBody.IsMain, responseBody.Data.IsMain)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/users/current/addresses", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*[]model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)
	lengthData := len(*responseBody.Data)
	for i, data := range *responseBody.Data {
		i = i + 1
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.NotNil(t, data.Delivery.ID)
		assert.Equal(t, fmt.Sprintf("Complete Address %+v", i), data.CompleteAddress)
		assert.Equal(t, fmt.Sprintf("https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 %+v", i), data.GoogleMapsLink)
		if i == lengthData {
			assert.Equal(t, true, data.IsMain)
		} else {
			assert.Equal(t, false, data.IsMain)
		}
	}
}

func TestGetAddressById(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	createDeliveryResponse := DoCreateDelivery(t, token)
	var getAddress model.AddressResponse
	for i := 1; i <= 3; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      createDeliveryResponse.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  fmt.Sprintf("https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 %+v", i),
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.AddressResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.DeliveryId, responseBody.Data.Delivery.ID)
		assert.Equal(t, requestBody.CompleteAddress, responseBody.Data.CompleteAddress)
		assert.Equal(t, requestBody.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
		assert.Equal(t, requestBody.IsMain, responseBody.Data.IsMain)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		if i == 1 {
			getAddress = responseBody.Data
		}
	}

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/users/current/addresses/%+v", getAddress.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, getAddress.Delivery.ID, responseBody.Data.Delivery.ID)
	assert.Equal(t, getAddress.CompleteAddress, responseBody.Data.CompleteAddress)
	assert.Equal(t, getAddress.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
	assert.Equal(t, false, responseBody.Data.IsMain)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestDeleteAddressByIds(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	createDeliveryResponse := DoCreateDelivery(t, token)
	var getIdAddress string
	var ids []uint64
	for i := 1; i <= 3; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      createDeliveryResponse.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  fmt.Sprintf("https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 %+v", i),
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.AddressResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.DeliveryId, responseBody.Data.Delivery.ID)
		assert.Equal(t, requestBody.CompleteAddress, responseBody.Data.CompleteAddress)
		assert.Equal(t, requestBody.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
		assert.Equal(t, requestBody.IsMain, responseBody.Data.IsMain)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		convertedToString := strconv.Itoa(int(responseBody.Data.ID))
		getIdAddress += convertedToString + ","
		ids = append(ids, responseBody.Data.ID)
	}

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/users/current/addresses?ids=%+v", getIdAddress), nil)
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
	assert.Equal(t, 3, len(ids))
	var result []entity.Address
	err = db.Unscoped().Where("id IN ?", ids).Find(&result).Error
	assert.Nil(t, err)

	for _, data := range result {
		assert.NotNil(t, data.DeletedAt)
	}
}

func TestFailedDeleteAddressByIds(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	createDeliveryResponse := DoCreateDelivery(t, token)
	getIdAddress := "e,1,3,b"
	for i := 1; i <= 3; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      createDeliveryResponse.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  fmt.Sprintf("https://maps.app.goo.gl/ftF7eEsBHa69uw3H6 %+v", i),
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.AddressResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.DeliveryId, responseBody.Data.Delivery.ID)
		assert.Equal(t, requestBody.CompleteAddress, responseBody.Data.CompleteAddress)
		assert.Equal(t, requestBody.GoogleMapsLink, responseBody.Data.GoogleMapsLink)
		assert.Equal(t, requestBody.IsMain, responseBody.Data.IsMain)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)
	}

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/users/current/addresses?ids=%+v", getIdAddress), nil)
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