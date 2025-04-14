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
			IsMain: true,
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
