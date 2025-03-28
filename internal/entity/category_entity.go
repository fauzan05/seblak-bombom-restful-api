package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint64         `gorm:"primary_key;column:id;autoIncrement"`
	Name        string         `gorm:"column:name"`
	Description string         `gorm:"column:description"`
	Created_At  time.Time      `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At  time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	Deleted_At  gorm.DeletedAt `gorm:"column:deleted_at"`
	Products    []Product      `gorm:"foreignKey:category_id;references:id"`
}

func (c *Category) TableName() string {
	return "categories"
}
