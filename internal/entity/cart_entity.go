package entity

import "time"

type Cart struct {
	ID         uint64    `gorm:"primary_key;column:id;autoIncrement"`
	UserID     uint64    `gorm:"column:user_id"`
	ProductID  uint64    `gorm:"column:product_id"`
	Name       string    `gorm:"column:name"`
	Quantity   int       `gorm:"column:quantity"`
	Price      float32   `gorm:"column:price"`
	TotalPrice float32   `gorm:"column:total_price"`
	Stock      int       `gorm:"column:stock"`
	Created_At time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Cart) TableName() string {
	return "carts"
}
