package entity

import "time"

type Cart struct {
	ID         uint64     `gorm:"primary_key;column:id;autoIncrement"`
	UserID     uint64     `gorm:"column:user_id"`
	Created_At time.Time  `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At time.Time  `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User       *User      `gorm:"foreignKey:user_id;references:id"`
	CartItems  []CartItem `gorm:"foreignKey:cart_id;references:id"`
}

func (c *Cart) TableName() string {
	return "carts"
}
