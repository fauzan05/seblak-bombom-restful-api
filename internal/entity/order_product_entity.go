package entity

import (
	"time"
)

type OrderProduct struct {
	ID                        uint64    `gorm:"primary_key;column:id;autoIncrement"`
	OrderId                   uint64    `gorm:"column:order_id"`
	ProductId                 uint64    `gorm:"column:product_id"`
	ProductName               string    `gorm:"column:product_name"`
	ProductFirstImagePosition string    `gorm:"column:product_first_image_position"`
	Category                  string    `gorm:"column:category"`
	Price                     float32   `gorm:"column:price"`
	Quantity                  int       `gorm:"column:quantity"`
	CreatedAt                 time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt                 time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Order                     *Order    `gorm:"foreignKey:order_id;references:id"`
	Product                   *Product  `gorm:"foreignKey:product_id;references:id"`
}

func (u *OrderProduct) TableName() string {
	return "order_products"
}
