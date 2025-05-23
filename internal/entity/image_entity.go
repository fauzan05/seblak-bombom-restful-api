package entity

import "time"

type Image struct {
	ID         uint64    `gorm:"primary_key;column:id;autoIncrement"`
	ProductId  uint64    `gorm:"column:product_id"`
	FileName   string    `gorm:"column:file_name"`
	Type       string    `gorm:"column:type"`
	Position   int       `gorm:"column:position"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Product    *Product  `gorm:"foreignKey:product_id;references:id"`
}

func (i *Image) TableName() string {
	return "images"
}
