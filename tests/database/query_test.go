package database

import (
	"fmt"
	"seblak-bombom-restful-api/internal/config"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("config")
	config.SetConfigType("json")
	config.AddConfigPath("../../")
	err := config.ReadInConfig()

	helper.HandleErrorWithPanic(err)

	return config
}

var db = config.NewDatabaseProd(NewViper(), config.NewLogger(NewViper()))
func TestCountEmailUser(t *testing.T) {
	var total int64
	err := db.Model(&entity.User{}).Where("email = ?", "fauzannurhidayat8@gmail.com").Count(&total).Error
	assert.Nil(t, err)
	fmt.Println(total)
}

func TestFindUser(t *testing.T) {
	user := entity.User{}
	result := db.First(&user)
	assert.Nil(t, result.Error)
	fmt.Println(user)
}

func TestCreateNewUser(t *testing.T) {
	now := time.Now()
	oneHours := now.Add(1 * time.Hour)

	user := entity.User{
		ID: 1,
		Name: entity.Name{
			FirstName: "Fauzan",
			LastName: "Nurhidayat",
		},
		Email: "fauzannurhidayat8@gmail.com",
		Phone: "081335457601",
		Password: "Fauzan123",
		Token: entity.Token{
			Token: "417313b7-ea06-4e7f-81ed-927ee23ff86c",
			UserId: 1,
			ExpiryDate: oneHours,
		},
		Addresses: []entity.Address{
			{
				UserId: 1,
				Regency: "Kebumen",
				SubDistrict: "Pejagoan",
				CompleteAddress: "Jl tembana-peniron km.12, Dukuh jetis, Desa Peniron RT01/05, Kecamatan Pejagoan, Kabupaten Kebumen, Provinsi Jawa Tengah 54361",
				Longitude: -74.00898606,
				Latitude: 40.71727401,
				IsMain: true,
			},
		},
		Role: helper.ADMIN,
	}
	result := db.Create(&user)
	assert.Nil(t, result.Error)
}

func TestRelation(t *testing.T) {
	user := new(entity.User)
	userResult := db.First(&user)
	assert.Nil(t, userResult.Error)

	// var token entity.Token
	fmt.Println(user)
	// result := db.Where("user_id = ?", 2).First(&token).Error
	// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 	fmt.Println("datanya kosong")
	// } else {
	// 	fmt.Println("datanya ada")
	// }
}

func TestGetUserByToken(t *testing.T) {
	user := new(entity.User)
	token := new(entity.Token)
	token_code := "417313b7-ea06-4e7f-81ed-927ee23ff86c"
	result := db.Where("token = ?", token_code).Joins("User").Find(&token)
	assert.Nil(t, result.Error)
	result = db.Where("id = ?", token.UserId).Preload("Token").Preload("Addresses").Find(&user)
	assert.Nil(t, result.Error)

	fmt.Println(user.Token.Token)
	for _, v := range user.Addresses {
		fmt.Println(v)
	}
}

func TestGetUserWithAddress(t *testing.T) {
	var user []entity.User
	findUser := db.Where("id = ?", 1).Find(&user)
	assert.Nil(t, findUser.Error)
	result := db.Preload("Addresses").Find(&user)
	assert.Nil(t, result.Error)
	fmt.Println(user)
}

func TestCreateNewCategory(t *testing.T) {
	for i := 0; i < 5; i++ {
		newCategory := new(entity.Category)
		// category := &entity.Category{
		// 	Name: "Category " + strconv.Itoa(i),
		// 	Description: "Description " + strconv.Itoa(i),
		// }
		newCategory.Name = "Category " + strconv.Itoa(i)
		newCategory.Description = "Description " + strconv.Itoa(i)
		db.Create(newCategory)
		fmt.Println(newCategory.ID)
	}
}

func TestFindAllCategories(t *testing.T) {
	categories := new([]entity.Category)
	result := db.Find(&categories)
	assert.Nil(t, result.Error)
	for _, v := range *categories {
		fmt.Println(v)
	}
}

func TestFindAndCountById(t *testing.T) {
	newCategory := new(entity.Category)
	newCategory.ID = uint64(5)
	var count int64
	result := db.Find(&newCategory).Count(&count)
	assert.Nil(t, result.Error)
	fmt.Println("Count : ", count)
	fmt.Println(newCategory.ID)
	fmt.Println(newCategory.Name)
	fmt.Println(newCategory.Description)
}

func TestFindOrderWithOrderProducts(t *testing.T) {
	selectedOrder := new(entity.Order)
	selectedOrder.ID = 1
	if err := db.Preload("OrderProducts").Find(selectedOrder).Error; err != nil {
		panic(err)
	}

	type Products struct {
		ID          uint64
		ProductName string
	}
	
	var products []Products
	for _, product := range selectedOrder.OrderProducts {
		newProduct := Products{
			ID:          product.ID,
			ProductName: product.ProductName,
		}
		products = append(products, newProduct)
	}

	fmt.Println(products)
}

func TestFindCurrentCartWithProducts(t *testing.T) {
	newCart := new(entity.Cart)

	if err := db.SetupJoinTable(&entity.Cart{}, "CartItems", &entity.CartItem{}); err != nil {
		panic(err)
	}

	for _, newCart := range newCart.CartItems {
		fmt.Println("DATANYA : ", newCart)
	}
}