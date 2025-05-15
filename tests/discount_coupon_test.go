package tests

import (
	"encoding/json"
	"fmt"
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

func TestCreateDiscountCoupon(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-01-01T00:00:01Z"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-12-30T23:59:59Z"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)

	for i := 1; i <= 5; i++ {
		requestBody := model.CreateDiscountCouponRequest{
			Name:            fmt.Sprintf("Diskon %+v", i),
			Description:     fmt.Sprintf("Discount Description %+v", i),
			Code:            fmt.Sprintf("ABC%+v", i),
			Value:           15,
			Type:            helper.PERCENT,
			Start:           helper.TimeRFC3339(parseStart),
			End:             helper.TimeRFC3339(parseEnd),
			MaxUsagePerUser: 5,
			UsedCount:       0,
			MinOrderValue:   20000,
			Status:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/discount-coupons", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.DiscountCouponResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.Name, responseBody.Data.Name)
		assert.Equal(t, requestBody.Description, responseBody.Data.Description)
		assert.Equal(t, requestBody.Code, responseBody.Data.Code)
		assert.Equal(t, requestBody.Value, responseBody.Data.Value)
		assert.Equal(t, requestBody.Type, responseBody.Data.Type)
		assert.Equal(t, requestBody.Start, responseBody.Data.Start)
		assert.Equal(t, requestBody.End, responseBody.Data.End)
		assert.Equal(t, requestBody.MaxUsagePerUser, responseBody.Data.MaxUsagePerUser)
		assert.Equal(t, requestBody.UsedCount, responseBody.Data.UsedCount)
		assert.Equal(t, requestBody.MinOrderValue, responseBody.Data.MinOrderValue)
		assert.Equal(t, requestBody.Status, responseBody.Data.Status)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)
	}
}

func TestCreateDiscountCouponFailed(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-01-01T00:00:01Z"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-12-30T23:59:59Z"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)

	for i := 1; i <= 5; i++ {
		requestBody := model.CreateDiscountCouponRequest{
			Name:            "",
			Description:     "",
			Code:            "",
			Value:           15,
			Type:            helper.PERCENT,
			Start:           helper.TimeRFC3339(parseStart),
			End:             helper.TimeRFC3339(parseEnd),
			MaxUsagePerUser: 5,
			UsedCount:       0,
			MinOrderValue:   20000,
			Status:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/discount-coupons", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.DiscountCouponResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	}
}

func TestGetDiscountCouponPagination(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateManyDiscountCoupon(t, token, 27, 1)

	request := httptest.NewRequest(http.MethodGet, "/api/discount-coupons?search=diskon&column=discount_coupons.id&sort_by=asc&per_page=5&page=3", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 3, responseBody.CurrentPages)
	assert.Equal(t, 5, len(*responseBody.Data))
	assert.Equal(t, int64(27), responseBody.TotalDatas)
	assert.Equal(t, 5, responseBody.DataPerPages)
	assert.Equal(t, 6, responseBody.TotalPages)
}

func TestGetDiscountCouponPaginationSortingColumnNotFound(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateManyDiscountCoupon(t, token, 27, 1)

	request := httptest.NewRequest(http.MethodGet, "/api/discount-coupons?search=diskon&column=discount_coupons.mama&sort_by=asc&per_page=5&page=3", nil)
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
	assert.Equal(t, "invalid sort column : discount_coupons.mama", responseBody.Error)
}

func TestGetDiscountCouponById(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)
	getDiscountCoupon := DoCreateManyDiscountCoupon(t, token, 27, 1)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/discount-coupons/%+v", getDiscountCoupon.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, getDiscountCoupon.Name, responseBody.Data.Name)
	assert.Equal(t, getDiscountCoupon.Description, responseBody.Data.Description)
	assert.Equal(t, getDiscountCoupon.Code, responseBody.Data.Code)
	assert.Equal(t, getDiscountCoupon.Start.ToTime(), responseBody.Data.Start.ToTime())
	assert.Equal(t, getDiscountCoupon.End.ToTime(), responseBody.Data.End.ToTime())
	assert.Equal(t, getDiscountCoupon.MaxUsagePerUser, responseBody.Data.MaxUsagePerUser)
	assert.Equal(t, getDiscountCoupon.MinOrderValue, responseBody.Data.MinOrderValue)
	assert.Equal(t, getDiscountCoupon.Type, responseBody.Data.Type)
	assert.Equal(t, getDiscountCoupon.Value, responseBody.Data.Value)
	assert.Equal(t, getDiscountCoupon.UsedCount, responseBody.Data.UsedCount)
	assert.NotNil(t, responseBody.Data.CreatedAt.ToTime())
	assert.NotNil(t, responseBody.Data.UpdatedAt.ToTime())
}

