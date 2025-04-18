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
