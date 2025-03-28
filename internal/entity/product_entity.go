package entity

import (
	"time"

	"gorm.io/gorm"
)

// product is a struct that represents a product entity in database table
type Product struct {
	ID          uint64          `gorm:"primary_key;column:id;autoIncrement"`
	CategoryId  uint64          `gorm:"column:category_id"`
	Name        string          `gorm:"column:name"`
	Description string          `gorm:"column:description"`
	Price       float32         `gorm:"column:price"`
	Stock       int             `gorm:"column:stock"`
	CreatedAt  time.Time       `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt  time.Time       `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt  `gorm:"column:deleted_at"`
	Category    *Category       `gorm:"foreignKey:category_id;references:id"`
	Images      []Image         `gorm:"foreignKey:product_id;references:id"`
	Reviews     []ProductReview `gorm:"foreignKey:product_id;references:id"`
	CartItems   []CartItem      `gorm:"foreignKey:product_id;references:id"`
}

func (p *Product) TableName() string {
	return "products"
}
