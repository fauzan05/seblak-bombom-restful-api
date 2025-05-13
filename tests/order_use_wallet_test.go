package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrderAsAdminWithoutDeliveryAndDiscount(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoCreateApplicationSetting(t, tokenAdmin)
	DoRegisterCustomer(t)
	tokenCustomer := DoLoginCustomer(t)
	DoSetBalanceManually(tokenCustomer, float32(150000))
	currentUser := GetCurrentUserByToken(t, tokenCustomer)
	delivery := DoCreateDelivery(t, tokenAdmin)
	DoCreateManyAddress(t, tokenCustomer, 2, 1, delivery)
	product := DoCreateProduct(t, tokenAdmin, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
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
	request := httptest.NewRequest(http.MethodPost, "/api/orders?lang=id", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCustomer)
	request.Host = "localhost"

	response, err := app.Test(request, int(time.Second) * 5)
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
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

	// cek saldo
	currentUser = GetCurrentUserByToken(t, tokenCustomer)
	assert.Equal(t, (float32(150000) - responseBody.Data.TotalFinalPrice), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestCreateOrderAsAdminWithDiscountButDeliveryDeleted(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	currentUser := GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
	newDelivery := new(entity.Delivery)
	db.Model(entity.Delivery{}).First(newDelivery)
	db.Delete(newDelivery)

	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
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
	// cek saldo
	currentUser = GetCurrentUserByToken(t, token)
	assert.Equal(t, (float32(150000) - responseBody.Data.TotalFinalPrice), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestCreateOrderAsAdminWithDeliveryAndNoDiscount(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	currentUser := GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	getDelivery := DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
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

	// cek saldo
	currentUser = GetCurrentUserByToken(t, token)
	assert.Equal(t, (float32(150000) - responseBody.Data.TotalFinalPrice), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestCreateOrderAsAdminWithDeliveryAndDeliveryDeleted(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)

	newDelivery := new(entity.Delivery)
	db.Model(entity.Delivery{}).First(newDelivery)
	db.Delete(newDelivery)

	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
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

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "delivery not found, please selected one!", responseBody.Error)
}

func TestCreateOrderAsAdminWithDeliveryAndDiscount(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	currentUser := GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	getDelivery := DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
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

	// cek saldo
	currentUser = GetCurrentUserByToken(t, token)
	assert.Equal(t, helper.RoundFloat32((float32(150000)-responseBody.Data.TotalFinalPrice), 1), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestCreateOrderAsAdminWithDeliveryAndDiscountUsageExceededLimit(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 2, 50000, true)

	currentUser := GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(1500000))
	delivery := DoCreateDelivery(t, token)
	getDelivery := DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)

	for i := 1; i <= 3; i++ {
		requestBody := model.CreateOrderRequest{
			DiscountId:     getDiscountCoupon.ID,
			PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
			PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
			ChannelCode:    helper.WALLET_CHANNEL_CODE,
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

		if i == 3 {
			responseBody := new(model.ErrorResponse[string])
			err = json.Unmarshal(bytes, responseBody)
			assert.Nil(t, err)

			assert.Equal(t, "the usage limit for this discount coupon has been exceeded!", responseBody.Error)
		} else {
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
			assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
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
		}
	}
}

func TestCreateOrderBalanceInsufficient(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoSetBalanceManually(token, float32(50000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     0,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
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

	// cek saldo
	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(50000), currentUser.Wallet.Balance)
}

func TestCreateOrderDiscountNotFound(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
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

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderDiscountNotYetActiveDate(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(1, 0, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(5, 0, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
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
	assert.Equal(t, fmt.Sprintf("discount is not yet valid. It will be active starting %+s", parseStart.Format("January 02 2006 at 15:04:05")), responseBody.Error)

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderDiscountExpired(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, -1, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(-1, 0, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
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
	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderDiscountMinOrder(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(0, 1, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
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

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderPaymentMethodNotValid(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(0, 1, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentMethod:  "KAKA",
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  4,
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
	assert.Equal(t, "invalid payment method!", responseBody.Error)

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderChannelCodeNotValid(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(0, 0, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    "KAKALA",
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  4,
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
	assert.Equal(t, "invalid channel code!", responseBody.Error)

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderPaymentGatewayNotValid(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(0, 0, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.WALLET_CHANNEL_CODE,
		PaymentGateway: "LALA",
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  4,
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
	assert.Equal(t, "invalid payment gateway!", responseBody.Error)

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderWalletWrongPaymentMethod(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(0, 0, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_EWALLET,
		ChannelCode:    helper.XENDIT_EWALLET_DANA_CHANNEL_CODE,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  4,
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
	assert.Equal(t, "payment method EWALLET is not available on payment gateway system!", responseBody.Error)

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestCreateOrderWalletWrongChannelCode(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 0, 1)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(0, 0, 0, 23, 59, 0)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	DoSetBalanceManually(token, float32(150000))
	delivery := DoCreateDelivery(t, token)
	DoCreateManyAddress(t, token, 2, 1, delivery)
	product := DoCreateProduct(t, token, 2, 1)
	requestBody := model.CreateOrderRequest{
		DiscountId:     getDiscountCoupon.ID,
		PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
		PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
		ChannelCode:    helper.XENDIT_EWALLET_DANA_CHANNEL_CODE,
		IsDelivery:     false,
		Note:           "Yang cepet ya!",
		OrderProducts: []model.CreateOrderProductRequest{
			{
				ProductId: product.ID,
				Quantity:  4,
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
	assert.Equal(t, "channel code EWALLET_DANA is not available on payment gateway system!", responseBody.Error)

	currentUser := GetCurrentUserByToken(t, token)
	assert.Equal(t, float32(150000), currentUser.Wallet.Balance)
}

func TestGetAllOrderPagination(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	totalOrder := 20
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), totalOrder*10, totalOrder, 20000, true)
	product := DoCreateProduct(t, token, 2, 1)
	delivery := DoCreateDelivery(t, token)
	DoCreateManyOrderUsingWalletPayment(t, token, totalOrder, getDiscountCoupon, product, delivery)

	request := httptest.NewRequest(http.MethodGet, "/api/orders?per_page=5&page=2", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(20), responseBody.TotalDatas)
	assert.Equal(t, 4, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	for _, order := range *responseBody.Data {
		assert.NotNil(t, order.ChannelCode)
		assert.NotNil(t, order.CompleteAddress)
		assert.NotNil(t, order.ID)
		for _, orderProduct := range order.OrderProducts {
			assert.NotNil(t, orderProduct.ID)
			assert.NotNil(t, orderProduct.OrderId)
			assert.NotNil(t, orderProduct.ProductName)
			for _, images := range orderProduct.Product.Images {
				assert.NotNil(t, images.FileName)
			}
		}
	}
}

func TestGetAllOrderPaginationSortingColumnDesc(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	totalOrder := 20
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), totalOrder*10, totalOrder, 20000, true)
	product := DoCreateProduct(t, token, 2, 1)
	delivery := DoCreateDelivery(t, token)
	DoCreateManyOrderUsingWalletPayment(t, token, totalOrder, getDiscountCoupon, product, delivery)

	request := httptest.NewRequest(http.MethodGet, "/api/orders?per_page=5&page=2&column=orders.id&sort_by=desc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(20), responseBody.TotalDatas)
	assert.Equal(t, 4, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	for _, order := range *responseBody.Data {
		assert.NotNil(t, order.ChannelCode)
		assert.NotNil(t, order.CompleteAddress)
		assert.NotNil(t, order.ID)
		for _, orderProduct := range order.OrderProducts {
			assert.NotNil(t, orderProduct.ID)
			assert.NotNil(t, orderProduct.OrderId)
			assert.NotNil(t, orderProduct.ProductName)
			for _, images := range orderProduct.Product.Images {
				assert.NotNil(t, images.FileName)
			}
		}
	}

	orders := *responseBody.Data
	for i := range len(orders) - 1 {
		assert.Greater(t, orders[i].ID, orders[i+1].ID)
	}
}

func TestGetAllOrderPaginationColumnNotFound(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	totalOrder := 20
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), totalOrder*10, totalOrder, 20000, true)
	product := DoCreateProduct(t, token, 2, 1)
	delivery := DoCreateDelivery(t, token)
	DoCreateManyOrderUsingWalletPayment(t, token, totalOrder, getDiscountCoupon, product, delivery)

	request := httptest.NewRequest(http.MethodGet, "/api/orders?per_page=5&page=2&column=orders.mama", nil)
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
	assert.Equal(t, "invalid sort column : orders.mama", responseBody.Error)
}

func TestGetAllOrderPaginationSomeProductDeleted(t *testing.T) {
	// product harus tetap berelasi meskipun produk sudah dihapus
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	totalOrder := 20
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), totalOrder*10, totalOrder, 20000, true)
	product := DoCreateProduct(t, token, 2, 1)
	delivery := DoCreateDelivery(t, token)
	DoCreateManyOrderUsingWalletPayment(t, token, totalOrder, getDiscountCoupon, product, delivery)
	newProduct := new(entity.Product)
	db.Model(entity.Product{}).First(newProduct)

	db.Delete(newProduct)

	request := httptest.NewRequest(http.MethodGet, "/api/orders?per_page=5&page=2", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(20), responseBody.TotalDatas)
	assert.Equal(t, 4, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	for _, order := range *responseBody.Data {
		assert.NotNil(t, order.ChannelCode)
		assert.NotNil(t, order.CompleteAddress)
		assert.NotNil(t, order.ID)
		for _, orderProduct := range order.OrderProducts {
			assert.NotNil(t, orderProduct.ID)
			assert.NotNil(t, orderProduct.OrderId)
			assert.NotNil(t, orderProduct.ProductName)
			for _, images := range orderProduct.Product.Images {
				assert.NotNil(t, images.FileName)
			}
		}
	}
}

func TestGetAllOrderPaginationSearchProduct(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	totalOrder := 20
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), totalOrder*10, totalOrder, 20000, true)
	product := DoCreateProduct(t, token, 2, 1)
	delivery := DoCreateDelivery(t, token)
	DoCreateManyOrderUsingWalletPayment(t, token, totalOrder, getDiscountCoupon, product, delivery)

	request := httptest.NewRequest(http.MethodGet, "/api/orders?per_page=10&page=1&search=produk", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(20), responseBody.TotalDatas)
	assert.Equal(t, 2, responseBody.TotalPages)
	assert.Equal(t, 1, responseBody.CurrentPages)
	for _, order := range *responseBody.Data {
		assert.NotNil(t, order.ChannelCode)
		assert.NotNil(t, order.CompleteAddress)
		assert.NotNil(t, order.ID)
		for _, orderProduct := range order.OrderProducts {
			assert.NotNil(t, orderProduct.ID)
			assert.NotNil(t, orderProduct.OrderId)
			assert.NotNil(t, orderProduct.ProductName)
			for _, images := range orderProduct.Product.Images {
				assert.NotNil(t, images.FileName)
			}
		}
	}
}

func TestGetAllOrderPaginationSearchProductNotFound(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	totalOrder := 20
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), totalOrder*10, totalOrder, 20000, true)
	product := DoCreateProduct(t, token, 2, 1)
	delivery := DoCreateDelivery(t, token)
	DoCreateManyOrderUsingWalletPayment(t, token, totalOrder, getDiscountCoupon, product, delivery)

	request := httptest.NewRequest(http.MethodGet, "/api/orders?per_page=10&page=1&search=alala", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(0), responseBody.TotalDatas)
	assert.Equal(t, 0, responseBody.TotalPages)
	assert.Equal(t, 1, responseBody.CurrentPages)
	for _, order := range *responseBody.Data {
		assert.NotNil(t, order.ChannelCode)
		assert.NotNil(t, order.CompleteAddress)
		assert.NotNil(t, order.ID)
		for _, orderProduct := range order.OrderProducts {
			assert.NotNil(t, orderProduct.ID)
			assert.NotNil(t, orderProduct.OrderId)
			assert.NotNil(t, orderProduct.ProductName)
			for _, images := range orderProduct.Product.Images {
				assert.NotNil(t, images.FileName)
			}
		}
	}
}

func TestGetAllCurrentUserOrderPagination(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	totalOrder := 20
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, token, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(10), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), totalOrder, totalOrder, 20000, true)
	product := DoCreateProduct(t, token, 2, 1)

	delivery := DoCreateDelivery(t, token)
	DoCreateManyOrderUsingWalletPayment(t, token, totalOrder, getDiscountCoupon, product, delivery)

	DoRegisterCustomer(t)
	token = DoLoginCustomer(t)
	DoCreateManyOrderUsingWalletPayment(t, token, 5, getDiscountCoupon, product, delivery)

	request := httptest.NewRequest(http.MethodGet, "/api/orders?per_page=10&page=1", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, int64(5), responseBody.TotalDatas)
	assert.Equal(t, 1, responseBody.TotalPages)
	assert.Equal(t, 1, responseBody.CurrentPages)
	for _, order := range *responseBody.Data {
		assert.NotNil(t, order.ChannelCode)
		assert.NotNil(t, order.CompleteAddress)
		assert.NotNil(t, order.ID)
		for _, orderProduct := range order.OrderProducts {
			assert.NotNil(t, orderProduct.ID)
			assert.NotNil(t, orderProduct.OrderId)
			assert.NotNil(t, orderProduct.ProductName)
			for _, images := range orderProduct.Product.Images {
				assert.NotNil(t, images.FileName)
			}
		}
	}
}

func TestUpdateOrderReceivedAsCustomerError(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.ORDER_RECEIVED

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
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

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "admin access only!", responseBody.Error)
}

func TestUpdateReceivedOrderAsAdmin(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.ORDER_RECEIVED

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenAdmin)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	currentUser := GetCurrentUserByToken(t, tokenCust)
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.ORDER_RECEIVED, responseBody.Data.OrderStatus)
	assert.Equal(t, true, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(order.DeliveryCost), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)

	assert.Equal(t, float32(100004), responseBody.Data.TotalProductPrice)
	assert.Equal(t, float32(99753.8), responseBody.Data.TotalFinalPrice)
	assert.Equal(t, 2, len(responseBody.Data.OrderProducts))
	for _, product := range responseBody.Data.OrderProducts {
		assert.NotNil(t, product.ID)
		assert.Equal(t, 2, product.Quantity)
	}

	// cek saldo
	currentUser = GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, helper.RoundFloat32((float32(150000)-responseBody.Data.TotalFinalPrice), 1), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestUpdateReadyForPickupOrderAsAdmin(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.READY_FOR_PICKUP

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenAdmin)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	currentUser := GetCurrentUserByToken(t, tokenCust)
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.READY_FOR_PICKUP, responseBody.Data.OrderStatus)
	assert.Equal(t, true, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(order.DeliveryCost), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)

	assert.Equal(t, float32(100004), responseBody.Data.TotalProductPrice)
	assert.Equal(t, float32(99753.8), responseBody.Data.TotalFinalPrice)
	assert.Equal(t, 2, len(responseBody.Data.OrderProducts))
	for _, product := range responseBody.Data.OrderProducts {
		assert.NotNil(t, product.ID)
		assert.Equal(t, 2, product.Quantity)
	}

	// cek saldo
	currentUser = GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, helper.RoundFloat32((float32(150000)-responseBody.Data.TotalFinalPrice), 1), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestUpdateReadyForPickupOrderAsCustomer(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.READY_FOR_PICKUP

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
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

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "admin access only!", responseBody.Error)
}

func TestUpdateBeingDeliveredOrderAsAdmin(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.ORDER_BEING_DELIVERED

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenAdmin)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	currentUser := GetCurrentUserByToken(t, tokenCust)
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.ORDER_BEING_DELIVERED, responseBody.Data.OrderStatus)
	assert.Equal(t, true, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(order.DeliveryCost), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)

	assert.Equal(t, float32(100004), responseBody.Data.TotalProductPrice)
	assert.Equal(t, float32(99753.8), responseBody.Data.TotalFinalPrice)
	assert.Equal(t, 2, len(responseBody.Data.OrderProducts))
	for _, product := range responseBody.Data.OrderProducts {
		assert.NotNil(t, product.ID)
		assert.Equal(t, 2, product.Quantity)
	}

	// cek saldo
	currentUser = GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, helper.RoundFloat32((float32(150000)-responseBody.Data.TotalFinalPrice), 1), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestUpdateBeingDeliveredOrderAsCustomer(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.ORDER_BEING_DELIVERED

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
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

	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
	assert.Equal(t, "admin access only!", responseBody.Error)
}

func TestUpdateDeliveredOrderAsAdmin(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.ORDER_DELIVERED

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenAdmin)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	currentUser := GetCurrentUserByToken(t, tokenCust)
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.ORDER_DELIVERED, responseBody.Data.OrderStatus)
	assert.Equal(t, true, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(order.DeliveryCost), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)

	assert.Equal(t, float32(100004), responseBody.Data.TotalProductPrice)
	assert.Equal(t, float32(99753.8), responseBody.Data.TotalFinalPrice)
	assert.Equal(t, 2, len(responseBody.Data.OrderProducts))
	for _, product := range responseBody.Data.OrderProducts {
		assert.NotNil(t, product.ID)
		assert.Equal(t, 2, product.Quantity)
	}

	// cek saldo
	currentUser = GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, helper.RoundFloat32((float32(150000)-responseBody.Data.TotalFinalPrice), 1), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}

func TestUpdateDeliveredOrderAsCustomer(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	tokenAdmin := DoLoginAdmin(t)
	DoRegisterCustomer(t)
	tokenCust := DoLoginCustomer(t)
	order := DoCreateOrderAsCustomerWithDeliveryAndDiscount(t, tokenAdmin, tokenCust)

	requestBodyUpdate := new(model.UpdateOrderRequest)
	requestBodyUpdate.OrderStatus = helper.ORDER_DELIVERED

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/orders/%d/status", order.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", tokenCust)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.OrderResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	currentUser := GetCurrentUserByToken(t, tokenCust)
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
	assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
	assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
	assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
	assert.Equal(t, helper.ORDER_DELIVERED, responseBody.Data.OrderStatus)
	assert.Equal(t, true, responseBody.Data.IsDelivery)
	assert.Equal(t, float32(order.DeliveryCost), responseBody.Data.DeliveryCost)
	for _, address := range currentUser.Addresses {
		if address.IsMain {
			assert.Equal(t, address.Delivery.Cost, responseBody.Data.DeliveryCost)
			assert.Equal(t, address.CompleteAddress, responseBody.Data.CompleteAddress)
			break
		}
	}
	assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)

	assert.Equal(t, float32(100004), responseBody.Data.TotalProductPrice)
	assert.Equal(t, float32(99753.8), responseBody.Data.TotalFinalPrice)
	assert.Equal(t, 2, len(responseBody.Data.OrderProducts))
	for _, product := range responseBody.Data.OrderProducts {
		assert.NotNil(t, product.ID)
		assert.Equal(t, 2, product.Quantity)
	}

	// cek saldo
	currentUser = GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, helper.RoundFloat32((float32(150000)-responseBody.Data.TotalFinalPrice), 1), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)
}