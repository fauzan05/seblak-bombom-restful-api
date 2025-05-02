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
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.Equal(t, 3, len(responseBody.Data.Images))
	for _, image := range responseBody.Data.Images {
		assert.NotNil(t, image.ID)
		assert.NotNil(t, image.FileName)
		assert.NotNil(t, image.Type)
		assert.NotNil(t, image.CreatedAt)
		assert.NotNil(t, image.UpdatedAt)
	}
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

func TestCreateProductPriceNegativeNumber(t *testing.T) {
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
	_ = writer.WriteField("price", "-1")
	_ = writer.WriteField("stock", "0")
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
	assert.Equal(t, responseBody.Error, "price must be positive number!")
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
	assert.Equal(t, responseBody.Error, "category is not found!")
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

func TestUpdateProduct(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	createCategory = DoCreateCategory(t, token, "Minuman", "Ini adalah minuman")
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1 Update")
	_ = writer.WriteField("description", "Ini adalah produk 1 Update")
	_ = writer.WriteField("price", "27000")
	_ = writer.WriteField("stock", "100")

	// CURRENT IMAGES
	for i, curImg := range responseBody.Data.Images {
		if i == 0 || i == 1 {
			_ = writer.WriteField("images_deleted", fmt.Sprintf("%d", curImg.ID))
		} else {
			_ = writer.WriteField("current_images", fmt.Sprintf("%d", curImg.ID))
			_ = writer.WriteField("current_positions", fmt.Sprintf("%d", curImg.Position))
		}
	}

	// NEW IMAGES
	for i := 1; i <= 2; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="new_images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	newPositions := []int{1, 2}
	for _, pos := range newPositions {
		_ = writer.WriteField("new_positions", fmt.Sprintf("%d", pos))
	}

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody = new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesUpdate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1 Update", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1 Update", responseBody.Data.Description)
	assert.Equal(t, float32(27000), responseBody.Data.Price)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, 100, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateProductCategoryNotFound(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", "1")
	_ = writer.WriteField("name", "Produk 1 Update")
	_ = writer.WriteField("description", "Ini adalah produk 1 Update")
	_ = writer.WriteField("price", "27000")
	_ = writer.WriteField("stock", "100")

	// CURRENT IMAGES
	for i, curImg := range responseBody.Data.Images {
		if i == 0 || i == 1 {
			_ = writer.WriteField("images_deleted", fmt.Sprintf("%d", curImg.ID))
		} else {
			_ = writer.WriteField("current_images", fmt.Sprintf("%d", curImg.ID))
			_ = writer.WriteField("current_positions", fmt.Sprintf("%d", curImg.Position))
		}
	}

	// NEW IMAGES
	for i := 1; i <= 2; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="new_images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	newPositions := []int{1, 2}
	for _, pos := range newPositions {
		_ = writer.WriteField("new_positions", fmt.Sprintf("%d", pos))
	}

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, "category is not found!", responseBodyUpdate.Error)
}

func TestUpdateProductStockNegativeNumber(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", "1")
	_ = writer.WriteField("name", "Seblak 1 Update")
	_ = writer.WriteField("description", "Desc 1 Update")
	_ = writer.WriteField("price", "0")
	_ = writer.WriteField("stock", "-1")

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBodyUpdate.Error, "stock must be positive number!")
}

func TestUpdateProductPriceNegativeNumber(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", "1")
	_ = writer.WriteField("name", "Seblak 1 Update")
	_ = writer.WriteField("description", "Desc 1 Update")
	_ = writer.WriteField("price", "-1")
	_ = writer.WriteField("stock", "0")

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBodyUpdate.Error, "price must be positive number!")
}

func TestUpdateProductDataBadRequest(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", "0")
	_ = writer.WriteField("name", "")
	_ = writer.WriteField("description", "")
	_ = writer.WriteField("price", "0")
	_ = writer.WriteField("stock", "0")

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Contains(t, responseBodyUpdate.Error, "invalid request body :")
}

