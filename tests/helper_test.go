package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func ClearAll() {
	DeleteAllApplicationImages()
	ClearApplicationsSetting()
	ClearOrderProducts()
	ClearOrders()
	DeleteAllProductImages()
	ClearImages()
	ClearProducts()
	ClearCategories()
	ClearDiscountUsages()
	ClearDiscountCoupons()
	ClearTokens()
	ClearWallets()
	ClearAddresses()
	ClearDeliveries()
	ClearCarts()
	ClearUsers()
}

func ClearApplicationsSetting() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Application{}).Error
	if err != nil {
		log.Fatalf("Failed clear application settings data : %+v", err)
	}
}

func ClearImages() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Image{}).Error
	if err != nil {
		log.Fatalf("Failed clear images data : %+v", err)
	}
}

func ClearOrders() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Order{}).Error
	if err != nil {
		log.Fatalf("Failed clear orders data : %+v", err)
	}
}

func ClearOrderProducts() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.OrderProduct{}).Error
	if err != nil {
		log.Fatalf("Failed clear order products data : %+v", err)
	}
}

func ClearProducts() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Product{}).Error
	if err != nil {
		log.Fatalf("Failed clear products data : %+v", err)
	}
}

func ClearDiscountCoupons() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.DiscountCoupon{}).Error
	if err != nil {
		log.Fatalf("Failed clear discount coupons data : %+v", err)
	}
}

func ClearDiscountUsages() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.DiscountUsage{}).Error
	if err != nil {
		log.Fatalf("Failed clear discount usages data : %+v", err)
	}
}

func ClearTokens() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Token{}).Error
	if err != nil {
		log.Fatalf("Failed clear token data : %+v", err)
	}
}

func ClearCarts() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Cart{}).Error
	if err != nil {
		log.Fatalf("Failed clear cart data : %+v", err)
	}
}

func ClearCategories() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Category{}).Error
	if err != nil {
		log.Fatalf("Failed clear categories data : %+v", err)
	}
}

func ClearDeliveries() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Delivery{}).Error
	if err != nil {
		log.Fatalf("Failed clear delivery data : %+v", err)
	}
}

func ClearAddresses() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Address{}).Error
	if err != nil {
		log.Fatalf("Failed clear address data : %+v", err)
	}
}

func ClearWallets() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.Wallet{}).Error
	if err != nil {
		log.Fatalf("Failed clear wallet data : %+v", err)
	}
}

func ClearUsers() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.User{}).Error
	if err != nil {
		log.Fatalf("Failed clear user data : %+v", err)
	}
}

func DoLoginAdmin(t *testing.T) string {
	requestBody := model.LoginUserRequest{
		Email:    "johndoe@email.com",
		Password: "johndoe123",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserTokenResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)

	return responseBody.Data.Token
}

func DoLoginCustomer(t *testing.T) string {
	requestBody := model.LoginUserRequest{
		Email:    "customer1@email.com",
		Password: "customer1",
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserTokenResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.Token)

	return responseBody.Data.Token
}

