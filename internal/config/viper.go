package config

import (
	"os"
	// "path/filepath"
	// "seblak-bombom-restful-api/internal/helper"

	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
    config := viper.New()

    // ⬇️ Ini akan ambil langsung dari env saat dijalankan (termasuk dari Railway)
    config.AutomaticEnv()

    // ⬇️ Jika kamu jalankan lokal dan ada file .env, ini bisa opsional:
    if _, err := os.Stat(".env"); err == nil {
        config.SetConfigFile(".env")
        config.SetConfigType("env")
        _ = config.ReadInConfig() // Boleh diabaikan kalau file tidak ada
    }

    return config
}

