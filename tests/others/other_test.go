package others

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type Waktu struct {
	CreatedAt string `json:"created_at"`
}

func TestIncludeStringIntoTime(t *testing.T) {
	data := `{
		"created_at": "2006-01-02 15:04:05"
	}`

	// Konversi data JSON string menjadi slice byte ([]byte)
	jsondata := []byte(data)

	// Unmarshal data JSON ke struktur Waktu
	var waktu Waktu
	err := json.Unmarshal(jsondata, &waktu)
	if err != nil {
		panic(err)
	}
	layout := "2006-01-02 15:04:05"
	createdAt, _ := time.Parse(layout, waktu.CreatedAt)
	fmt.Printf("%T", createdAt)
	// Tampilkan hasil unmarshal
	fmt.Println("Created At:", waktu.CreatedAt)
}

func TestCalculateAmount(t *testing.T) {
	total := float32(65000)
	discount := float32(10) / float32(100)
	afterDiscount := total * discount
	result := total - afterDiscount
	fmt.Println(result)
}
