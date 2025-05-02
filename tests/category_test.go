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

func TestCreateCategory(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.CreateCategoryRequest{
		Name:        "Makanan",
		Description: "Ini adalah makanan",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Description, responseBody.Data.Description)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestCreateCategoryBadRequest(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.CreateCategoryRequest{
		Name:        "",
		Description: "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	requestBody = model.CreateCategoryRequest{
		Name:        "There are many variations of passages of Lorem Ipsum available, but the majority have suffered alteration in some form, by injected humour, or randomised words which don't look even slightly believable. If you are going to use a passage of Lorem Ipsum, you need to be sure there isn't anything embarrassing hidden in the middle of text. All the Lorem Ipsum generators on the Internet tend to repeat predefined chunks as necessary, making this the first true generator on the Internet. It uses a dictionary of over 200 Latin words, combined with a handful of model sentence structures, to generate Lorem Ipsum which looks reasonable. The generated Lorem Ipsum is therefore always free from repetition, injected humour, or non-characteristic words etc.",
		Description: "",
	}

	bodyJson, err = json.Marshal(requestBody)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateCategory(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBodyCreate := model.CreateCategoryRequest{
		Name:        "Makanan",
		Description: "Ini adalah makanan",
	}

	bodyJson, err := json.Marshal(requestBodyCreate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBodyCreate.Name, responseBody.Data.Name)
	assert.Equal(t, requestBodyCreate.Description, responseBody.Data.Description)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	requestBodyUpdate := model.CreateCategoryRequest{
		Name:        "Minuman",
		Description: "Ini adalah minuman",
	}

	bodyJson, err = json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/categories/%+v", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBodyUpdate.Name, responseBody.Data.Name)
	assert.NotEqual(t, requestBodyCreate.Name, responseBody.Data.Name)
	assert.Equal(t, requestBodyUpdate.Description, responseBody.Data.Description)
	assert.NotEqual(t, requestBodyCreate.Description, responseBody.Data.Description)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateCategoryBadRequest(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBodyCreate := model.CreateCategoryRequest{
		Name:        "Makanan",
		Description: "Ini adalah makanan",
	}

	bodyJson, err := json.Marshal(requestBodyCreate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, requestBodyCreate.Name, responseBody.Data.Name)
	assert.Equal(t, requestBodyCreate.Description, responseBody.Data.Description)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	requestBodyUpdate := model.UpdateCategoryRequest{
		ID:          responseBody.Data.ID,
		Name:        "",
		Description: "",
	}

	bodyJson, err = json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/categories/%+v", responseBody.Data.ID), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestUpdateCategoryNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBodyUpdate := model.UpdateCategoryRequest{
		Name:        "Category New",
		Description: "Desc NEw",
	}

	bodyJson, err := json.Marshal(requestBodyUpdate)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/categories/%+v", -9), strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestGetCategoryById(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	requestBody := model.CreateCategoryRequest{
		Name:        "Makanan",
		Description: "Ini adalah makanan",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyCreate := new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBodyCreate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBodyCreate.Data.ID)
	assert.Equal(t, requestBody.Name, responseBodyCreate.Data.Name)
	assert.Equal(t, requestBody.Description, responseBodyCreate.Data.Description)
	assert.NotNil(t, responseBodyCreate.Data.CreatedAt)
	assert.NotNil(t, responseBodyCreate.Data.UpdatedAt)

	// get by id
	request = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/categories/%+v", responseBodyCreate.Data.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, responseBodyCreate.Data.ID, responseBody.Data.ID)
	assert.Equal(t, responseBodyCreate.Data.Name, responseBody.Data.Name)
	assert.Equal(t, responseBodyCreate.Data.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, responseBodyCreate.Data.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetAllCategoryPagination(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	for i := 1; i <= 25; i++ {
		requestBody := model.CreateCategoryRequest{
			Name:        fmt.Sprintf("Makanan %+v", i),
			Description: fmt.Sprintf("Ini adalah makanan %+v", i),
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBodyCreate := new(model.ApiResponse[model.CategoryResponse])
		err = json.Unmarshal(bytes, responseBodyCreate)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotNil(t, responseBodyCreate.Data.ID)
		assert.Equal(t, requestBody.Name, responseBodyCreate.Data.Name)
		assert.Equal(t, requestBody.Description, responseBodyCreate.Data.Description)
		assert.NotNil(t, responseBodyCreate.Data.CreatedAt)
		assert.NotNil(t, responseBodyCreate.Data.UpdatedAt)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/categories?per_page=10&page=2&search=makanan&column=id&sort_by=desc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 10, len(*responseBody.Data))
	assert.Equal(t, int64(25), responseBody.TotalDatas)
	assert.Equal(t, 3, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	assert.Equal(t, 10, responseBody.DataPerPages)
}

func TestGetAllCategoryPaginationSearchNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	for i := 1; i <= 25; i++ {
		requestBody := model.CreateCategoryRequest{
			Name:        fmt.Sprintf("Makanan %+v", i),
			Description: fmt.Sprintf("Ini adalah makanan %+v", i),
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBodyCreate := new(model.ApiResponse[model.CategoryResponse])
		err = json.Unmarshal(bytes, responseBodyCreate)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotNil(t, responseBodyCreate.Data.ID)
		assert.Equal(t, requestBody.Name, responseBodyCreate.Data.Name)
		assert.Equal(t, requestBody.Description, responseBodyCreate.Data.Description)
		assert.NotNil(t, responseBodyCreate.Data.CreatedAt)
		assert.NotNil(t, responseBodyCreate.Data.UpdatedAt)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/categories?per_page=10&page=1&search=zzz&column=id&sort_by=desc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 0, len(*responseBody.Data))
	assert.Equal(t, int64(0), responseBody.TotalDatas)
	assert.Equal(t, 0, responseBody.TotalPages)
	assert.Equal(t, 1, responseBody.CurrentPages)
	assert.Equal(t, 10, responseBody.DataPerPages)
}

func TestGetAllCategoryPaginationSortingColumn(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	for i := 1; i <= 25; i++ {
		requestBody := model.CreateCategoryRequest{
			Name:        fmt.Sprintf("Makanan %+v", i),
			Description: fmt.Sprintf("Ini adalah makanan %+v", i),
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBodyCreate := new(model.ApiResponse[model.CategoryResponse])
		err = json.Unmarshal(bytes, responseBodyCreate)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotNil(t, responseBodyCreate.Data.ID)
		assert.Equal(t, requestBody.Name, responseBodyCreate.Data.Name)
		assert.Equal(t, requestBody.Description, responseBodyCreate.Data.Description)
		assert.NotNil(t, responseBodyCreate.Data.CreatedAt)
		assert.NotNil(t, responseBodyCreate.Data.UpdatedAt)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/categories?per_page=10&page=1&search=makanan&column=categories.name&sort_by=desc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.CategoryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 10, len(*responseBody.Data))
	assert.Equal(t, int64(25), responseBody.TotalDatas)
	assert.Equal(t, 3, responseBody.TotalPages)
	assert.Equal(t, 1, responseBody.CurrentPages)
	assert.Equal(t, 10, responseBody.DataPerPages)
}

func TestDeleteCategoryById(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	var getAllIds string
	for i := 1; i <= 25; i++ {
		requestBody := model.CreateCategoryRequest{
			Name:        fmt.Sprintf("Makanan %+v", i),
			Description: fmt.Sprintf("Ini adalah makanan %+v", i),
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/categories", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBodyCreate := new(model.ApiResponse[model.CategoryResponse])
		err = json.Unmarshal(bytes, responseBodyCreate)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.NotNil(t, responseBodyCreate.Data.ID)
		assert.Equal(t, requestBody.Name, responseBodyCreate.Data.Name)
		assert.Equal(t, requestBody.Description, responseBodyCreate.Data.Description)
		assert.NotNil(t, responseBodyCreate.Data.CreatedAt)
		assert.NotNil(t, responseBodyCreate.Data.UpdatedAt)
		getAllIds += fmt.Sprintf("%+v,", responseBodyCreate.Data.ID)
	}

	request := httptest.NewRequest(http.MethodDelete, "/api/categories?ids="+getAllIds, nil)
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

func TestDeleteCategoriesIdNotValid(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	getAllIds := "3,s,r,t,"
	request := httptest.NewRequest(http.MethodDelete, "/api/categories?ids="+getAllIds, nil)
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
