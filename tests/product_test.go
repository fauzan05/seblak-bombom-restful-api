package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"seblak-bombom-restful-api/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1000")
	positions := []int{1, 2, 3}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 3; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestCreateProductImageSizeExceedsLimit(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1000")
	positions := []int{1, 2, 3}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 3; i++ {
		filename, content, err := GenerateDummyJPEG(2 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBody.Error, "file size is too large, maxium is 1MB")
}

// position is not equal to number of images
func TestCreateProductImagePositionBadRequest(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1000")
	positions := []int{1, 2, 3, 2, 9}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 3; i++ {
		filename, content, err := GenerateDummyJPEG(2 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBody.Error, "each uploaded image must have a corresponding position!")
}

func TestCreateProductFileNotImage(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1000")
	positions := []int{1, 2, 3}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 3; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "application/pdf")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBody.Error, "file type isn't valid, just only JPEG, PNG, JPG, WEBP and GIF are allowed")
}

func TestCreateProductImagesMoreThanFive(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1000")
	positions := []int{1, 2, 3, 4, 5, 6}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 6; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBody.Error, "you can upload up to 5 images only!")
}

func TestCreateProductImagesPositionNotIncluded(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1000")

	// Buat file image dummy
	for i := 1; i <= 6; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBody.Error, "image position must be included!")
}

func TestCreateProductStockNegativeNumber(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "-1")
	positions := []int{1, 2}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 2; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBody.Error, "stock must be positive number!")
}

func TestCreateProductImagesNotUploaded(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1")
	positions := []int{1, 2}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBody.Error, "images must be uploaded!")
}

func TestCreateProductCategoryNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", "1")
	_ = writer.WriteField("name", "Produk 1")
	_ = writer.WriteField("description", "Ini adalah produk 1")
	_ = writer.WriteField("price", "25000")
	_ = writer.WriteField("stock", "1")
	positions := []int{1, 2}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 2; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, responseBody.Error, "category not found!")
}

func TestCreateProductDataBadRequest(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", "0")
	_ = writer.WriteField("name", "")
	_ = writer.WriteField("description", "")
	_ = writer.WriteField("price", "")
	_ = writer.WriteField("stock", "0")

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/products", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, responseBody.Error, "invalid request body :")
}