func DoCreateDelivery(t *testing.T, token string) model.DeliveryResponse {
	requestBody := model.CreateDeliveryRequest{
		City:     "Kebumen",
		District: "Pejagoan",
		Village:  "Peniron",
		Hamlet:   "Jetis",
		Cost:     5000,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	request := httptest.NewRequest(http.MethodPost, "/api/deliveries", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.DeliveryResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, responseBody.Data.City, requestBody.City)
	assert.Equal(t, responseBody.Data.District, requestBody.District)
	assert.Equal(t, responseBody.Data.Village, requestBody.Village)
	assert.Equal(t, responseBody.Data.Hamlet, requestBody.Hamlet)
	assert.Equal(t, responseBody.Data.Cost, requestBody.Cost)

	return responseBody.Data
}

func DoCreateManyDelivery(t *testing.T, totalData int) string {
	token := DoLoginAdmin(t)
	totalIds := ""

	for i := 1; i <= totalData; i++ {
		cost := 5000 * float32(i)
		requestBody := model.CreateDeliveryRequest{
			City:     fmt.Sprintf("Kebumen %+v", i),
			District: fmt.Sprintf("Pejagoan %+v", i),
			Village:  fmt.Sprintf("Peniron %+v", i),
			Hamlet:   fmt.Sprintf("Jetis %+v", i),
			Cost:     cost,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)

		request := httptest.NewRequest(http.MethodPost, "/api/deliveries", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request)
		assert.Nil(t, err)

		bytes, err := io.ReadAll(response.Body)
		assert.Nil(t, err)

		responseBody := new(model.ApiResponse[model.DeliveryResponse])
		err = json.Unmarshal(bytes, responseBody)
		assert.Nil(t, err)

		assert.Equal(t, http.StatusCreated, response.StatusCode)
		assert.Equal(t, responseBody.Data.City, requestBody.City)
		assert.Equal(t, responseBody.Data.District, requestBody.District)
		assert.Equal(t, responseBody.Data.Village, requestBody.Village)
		assert.Equal(t, responseBody.Data.Hamlet, requestBody.Hamlet)
		assert.Equal(t, responseBody.Data.Cost, requestBody.Cost)

		convertIdToString := strconv.Itoa(int(responseBody.Data.ID))
		totalIds += convertIdToString + ","
	}

	return totalIds
}

func DoCreateManyDiscountCoupon(t *testing.T, token string, totalData int, returnDataByI int) *model.DiscountCouponResponse {
	start := "2025-01-01T00:00:01Z"
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := "2025-12-30T23:59:59Z"
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)

	// Ubah waktu ke WIB
	startWIB := parseStart.Local()
	endWIB := parseEnd.Local()

	var getDiscountCoupon *model.DiscountCouponResponse
	for i := 1; i <= totalData; i++ {
		requestBody := model.CreateDiscountCouponRequest{
			Name:            fmt.Sprintf("Diskon %+v", i),
			Description:     fmt.Sprintf("Discount Description %+v", i),
			Code:            fmt.Sprintf("ABC%+v", i),
			Value:           15,
			Type:            helper.PERCENT,
			Start:           helper.TimeRFC3339(startWIB),
			End:             helper.TimeRFC3339(endWIB),
			TotalMaxUsage:   100,
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
		assert.Equal(t, requestBody.Start.ToTime(), responseBody.Data.Start.ToTime())
		assert.Equal(t, requestBody.End.ToTime(), responseBody.Data.End.ToTime())
		assert.Equal(t, requestBody.TotalMaxUsage, responseBody.Data.TotalMaxUsage)
		assert.Equal(t, requestBody.MaxUsagePerUser, responseBody.Data.MaxUsagePerUser)
		assert.Equal(t, requestBody.UsedCount, responseBody.Data.UsedCount)
		assert.Equal(t, requestBody.MinOrderValue, responseBody.Data.MinOrderValue)
		assert.Equal(t, requestBody.Status, responseBody.Data.Status)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		if i == returnDataByI {
			getDiscountCoupon = &responseBody.Data
		}
	}

	return getDiscountCoupon
}

func DoCreateCategory(t *testing.T, token string, categoryName string, categoryDesc string) *model.CategoryResponse {
	requestBody := model.CreateProductRequest{
		Name:        categoryName,
		Description: categoryDesc,
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
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.Equal(t, requestBody.Description, responseBody.Data.Description)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	return &responseBody.Data
}

func GenerateDummyJPEG(sizeInBytes int) (filename string, content []byte, err error) {
	// Estimasi kasar dimensi berdasarkan ukuran (JPEG compress, jadi ini bukan ukuran pasti)
	// Makin besar dimensi, makin besar file
	scale := int(math.Sqrt(float64(sizeInBytes) / 3)) // 3 byte per pixel (RGB)
	width := scale
	height := scale

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Warna putih polos
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 10}) // Bisa ganti quality
	if err != nil {
		return "", nil, err
	}

	// Tambahkan padding kalau masih kurang
	for buf.Len() < sizeInBytes {
		buf.WriteByte(0) // dummy padding
	}

	filename = fmt.Sprintf("dummy_%d.jpg", sizeInBytes)
	return filename, buf.Bytes(), nil
}

func DeleteAllProductImages() {
	folderPath := "../uploads/images/products/"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			err := os.Remove(filepath.Join(folderPath, file.Name()))
			if err != nil {
				log.Printf("Failed to delete file %s: %v", file.Name(), err)
			}
		}
	}
}

func DeleteAllApplicationImages() {
	folderPath := "../uploads/images/application/"

	files, err := os.ReadDir(folderPath)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			err := os.Remove(filepath.Join(folderPath, file.Name()))
			if err != nil {
				log.Printf("Failed to delete file %s: %v", file.Name(), err)
			}
		}
	}
}

func DoCreateProduct(t *testing.T, token string, totalData int, getProductByIndex int) model.ProductResponse {
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	var getProduct model.ProductResponse
	for i := 1; i <= totalData; i++ {
		// Simulasi multipart body
		var b bytes.Buffer
		writer := multipart.NewWriter(&b)

		// Tambahkan field JSON sebagai string field biasa
		_ = writer.WriteField("category_id", fmt.Sprintf("%d", createCategory.ID))
		_ = writer.WriteField("name", fmt.Sprintf("Produk %d", i))
		_ = writer.WriteField("description", fmt.Sprintf("Ini adalah produk %d", i))
		_ = writer.WriteField("price", fmt.Sprintf("2500%d", i))
		_ = writer.WriteField("stock", fmt.Sprintf("100%d", i))
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
		assert.Equal(t, fmt.Sprintf("Produk %d", i), responseBody.Data.Name)
		assert.Equal(t, fmt.Sprintf("Ini adalah produk %d", i), responseBody.Data.Description)
		convertPriceToFloat32, err := strconv.Atoi(fmt.Sprintf("2500%d", i))
		assert.Nil(t, err)
		assert.Equal(t, float32(convertPriceToFloat32), responseBody.Data.Price)
		convertStockToInt, err := strconv.Atoi(fmt.Sprintf("100%d", i))
		assert.Nil(t, err)
		assert.Equal(t, convertStockToInt, responseBody.Data.Stock)
		assert.NotNil(t, responseBody.Data.CreatedAt)
		assert.NotNil(t, responseBody.Data.UpdatedAt)

		if i == getProductByIndex {
			getProduct = responseBody.Data
		}
	}

	return getProduct
}
