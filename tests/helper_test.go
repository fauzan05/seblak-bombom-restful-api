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
	ClearDiscountCouponUsages()
	ClearXenditTransactions()
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

func ClearDiscountCouponUsages() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.DiscountUsage{}).Error
	if err != nil {
		log.Fatalf("Failed clear discount usages data : %+v", err)
	}
}

func ClearXenditTransactions() {
	err := db.Unscoped().Where("1 = 1").Delete(&entity.XenditTransactions{}).Error
	if err != nil {
		log.Fatalf("Failed clear xendit transactions data : %+v", err)
	}
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

	response, err := app.Test(request, int(time.Second) * 5)
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

	response, err := app.Test(request, int(time.Second) * 5)
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

func DoCreateDelivery(t *testing.T, token string) *model.DeliveryResponse {
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

	response, err := app.Test(request, int(time.Second) * 5)
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

	return &responseBody.Data
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

		response, err := app.Test(request, int(time.Second) * 5)
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

func DoCreateManyAddress(t *testing.T, token string, totalData int, returnDataByIndex int, delivery *model.DeliveryResponse) *model.AddressResponse {
	var addresses *model.AddressResponse
	for i := 1; i <= totalData; i++ {
		requestBody := model.AddressCreateRequest{
			DeliveryId:      delivery.ID,
			CompleteAddress: fmt.Sprintf("Complete Address %+v", i),
			GoogleMapsLink:  "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6",
			IsMain:          true,
		}

		bodyJson, err := json.Marshal(requestBody)
		assert.Nil(t, err)
		request := httptest.NewRequest(http.MethodPost, "/api/users/current/addresses", strings.NewReader(string(bodyJson)))
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Accept", "application/json")
		request.Header.Set("Authorization", token)

		response, err := app.Test(request, int(time.Second) * 5)
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

		if i == returnDataByIndex {
			addresses = &responseBody.Data
		}
	}
	return addresses
}

func DoCreateManyDiscountCoupon(t *testing.T, token string, totalData int, returnDataByIndex int) *model.DiscountCouponResponse {
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

		response, err := app.Test(request, int(time.Second) * 5)
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

		if i == returnDataByIndex {
			getDiscountCoupon = &responseBody.Data
		}
	}

	return getDiscountCoupon
}

func DoCreateDiscountCouponCustom(t *testing.T, token string, name string, desc string, code string, tipe helper.DiscountType, value float32, start helper.TimeRFC3339, end helper.TimeRFC3339, totalMaxUsage int, maxUsagePerUser int, minOrderValue float32, status bool) *model.DiscountCouponResponse {
	requestBody := model.CreateDiscountCouponRequest{
		Name:            name,
		Description:     desc,
		Code:            code,
		Value:           value,
		Type:            tipe,
		Start:           start,
		End:             end,
		TotalMaxUsage:   totalMaxUsage,
		MaxUsagePerUser: maxUsagePerUser,
		UsedCount:       0,
		MinOrderValue:   minOrderValue,
		Status:          status,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/discount-coupons", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second) * 5)
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

	return &responseBody.Data
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

	response, err := app.Test(request, int(time.Second) * 5)
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

func DoCreateProduct(t *testing.T, token string, totalData int, getProductByIndex int) *model.ProductResponse {
	createCategory := DoCreateCategory(t, token, "Makanan", "Ini adalah makanan")
	var getProduct *model.ProductResponse
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

		response, err := app.Test(request, int(time.Second) * 5)
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
			getProduct = &responseBody.Data
		}
	}

	return getProduct
}

func DoSetBalanceManually(token string, balance_value float32) {
	userEntity := new(entity.User)
	db.Model(entity.User{}).Joins("left join tokens on tokens.user_id = users.id").Where("tokens.token = ?", token).Scan(&userEntity)
	// update balance
	db.Model(entity.Wallet{}).Where("user_id = ?", userEntity.ID).Update("balance", balance_value)
}

func GetCurrentUserByToken(t *testing.T, token string) *model.UserResponse {
	request := httptest.NewRequest(http.MethodGet, "/api/users/current", nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", token)

	response, err := app.Test(request, int(time.Second) * 5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[*model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)
	return responseBody.Data
}

func getRFC3339WithOffsetAndTime(days, weeks, months, hour, minute, second int) string {
	loc, _ := time.LoadLocation("Asia/Jakarta") // GMT+7

	// Ambil sekarang lalu tambah offset
	base := time.Now().In(loc)
	offset := base.AddDate(0, months, days+(weeks*7))

	// Set waktu (jam, menit, detik) ke yang diinginkan
	final := time.Date(
		offset.Year(), offset.Month(), offset.Day(),
		hour, minute, second, 0, loc,
	)

	return final.Format(time.RFC3339)
}

func DoCreateManyOrderUsingWalletPayment(t *testing.T, token string, totalOrder int, discountCoupon *model.DiscountCouponResponse, product *model.ProductResponse, delivery *model.DeliveryResponse) {
	GetCurrentUserByToken(t, token)
	DoSetBalanceManually(token, float32(150000*totalOrder))
	getDelivery := DoCreateManyAddress(t, token, 2, 1, delivery)

	for i := 1; i <= totalOrder; i++ {
		requestBody := model.CreateOrderRequest{
			DiscountId:     discountCoupon.ID,
			PaymentGateway: helper.PAYMENT_GATEWAY_SYSTEM,
			PaymentMethod:  helper.PAYMENT_METHOD_WALLET,
			ChannelCode:    helper.WALLET_CHANNEL_CODE,
			IsDelivery:     true,
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
		assert.Equal(t, discountCoupon.Value, responseBody.Data.DiscountValue)
		assert.Equal(t, helper.PAYMENT_GATEWAY_SYSTEM, responseBody.Data.PaymentGateway)
		assert.Equal(t, helper.PAYMENT_METHOD_WALLET, responseBody.Data.PaymentMethod)
		assert.Equal(t, helper.PAID_PAYMENT, responseBody.Data.PaymentStatus)
		assert.Equal(t, helper.WALLET_CHANNEL_CODE, responseBody.Data.ChannelCode)
		assert.Equal(t, helper.ORDER_PENDING, responseBody.Data.OrderStatus)
		assert.Equal(t, true, responseBody.Data.IsDelivery)
		assert.Equal(t, float32(getDelivery.Delivery.Cost), responseBody.Data.DeliveryCost)
		assert.Equal(t, "Yang cepet ya!", responseBody.Data.Note)
		var totalProductPrice float32 = product.Price * 1

		assert.Equal(t, totalProductPrice, responseBody.Data.TotalProductPrice)
		assert.Equal(t, totalProductPrice+getDelivery.Delivery.Cost-responseBody.Data.TotalDiscount, responseBody.Data.TotalFinalPrice)
		assert.Equal(t, len(requestBody.OrderProducts), len(responseBody.Data.OrderProducts))
		for i, product := range responseBody.Data.OrderProducts {
			assert.Equal(t, requestBody.OrderProducts[i].ProductId, product.ProductId)
			assert.Equal(t, requestBody.OrderProducts[i].Quantity, product.Quantity)
		}

		assert.Nil(t, responseBody.Data.XenditTransaction)
	}
}

func DoRegisterAdmin(t *testing.T) {
	requestBody := model.RegisterUserRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@email.com",
		Phone:     "08123456789",
		Password:  "johndoe123",
		Role:      helper.ADMIN,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, int(time.Second) * 5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, requestBody.LastName, responseBody.Data.LastName)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.Equal(t, requestBody.Role, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func DoRegisterCustomer(t *testing.T) {
	requestBody := model.RegisterUserRequest{
		FirstName: "Customer",
		LastName:  "1",
		Email:     "customer1@email.com",
		Phone:     "0982131244",
		Password:  "customer1",
		Role:      helper.CUSTOMER,
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)
	request := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := app.Test(request, int(time.Second) * 5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, requestBody.LastName, responseBody.Data.LastName)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.Equal(t, requestBody.Role, responseBody.Data.Role)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func DoCreateOrderAsCustomerWithDeliveryAndDiscount(t *testing.T, tokenAdmin string, tokenCust string) *model.OrderResponse  {
	start := getRFC3339WithOffsetAndTime(0, 0, 0, 0, 1, 0)
	parseStart, err := time.Parse(time.RFC3339, start)
	assert.Nil(t, err)

	end := getRFC3339WithOffsetAndTime(15, 0, 0, 23, 59, 59)
	parseEnd, err := time.Parse(time.RFC3339, end)
	assert.Nil(t, err)
	getDiscountCoupon := DoCreateDiscountCouponCustom(t, tokenAdmin, "Lima-Promo", "Ini discount 5%", "#ABC5", helper.PERCENT, float32(5), helper.TimeRFC3339(parseStart), helper.TimeRFC3339(parseEnd), 100, 3, 50000, true)

	delivery := DoCreateDelivery(t, tokenAdmin)
	currentUser := GetCurrentUserByToken(t, tokenCust)
	DoSetBalanceManually(tokenCust, float32(150000))
	getDelivery := DoCreateManyAddress(t, tokenCust, 2, 1, delivery)
	product := DoCreateProduct(t, tokenAdmin, 2, 1)
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
	request.Header.Set("Authorization", tokenCust)

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
	currentUser = GetCurrentUserByToken(t, tokenCust)
	assert.Equal(t, helper.RoundFloat32((float32(150000)-responseBody.Data.TotalFinalPrice), 1), currentUser.Wallet.Balance)

	assert.Nil(t, responseBody.Data.XenditTransaction)

	return &responseBody.Data
}

func DoCreateApplicationSetting(t *testing.T, tokenAdmin string) {
	// Simulasi multipart body
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Tambahkan field JSON sebagai string field biasa
	_ = writer.WriteField("app_name", "Warung Seblak Jaman Now")
	_ = writer.WriteField("opening_hours", "07:00:00")
	_ = writer.WriteField("closing_hours", "20:00:00")
	_ = writer.WriteField("address", "Ini adalah alamat")
	_ = writer.WriteField("google_maps_link", "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6")
	_ = writer.WriteField("description", "Ini adalah deskripsi")
	_ = writer.WriteField("phone_number", "08133546789")
	_ = writer.WriteField("email", "seblak@mail.com")
	_ = writer.WriteField("instagram_name", "fauzan.hidayat-instagram")
	_ = writer.WriteField("instagram_link", "https://www.instagram.com/")
	_ = writer.WriteField("twitter_name", "fauzan.hidayat-twitter")
	_ = writer.WriteField("twitter_link", "https://www.twitter.com/")
	_ = writer.WriteField("facebook_name", "fauzan.hidayat-facebook")
	_ = writer.WriteField("facebook_link", "https://www.facebook.com/")

	filename, content, err := GenerateDummyJPEG(1 * 1024 * 1024) // 1 MB
	assert.Nil(t, err)

	partHeader := textproto.MIMEHeader{}
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="logo_filename"; filename="%s"`, filename))
	partHeader.Set("Content-Type", "image/jpeg")

	fileWriter, err := writer.CreatePart(partHeader)
	assert.Nil(t, err)

	_, err = fileWriter.Write(content)
	assert.Nil(t, err)

	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/api/applications", &b)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", tokenAdmin)

	response, err := app.Test(request, int(time.Second) * 5)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)

	responseBody := new(model.ApiResponse[model.ApplicationResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Data.ID)
	assert.Equal(t, "Warung Seblak Jaman Now", responseBody.Data.AppName)
	assert.Equal(t, "07:00:00", responseBody.Data.OpeningHours)
	assert.Equal(t, "20:00:00", responseBody.Data.ClosingHours)
	assert.Equal(t, "Ini adalah alamat", responseBody.Data.Address)
	assert.Equal(t, "https://maps.app.goo.gl/ftF7eEsBHa69uw3H6", responseBody.Data.GoogleMapsLink)
	assert.Equal(t, "Ini adalah deskripsi", responseBody.Data.Description)
	assert.Equal(t, "08133546789", responseBody.Data.PhoneNumber)
	assert.Equal(t, "seblak@mail.com", responseBody.Data.Email)
	assert.Equal(t, "fauzan.hidayat-instagram", responseBody.Data.InstagramName)
	assert.Equal(t, "https://www.instagram.com/", responseBody.Data.InstagramLink)
	assert.Equal(t, "fauzan.hidayat-twitter", responseBody.Data.TwitterName)
	assert.Equal(t, "https://www.twitter.com/", responseBody.Data.TwitterLink)
	assert.Equal(t, "fauzan.hidayat-facebook", responseBody.Data.FacebookName)
	assert.Equal(t, "https://www.facebook.com/", responseBody.Data.FacebookLink)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}