func TestUpdateProductNotFound(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	createCategory = DoCreateCategory(t, token, "Minuman", "Ini adalah minuman")
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1 Update")
	_ = writer.WriteField("description", "Ini adalah produk 1 Update")
	_ = writer.WriteField("price", "27000")
	_ = writer.WriteField("stock", "100")

	// CURRENT IMAGES
	for i, curImg := range responseBody.Data.Images {
		if i == 0 || i == 1 {
			_ = writer.WriteField("images_deleted", fmt.Sprintf("%d", curImg.ID))
		} else {
			_ = writer.WriteField("current_images", fmt.Sprintf("%d", curImg.ID))
			_ = writer.WriteField("current_positions", fmt.Sprintf("%d", curImg.Position))
		}
	}

	// NEW IMAGES
	for i := 1; i <= 2; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="new_images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	newPositions := []int{1, 2}
	for _, pos := range newPositions {
		_ = writer.WriteField("new_positions", fmt.Sprintf("%d", pos))
	}

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", 1), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	assert.Equal(t, responseBodyUpdate.Error, "product is not found!")
}

func TestUpdateProductNewImagesIsNotSameWithNewImagePositions(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	createCategory = DoCreateCategory(t, token, "Minuman", "Ini adalah minuman")
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1 Update")
	_ = writer.WriteField("description", "Ini adalah produk 1 Update")
	_ = writer.WriteField("price", "27000")
	_ = writer.WriteField("stock", "100")

	// CURRENT IMAGES
	for i, curImg := range responseBody.Data.Images {
		if i == 0 || i == 1 {
			_ = writer.WriteField("images_deleted", fmt.Sprintf("%d", curImg.ID))
		} else {
			_ = writer.WriteField("current_images", fmt.Sprintf("%d", curImg.ID))
			_ = writer.WriteField("current_positions", fmt.Sprintf("%d", curImg.Position))
		}
	}

	// NEW IMAGES
	for i := 1; i <= 2; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="new_images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	newPositions := []int{1, 2, 3}
	for _, pos := range newPositions {
		_ = writer.WriteField("new_positions", fmt.Sprintf("%d", pos))
	}

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBodyUpdate.Error, "each new uploaded image must have a corresponding position!")
}

func TestUpdateProductImageMoreThanFive(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	createCategory = DoCreateCategory(t, token, "Minuman", "Ini adalah minuman")
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1 Update")
	_ = writer.WriteField("description", "Ini adalah produk 1 Update")
	_ = writer.WriteField("price", "27000")
	_ = writer.WriteField("stock", "100")

	// CURRENT IMAGES
	for i, curImg := range responseBody.Data.Images {
		if i == 0 || i == 1 {
			_ = writer.WriteField("images_deleted", fmt.Sprintf("%d", curImg.ID))
		} else {
			_ = writer.WriteField("current_images", fmt.Sprintf("%d", curImg.ID))
			_ = writer.WriteField("current_positions", fmt.Sprintf("%d", curImg.Position))
		}
	}

	// NEW IMAGES
	for i := 1; i <= 3; i++ {
		filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="new_images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	newPositions := []int{1, 2, 3}
	for _, pos := range newPositions {
		_ = writer.WriteField("new_positions", fmt.Sprintf("%d", pos))
	}

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBodyUpdate.Error, "you can upload up to 5 images only!")
}

