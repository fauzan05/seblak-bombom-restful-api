package database

import (
	"fmt"
	"seblak-bombom-restful-api/internal/config"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
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

var db = config.NewDatabaseTest(NewViper(), config.NewLogger(NewViper()))
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
				GoogleMapLink: "https://maps.app.goo.gl/UBRaYVdBxkUDkHMW7",
				IsMain: true,
			},
		},
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
	result := db.Model(&user).Where("token = ?", token_code).Association("Token").Find(&token)
	fmt.Println(result)
	// fmt.Println(user)
	fmt.Println(token)
}

func TestGetUserWithAddress(t *testing.T) {
	var user []entity.User
	findUser := db.Where("id = ?", 2).Find(&user)
	assert.Nil(t, findUser.Error)
	result := db.Preload("Addresses").Find(&user)
	assert.Nil(t, result.Error)
	fmt.Println(user)
}
