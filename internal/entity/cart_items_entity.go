package entity

import "time"

type CartItem struct {
	ID        uint64    `gorm:"primary_key;column:id;autoIncrement"`
	CartId    uint64    `gorm:"column:cart_id"`
	ProductID uint64    `gorm:"column:product_id"`
	Quantity  int       `gorm:"column:quantity"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Cart      *Cart     `gorm:"foreignKey:cart_id;references:id"`
	Product   *Product  `gorm:"foreignKey:product_id;references:id"`
}

func (c *CartItem) TableName() string {
	return "cart_items"
}
