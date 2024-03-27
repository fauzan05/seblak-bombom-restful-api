package database

import (
	"fmt"
	"seblak-bombom-restful-api/internal/config"
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"testing"

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

