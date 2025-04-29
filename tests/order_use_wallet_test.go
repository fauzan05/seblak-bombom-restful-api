package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrderAsAdminWithoutDeliveryAndDiscount(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	currentUser := GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000))

	DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  2,
			},
			{
				ProductId: product.ID,
				Quantity:  2,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.Invoice)
	assert.Equal(t, helper.PERCENT, responseBody.Data.DiscountType)
	assert.Equal(t, float32(0), responseBody.Data.DiscountValue)
	assert.Equal(t, float32(0), responseBody.Data.TotalDiscount)
	assert.Equal(t, currentUser.ID, responseBody.Data.UserId)
	assert.Equal(t, currentUser.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, currentUser.LastName, responseBody.Data.LastName)
	assert.Equal(t, currentUser.Email, responseBody.Data.Email)
	assert.Equal(t, currentUser.Phone, responseBody.Data.Phone)
	assert.Equal(t, helper.PAYMENT_GATEWAY_SYSTEM, responseBody.Data.PaymentGateway)
	assert.Equal(t, helper.PAYMENT_METHOD_EWALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.ORDER_PENDING, responseBody.Data.OrderStatus)
	assert.Equal(t, false, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(0), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)
	var totalProductPrice float32 = product.Price * 4

	assert.Equal(t, totalProductPrice, responseBody.Data.TotalProductPrice)
	assert.Equal(t, totalProductPrice+0-responseBody.Data.TotalDiscount, responseBody.Data.TotalFinalPrice)
	assert.Equal(t, len(requestBody.OrderProducts), len(responseBody.Data.OrderProducts))
	for i, product := range responseBody.Data.OrderProducts {
		assert.Equal(t, requestBody.OrderProducts[i].ProductId, product.ProductId)
		assert.Equal(t, requestBody.OrderProducts[i].Quantity, product.Quantity)
	}

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestCreateOrderAsAdminWithDeliveryAndNoDiscount(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	currentUser := GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000))

	getDelivery := DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     true,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  2,
			},
			{
				ProductId: product.ID,
				Quantity:  2,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.Invoice)
	assert.Equal(t, helper.PERCENT, responseBody.Data.DiscountType)
	assert.Equal(t, float32(0), responseBody.Data.DiscountValue)
	assert.Equal(t, float32(0), responseBody.Data.TotalDiscount)
	assert.Equal(t, currentUser.ID, responseBody.Data.UserId)
	assert.Equal(t, currentUser.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, currentUser.LastName, responseBody.Data.LastName)
	assert.Equal(t, currentUser.Email, responseBody.Data.Email)
	assert.Equal(t, currentUser.Phone, responseBody.Data.Phone)
	assert.Equal(t, helper.PAYMENT_GATEWAY_SYSTEM, responseBody.Data.PaymentGateway)
	assert.Equal(t, helper.PAYMENT_METHOD_EWALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.ORDER_PENDING, responseBody.Data.OrderStatus)
	assert.Equal(t, true, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(getDelivery.Delivery.Cost), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)
	var totalProductPrice float32 = product.Price * 4

	assert.Equal(t, totalProductPrice, responseBody.Data.TotalProductPrice)
	assert.Equal(t, totalProductPrice+getDelivery.Delivery.Cost-responseBody.Data.TotalDiscount, responseBody.Data.TotalFinalPrice)
	assert.Equal(t, len(requestBody.OrderProducts), len(responseBody.Data.OrderProducts))
	for i, product := range responseBody.Data.OrderProducts {
		assert.Equal(t, requestBody.OrderProducts[i].ProductId, product.ProductId)
		assert.Equal(t, requestBody.OrderProducts[i].Quantity, product.Quantity)
	}

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestCreateOrderAsAdminWithDeliveryAndDiscount(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-04-01T00:00:01+07:00"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-04-29T23:59:59+07:00"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	currentUser := GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000))

	getDelivery := DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     true,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  2,
			},
			{
				ProductId: product.ID,
				Quantity:  2,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.Invoice)
	assert.Equal(t, helper.PERCENT, responseBody.Data.DiscountType)
	assert.Equal(t, float32(5), responseBody.Data.DiscountValue)
	assert.Equal(t, float32(5250.2), responseBody.Data.TotalDiscount)
	assert.Equal(t, currentUser.ID, responseBody.Data.UserId)
	assert.Equal(t, currentUser.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, currentUser.LastName, responseBody.Data.LastName)
	assert.Equal(t, currentUser.Email, responseBody.Data.Email)
	assert.Equal(t, currentUser.Phone, responseBody.Data.Phone)
	assert.Equal(t, helper.PAYMENT_GATEWAY_SYSTEM, responseBody.Data.PaymentGateway)
	assert.Equal(t, helper.PAYMENT_METHOD_EWALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.ORDER_PENDING, responseBody.Data.OrderStatus)
	assert.Equal(t, true, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(getDelivery.Delivery.Cost), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)
	var totalProductPrice float32 = product.Price * 4

	assert.Equal(t, totalProductPrice, responseBody.Data.TotalProductPrice)
	assert.Equal(t, totalProductPrice+getDelivery.Delivery.Cost-responseBody.Data.TotalDiscount, responseBody.Data.TotalFinalPrice)
	assert.Equal(t, len(requestBody.OrderProducts), len(responseBody.Data.OrderProducts))
	for i, product := range responseBody.Data.OrderProducts {
		assert.Equal(t, requestBody.OrderProducts[i].ProductId, product.ProductId)
		assert.Equal(t, requestBody.OrderProducts[i].Quantity, product.Quantity)
	}

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestCreateOrderBalanceInsufficient(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoSetBalanceManually(token, float32(50000))

	DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  2,
			},
			{
				ProductId: product.ID,
				Quantity:  2,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "your balance is insufficient to perform this transaction!", responseBody.Error)
}

func TestCreateOrderDiscountNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoSetBalanceManually(token, float32(150000))

	DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     1,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  2,
			},
			{
				ProductId: product.ID,
				Quantity:  2,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "discount has disabled or doesn't exists!", responseBody.Error)
}

func TestCreateOrderDiscountNotYetActiveDate(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-05-01T00:00:01+07:00"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-05-15T23:59:59+07:00"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))

	DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  2,
			},
			{
				ProductId: product.ID,
				Quantity:  2,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "discount is not yet valid. It will be active starting May 25 2025 at 00:00:01", responseBody.Error)
}

func TestCreateOrderDiscountExpired(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-04-01T00:00:01+07:00"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-04-28T23:59:59+07:00"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))

	DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  2,
			},
			{
				ProductId: product.ID,
				Quantity:  2,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "discount has expired and is no longer available!", responseBody.Error)
}

func TestCreateOrderDiscountMinOrder(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-04-01T00:00:01+07:00"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-04-30T23:59:59+07:00"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))

	DoCreateManyAddress(t, token, 2, 1)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  1,
			},
		},
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/orders", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, "the order does not meet the minimum purchase requirements for this discount coupon!", responseBody.Error)
}