func TestUpdateProductNewImageSizeExceedsLimit(t *testing.T) {
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
	positions := []int{1, 2, 3, 4, 5}
	for _, pos := range positions {
		_ = writer.WriteField("positions", fmt.Sprintf("%d", pos))
	}

	// Buat file image dummy
	for i := 1; i <= 5; i++ {
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

	bytesCreate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytesCreate, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
	assert.Equal(t, "Produk 1", responseBody.Data.Name)
	assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
	assert.Equal(t, 5, len(responseBody.Data.Images))
	assert.Equal(t, float32(25000), responseBody.Data.Price)
	assert.Equal(t, 1000, responseBody.Data.Stock)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	// UPDATE
	var c bytes.Buffer
	createCategory = DoCreateCategory(t, token, "Minuman", "Ini adalah minuman")
	// Simulasi multipart body
	writer = multipart.NewWriter(&c)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
	_ = writer.WriteField("name", "Produk 1 Update")
	_ = writer.WriteField("description", "Ini adalah produk 1 Update")
	_ = writer.WriteField("price", "27000")
	_ = writer.WriteField("stock", "100")

	// CURRENT IMAGES
	for i, curImg := range responseBody.Data.Images {
		if i == 0 || i == 1 {
			_ = writer.WriteField("images_deleted", fmt.Sprintf("%d", curImg.ID))
		} else {
			_ = writer.WriteField("current_images", fmt.Sprintf("%d", curImg.ID))
			_ = writer.WriteField("current_positions", fmt.Sprintf("%d", curImg.Position))
		}
	}

	// NEW IMAGES
	for i := 1; i <= 2; i++ {
		filename, content, err := GenerateDummyJPEG(2 * 1024 * 1024) // 1 MB
		assert.Nil(t, err)

		partHeader := textproto.MIMEHeader{}
		partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="new_images"; filename="%s"`, filename))
		partHeader.Set("Content-Type", "image/jpeg")

		fileWriter, err := writer.CreatePart(partHeader)
		assert.Nil(t, err)

		_, err = fileWriter.Write(content)
		assert.Nil(t, err)
	}

	newPositions := []int{1, 2}
	for _, pos := range newPositions {
		_ = writer.WriteField("new_positions", fmt.Sprintf("%d", pos))
	}

	writer.Close()

	request = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/products/%+v", responseBody.Data.ID), &c)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", token)

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytesUpdate, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyUpdate := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytesUpdate, responseBodyUpdate)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	assert.Equal(t, responseBodyUpdate.Error, "file size is too large, maxium is 1MB")
}

func TestGetProductPagination(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateProduct(t, token, 27, 0)

	request := httptest.NewRequest(http.MethodGet, "/api/products?per_page=5&page=2&search=produk&column=products.id&sort_by=desc", nil)
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
	assert.Equal(t, 5, len(*responseBody.Data))
	assert.Equal(t, int64(27), responseBody.TotalDatas)
	assert.Equal(t, 6, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	assert.Equal(t, 5, responseBody.DataPerPages)
}

func TestGetProductPaginationSortingDesc(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateProduct(t, token, 27, 0)

	request := httptest.NewRequest(http.MethodGet, "/api/products?per_page=5&page=2&search=produk&column=products.id&sort_by=desc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.ProductResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(*responseBody.Data))
	assert.Equal(t, int64(27), responseBody.TotalDatas)
	assert.Equal(t, 6, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	assert.Equal(t, 5, responseBody.DataPerPages)

	for _, product := range *responseBody.Data {
		assert.NotEmpty(t, product.Images)
		assert.NotEmpty(t, product.Category)
	}

	products := *responseBody.Data
	for i := range len(products) - 1 {
		assert.Greater(t, products[i].ID, products[i+1].ID)
	}
}

func TestGetProductPaginationSortingDescColumnNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateProduct(t, token, 27, 0)

	request := httptest.NewRequest(http.MethodGet, "/api/products?per_page=5&page=2&search=produk&column=products.lala&sort_by=desc", nil)
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
	assert.Equal(t, "invalid sort column : products.lala", responseBody.Error)
}

func TestGetProductPaginationSortingAsc(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateProduct(t, token, 27, 0)

	request := httptest.NewRequest(http.MethodGet, "/api/products?per_page=5&page=2&search=produk&column=products.id&sort_by=asc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.ProductResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 5, len(*responseBody.Data))
	assert.Equal(t, int64(27), responseBody.TotalDatas)
	assert.Equal(t, 6, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	assert.Equal(t, 5, responseBody.DataPerPages)

	for _, product := range *responseBody.Data {
		assert.NotEmpty(t, product.Images)
		assert.NotEmpty(t, product.Category)
	}

	products := *responseBody.Data
	for i := range len(products) - 1 {
		assert.Less(t, products[i].ID, products[i+1].ID)
	}
}

func TestGetProductPaginationSearchEmptyResult(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	DoCreateProduct(t, token, 27, 0)

	request := httptest.NewRequest(http.MethodGet, "/api/products?per_page=5&page=2&search=lala&column=products.id&sort_by=desc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponsePagination[*[]model.ProductResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 0, len(*responseBody.Data))
	assert.Equal(t, int64(0), responseBody.TotalDatas)
	assert.Equal(t, 0, responseBody.TotalPages)
	assert.Equal(t, 2, responseBody.CurrentPages)
	assert.Equal(t, 5, responseBody.DataPerPages)
}

func TestGetProductById(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)
	getProduct := DoCreateProduct(t, token, 5, 1)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/%d", getProduct.ID), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ProductResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, responseBody.Data.Name, getProduct.Name)
	assert.Equal(t, responseBody.Data.Description, getProduct.Description)
	assert.Equal(t, responseBody.Data.Category.ID, getProduct.Category.ID)
	assert.Equal(t, responseBody.Data.Category.Name, getProduct.Category.Name)
	assert.Equal(t, responseBody.Data.Category.Description, getProduct.Category.Description)
	assert.Equal(t, responseBody.Data.Category.CreatedAt, getProduct.Category.CreatedAt)
	assert.Equal(t, responseBody.Data.Category.UpdatedAt, getProduct.Category.UpdatedAt)
	assert.Equal(t, responseBody.Data.Price, getProduct.Price)
	assert.Equal(t, responseBody.Data.Stock, getProduct.Stock)
	assert.Equal(t, responseBody.Data.CreatedAt, getProduct.CreatedAt)
	assert.Equal(t, responseBody.Data.UpdatedAt, getProduct.UpdatedAt)
	for i, image := range responseBody.Data.Images {
		assert.Equal(t, getProduct.Images[i].ID, image.ID)
		assert.Equal(t, getProduct.Images[i].FileName, image.FileName)
		assert.Equal(t, getProduct.Images[i].Position, image.Position)
		assert.Equal(t, getProduct.Images[i].ProductId, image.ProductId)
		assert.Equal(t, getProduct.Images[i].Type, image.Type)
		assert.Equal(t, getProduct.Images[i].CreatedAt, image.CreatedAt)
		assert.Equal(t, getProduct.Images[i].UpdatedAt, image.UpdatedAt)
	}
}

func TestGetProductByIdNotFound(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)

	request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/products/%d", 1), nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ErrorResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

func TestDeleteProduct(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	var getAllIds string
	for i := 1; i <= 5; i++ {
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
		assert.NotNil(t, responseBody.Data.ID)
		assert.Equal(t, createCategory.ID, responseBody.Data.Category.ID)
		assert.Equal(t, "Produk 1", responseBody.Data.Name)
		assert.Equal(t, "Ini adalah produk 1", responseBody.Data.Description)
		assert.Equal(t, float32(25000), responseBody.Data.Price)
		assert.Equal(t, 1000, responseBody.Data.Stock)
		assert.Equal(t, 3, len(responseBody.Data.Images))
		for _, image := range responseBody.Data.Images {
			assert.NotNil(t, image.ID)
			assert.NotNil(t, image.FileName)
			assert.NotNil(t, image.Type)
			assert.NotNil(t, image.CreatedAt)
			assert.NotNil(t, image.UpdatedAt)
		}
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)
		getAllIds += fmt.Sprintf("%d,", responseBody.Data.ID)
	}

	request := httptest.NewRequest(http.MethodDelete, "/api/products?ids="+getAllIds, nil)
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

	// cek apakah produk masih ada
	request = httptest.NewRequest(http.MethodGet, "/api/products?per_page=5&page=2&column=products.id&sort_by=desc", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err = app.Test(request)
	assert.Nil(t, err)

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBodyPagination := new(model.ApiResponsePagination[*[]model.ProductResponse])
	err = json.Unmarshal(bytes, responseBodyPagination)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, 0, len(*responseBodyPagination.Data))
	assert.Equal(t, int64(0), responseBodyPagination.TotalDatas)
	assert.Equal(t, 0, responseBodyPagination.TotalPages)
	assert.Equal(t, 2, responseBodyPagination.CurrentPages)
	assert.Equal(t, 5, responseBodyPagination.DataPerPages)
}

func TestDeleteProductIdNotValid(t *testing.T) {
	ClearAll()
	TestRegisterAdmin(t)
	token := DoLoginAdmin(t)

	var getAllIds string = "b,3,#,m"
	request := httptest.NewRequest(http.MethodDelete, "/api/products?ids="+getAllIds, nil)
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
	assert.Contains(t, responseBody.Error, "invalid product ID :")
}
