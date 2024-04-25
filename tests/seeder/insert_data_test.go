package seeder

import (
	"seblak-bombom-restful-api/internal/config"
	// "seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"testing"

	"github.com/spf13/viper"
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

// Insert data user admin dan 1 customer user

func TestInsertUser(t *testing.T) {
	// adminUser := entity.User{
	// 	Name: entity.Name{
	// 		FirstName: "Fauzan",
	// 		LastName: "Nurhidayat",
	// 	},
	// 	Email: "fauzannurhidayat8@gmail.com",
	// 	Phone: "081335457601",
	// 	Password: ,
	// }

}