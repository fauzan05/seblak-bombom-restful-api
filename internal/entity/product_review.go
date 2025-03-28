package entity

import "time"

type ProductReview struct {
	ID         uint64    `gorm:"primary_key;column:id;autoIncrement"`
	ProductId  uint64    `gorm:"column:product_id"`
	UserId     uint64    `gorm:"column:user_id"`
	Rate       int       `gorm:"column:rate"`
	Comment    string    `gorm:"column:comment"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Product    *Product  `gorm:"foreignKey:product_id;references:id"`
	User       *User     `gorm:"foreignKey:user_id;references:id"`
}

func (p *ProductReview) TableName() string {
	return "product_reviews"
}
