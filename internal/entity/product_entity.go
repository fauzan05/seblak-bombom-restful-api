package entity

import "time"

// product is a struct that represents a product entity in database table
type Product struct {
	ID          uint64    `gorm:"primary_key;column:id;autoIncrement"`
	CategoryId  uint64    `gorm:"column:category_id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Price       int       `gorm:"column:price"`
	Stock       int       `gorm:"column:stock"`
	Created_At  time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Category    *Category `gorm:"foreignKey:category_id;references:id"`
}

func (p *Product) TableName() string {
	return "products"
}