func TestGetDiscountCouponByIdFailed(t *testing.T) {
	ClearAll()

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/discount-coupons/%+v", "e"), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestUpdateDiscountCouponById(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-01-01T00:00:01Z"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-12-30T23:59:59Z"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)

	requestBodyCreate := model.CreateDiscountCouponRequest{
		Name:            fmt.Sprintf("Diskon %+v", 1),
		Description:     fmt.Sprintf("Discount Description %+v", 1),
		Code:            fmt.Sprintf("ABC%+v", 1),
		Value:           15,
		Type:            helper.PERCENT,
		Start:           helper.TimeRFC3339(parseStart),
		End:             helper.TimeRFC3339(parseEnd),
		MaxUsagePerUser: 5,
		UsedCount:       0,
		MinOrderValue:   20000,
		Status:          true,
	}

	bodyJson, err := json.Marshal(requestBodyCreate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/discount-coupons", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyCreate := new(model.ApiResponse[model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBodyCreate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBodyCreate.Name, responseBodyCreate.Data.Name)
	assert.Equal(t, requestBodyCreate.Description, responseBodyCreate.Data.Description)
	assert.Equal(t, requestBodyCreate.Code, responseBodyCreate.Data.Code)
	assert.Equal(t, requestBodyCreate.Value, responseBodyCreate.Data.Value)
	assert.Equal(t, requestBodyCreate.Type, responseBodyCreate.Data.Type)
	assert.Equal(t, requestBodyCreate.Start, responseBodyCreate.Data.Start)
	assert.Equal(t, requestBodyCreate.End, responseBodyCreate.Data.End)
	assert.Equal(t, requestBodyCreate.MaxUsagePerUser, responseBodyCreate.Data.MaxUsagePerUser)
	assert.Equal(t, requestBodyCreate.UsedCount, responseBodyCreate.Data.UsedCount)
	assert.Equal(t, requestBodyCreate.MinOrderValue, responseBodyCreate.Data.MinOrderValue)
	assert.Equal(t, requestBodyCreate.Status, responseBodyCreate.Data.Status)
	assert.NotNil(t, responseBodyCreate.Data.CreatedAt)
	assert.NotNil(t, responseBodyCreate.Data.UpdatedAt)

	requestBodyUpdate := model.UpdateDiscountCouponRequest{
		Name:            fmt.Sprintf("Diskon %+v", 2),
		Description:     fmt.Sprintf("Discount Description %+v", 2),
		Code:            fmt.Sprintf("ABC%+v", 2),
		Value:           5000,
		Type:            helper.NOMINAL,
		Start:           helper.TimeRFC3339(parseStart),
		End:             helper.TimeRFC3339(parseEnd),
		MaxUsagePerUser: 3,
		UsedCount:       0,
		MinOrderValue:   25000,
		Status:          false,
	}

	bodyJson, err = json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/discount-coupons/%+v", responseBodyCreate.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ApiResponse[model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, requestBodyUpdate.Name, responseBodyUpdate.Data.Name)
	assert.Equal(t, requestBodyUpdate.Description, responseBodyUpdate.Data.Description)
	assert.Equal(t, requestBodyUpdate.Code, responseBodyUpdate.Data.Code)
	assert.Equal(t, requestBodyUpdate.Value, responseBodyUpdate.Data.Value)
	assert.Equal(t, requestBodyUpdate.Type, responseBodyUpdate.Data.Type)
	assert.Equal(t, requestBodyUpdate.Start, responseBodyUpdate.Data.Start)
	assert.Equal(t, requestBodyUpdate.End, responseBodyUpdate.Data.End)
	assert.Equal(t, requestBodyUpdate.MaxUsagePerUser, responseBodyUpdate.Data.MaxUsagePerUser)
	assert.Equal(t, requestBodyUpdate.UsedCount, responseBodyUpdate.Data.UsedCount)
	assert.Equal(t, requestBodyUpdate.MinOrderValue, responseBodyUpdate.Data.MinOrderValue)
	assert.Equal(t, requestBodyUpdate.Status, responseBodyUpdate.Data.Status)
}

func TestUpdateDiscountCouponByIdBadRequest(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-01-01T00:00:01Z"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-12-30T23:59:59Z"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)

	requestBodyCreate := model.CreateDiscountCouponRequest{
		Name:            fmt.Sprintf("Diskon %+v", 1),
		Description:     fmt.Sprintf("Discount Description %+v", 1),
		Code:            fmt.Sprintf("ABC%+v", 1),
		Value:           15,
		Type:            helper.PERCENT,
		Start:           helper.TimeRFC3339(parseStart),
		End:             helper.TimeRFC3339(parseEnd),
		MaxUsagePerUser: 5,
		UsedCount:       0,
		MinOrderValue:   20000,
		Status:          true,
	}

	bodyJson, err := json.Marshal(requestBodyCreate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/discount-coupons", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyCreate := new(model.ApiResponse[model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBodyCreate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBodyCreate.Name, responseBodyCreate.Data.Name)
	assert.Equal(t, requestBodyCreate.Description, responseBodyCreate.Data.Description)
	assert.Equal(t, requestBodyCreate.Code, responseBodyCreate.Data.Code)
	assert.Equal(t, requestBodyCreate.Value, responseBodyCreate.Data.Value)
	assert.Equal(t, requestBodyCreate.Type, responseBodyCreate.Data.Type)
	assert.Equal(t, requestBodyCreate.Start, responseBodyCreate.Data.Start)
	assert.Equal(t, requestBodyCreate.End, responseBodyCreate.Data.End)
	assert.Equal(t, requestBodyCreate.MaxUsagePerUser, responseBodyCreate.Data.MaxUsagePerUser)
	assert.Equal(t, requestBodyCreate.UsedCount, responseBodyCreate.Data.UsedCount)
	assert.Equal(t, requestBodyCreate.MinOrderValue, responseBodyCreate.Data.MinOrderValue)
	assert.Equal(t, requestBodyCreate.Status, responseBodyCreate.Data.Status)
	assert.NotNil(t, responseBodyCreate.Data.CreatedAt)
	assert.NotNil(t, responseBodyCreate.Data.UpdatedAt)

	requestBodyUpdate := model.UpdateDiscountCouponRequest{
		Name:            "",
		Description:     "",
		Code:            "",
		Value:           15,
		Type:            helper.PERCENT,
		Start:           helper.TimeRFC3339(parseStart),
		End:             helper.TimeRFC3339(parseEnd),
		MaxUsagePerUser: 5,
		UsedCount:       0,
		MinOrderValue:   20000,
		Status:          true,
	}

	bodyJson, err = json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/discount-coupons/%+v", responseBodyCreate.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ApiResponse[model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestUpdateDiscountCouponByIdNotFound(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-01-01T00:00:01Z"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-12-30T23:59:59Z"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)

	requestBodyUpdate := model.UpdateDiscountCouponRequest{
		Name:            fmt.Sprintf("Diskon %+v", 2),
		Description:     fmt.Sprintf("Discount Description %+v", 2),
		Code:            fmt.Sprintf("ABC%+v", 2),
		Value:           5000,
		Type:            helper.NOMINAL,
		Start:           helper.TimeRFC3339(parseStart),
		End:             helper.TimeRFC3339(parseEnd),
		MaxUsagePerUser: 3,
		UsedCount:       0,
		MinOrderValue:   25000,
		Status:          false,
	}

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/discount-coupons/%+v", -999), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ApiResponse[model.DiscountCouponResponse])
	err = json.Unmarshal(bytes, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeleteDiscountCoupon(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	start := "2025-01-01T00:00:01Z"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-12-30T23:59:59Z"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)

	var getAllIds string
	for i := 1; i <= 5; i++ {
		requestBody := model.CreateDiscountCouponRequest{
			Name:            fmt.Sprintf("Diskon %+v", i),
			Description:     fmt.Sprintf("Discount Description %+v", i),
			Code:            fmt.Sprintf("ABC%+v", i),
			Value:           15,
			Type:            helper.PERCENT,
			Start:           helper.TimeRFC3339(parseStart),
			End:             helper.TimeRFC3339(parseEnd),
			MaxUsagePerUser: 5,
			UsedCount:       0,
			MinOrderValue:   20000,
			Status:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/discount-coupons", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.DiscountCouponResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, requestBody.Name, responseBody.Data.Name)
		assert.Equal(t, requestBody.Description, responseBody.Data.Description)
		assert.Equal(t, requestBody.Code, responseBody.Data.Code)
		assert.Equal(t, requestBody.Value, responseBody.Data.Value)
		assert.Equal(t, requestBody.Type, responseBody.Data.Type)
		assert.Equal(t, requestBody.Start, responseBody.Data.Start)
		assert.Equal(t, requestBody.End, responseBody.Data.End)
		assert.Equal(t, requestBody.MaxUsagePerUser, responseBody.Data.MaxUsagePerUser)
		assert.Equal(t, requestBody.UsedCount, responseBody.Data.UsedCount)
		assert.Equal(t, requestBody.MinOrderValue, responseBody.Data.MinOrderValue)
		assert.Equal(t, requestBody.Status, responseBody.Data.Status)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		getAllIds += fmt.Sprintf("%+v,", responseBody.Data.ID)
	}

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/discount-coupons?ids=%+v", getAllIds), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestDeleteDiscountCouponIdsNotValid(t *testing.T) {
	ClearAll()
	DoRegisterAdmin(t)
	token := DoLoginAdmin(t)

	var getAllIds string = "b,s[];.,asd"

	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/discount-coupons?ids=%+v", getAllIds), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ApiResponse[bool])
	err = json.Unmarshal(bytes, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}
