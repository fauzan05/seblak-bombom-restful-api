package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddProductToCart(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	product := DoCreateProduct(t, tokenAdmin, 1, 1)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	requestBody := model.CreateCartRequest{
		ProductID: product.ID,
		Quantity:  3,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CartResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	currentUser := GetCurrentUserByToken(t, tokenCust)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, currentUser.ID, responseBody.Data.UserID)
	assert.Equal(t, 1, len(responseBody.Data.CartItems))
	for _, cartItem := range responseBody.Data.CartItems {
		assert.Equal(t, 3, cartItem.Quantity)
	}
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestAddProductProductNotFound(t *testing.T) {
	ClearAll()

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	requestBody := model.CreateCartRequest{
		ProductID: 1,
		Quantity:  3,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "failed to find product by id into product table : record not found", responseBody.Error)
}

func TestAddProductQuantityExceeded(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	product := DoCreateProduct(t, tokenAdmin, 1, 1)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	requestBody := model.CreateCartRequest{
		ProductID: product.ID,
		Quantity:  99999,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "quantity request exceeds available stock for product: Requested (99999), Available (-98998)", responseBody.Error)
}

func TestAddProductQuantityLowerThan1(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	product := DoCreateProduct(t, tokenAdmin, 1, 1)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	requestBody := model.CreateCartRequest{
		ProductID: product.ID,
		Quantity:  -1,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "quantity must be more than 0 if first time adding product to cart!", responseBody.Error)
}

func TestAddProductStockLowerThanQuantityRequest(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	product := DoCreateProduct(t, tokenAdmin, 1, 1)

	// update quantity
	newProduct := new(entity.Product)
	err := db.First(newProduct).Error
	assert.Nil(t, err)

	newProduct.Stock = 3
	db.Save(newProduct)

	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	requestBody := model.CreateCartRequest{
		ProductID: product.ID,
		Quantity:  10,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "quantity request exceeds available stock for product: Requested (10), Available (3)", responseBody.Error)
}

func TestAllCartItemByCurrentUser(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	product1 := DoCreateProduct(t, tokenAdmin, 1, 1)
	product2 := DoCreateProduct(t, tokenAdmin, 1, 1)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)

	for i := 1; i <= 5; i++ {
		requestBody := model.CreateCartRequest{
			ProductID: product1.ID,
			Quantity:  1,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", tokenCust)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.CartResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		currentUser := GetCurrentUserByToken(t, tokenCust)
		assert.NotNil(t, responseBody.Data.ID)
		assert.Equal(t, currentUser.ID, responseBody.Data.UserID)
		assert.Equal(t, 1, len(responseBody.Data.CartItems))
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)
	}

	for i := 1; i <= 5; i++ {
		requestBody := model.CreateCartRequest{
			ProductID: product2.ID,
			Quantity:  1,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", tokenCust)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.CartResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		currentUser := GetCurrentUserByToken(t, tokenCust)
		assert.NotNil(t, responseBody.Data.ID)
		assert.Equal(t, currentUser.ID, responseBody.Data.UserID)
		assert.Equal(t, 2, len(responseBody.Data.CartItems))
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/carts", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CartResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 2, len(responseBody.Data.CartItems))
	for _, cartItem := range responseBody.Data.CartItems {
		assert.Equal(t, 5, cartItem.Quantity)
	}
}

func TestUpdateCartItemByCurrentUser(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	product1 := DoCreateProduct(t, tokenAdmin, 1, 1)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)

	getCartItemId := 0
	for i := 1; i <= 5; i++ {
		requestBody := model.CreateCartRequest{
			ProductID: product1.ID,
			Quantity:  1,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", tokenCust)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.CartResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		currentUser := GetCurrentUserByToken(t, tokenCust)
		assert.NotNil(t, responseBody.Data.ID)
		assert.Equal(t, currentUser.ID, responseBody.Data.UserID)
		assert.Equal(t, 1, len(responseBody.Data.CartItems))
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		for _, cartItem := range responseBody.Data.CartItems {
			getCartItemId = int(cartItem.ID)
		}
	}

	newUpdateQuantity := new(model.UpdateCartRequest)
	newUpdateQuantity.Quantity = -2
	bodyJson, err := json.Marshal(newUpdateQuantity)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/carts/cart-items/%d", getCartItemId), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CartResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 1, len(responseBody.Data.CartItems))
	for _, cartItem := range responseBody.Data.CartItems {
		assert.Equal(t, 3, cartItem.Quantity)
	}
}

func TestDropProductFromCart(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	product1 := DoCreateProduct(t, tokenAdmin, 1, 1)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)

	requestBody := model.CreateCartRequest{
		ProductID: product1.ID,
		Quantity:  5,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/carts", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CartResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	currentUser := GetCurrentUserByToken(t, tokenCust)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, currentUser.ID, responseBody.Data.UserID)
	assert.Equal(t, 1, len(responseBody.Data.CartItems))
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	getCartId := uint64(0)
	for _, cartItem := range responseBody.Data.CartItems {
		getCartId = cartItem.ID
	}

	// DELETE
	request = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/carts/cart-items/%d", getCartId), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.ApiResponse[model.CartResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 0, len(responseBody.Data.CartItems))
	for _, cartItem := range responseBody.Data.CartItems {
		assert.Equal(t, 4, cartItem.Quantity)
	}
}
