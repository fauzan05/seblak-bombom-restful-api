package entity

import (
	"time"
)

// user is a struct that represents a user entity in database table
type OrderProduct struct {
	ID          uint64    `gorm:"primary_key;column:id;autoIncrement"`
	OrderId     uint64    `gorm:"column:order_id"`
	ProductId   uint64    `gorm:"column:product_id"`
	ProductName string    `gorm:"column:product_name"`
	Category    string    `gorm:"column:category"`
	Price       float32   `gorm:"column:price"`
	Quantity    int       `gorm:"column:quantity"`
	Created_At  time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Order       *Order    `gorm:"foreignKey:order_id;references:id"`
}

func (u *OrderProduct) TableName() string {
	return "order_products"
}
