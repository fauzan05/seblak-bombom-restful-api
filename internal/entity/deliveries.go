package entity

import "time"

type Delivery struct {
	ID         uint64    `gorm:"primary_key;column:id;autoIncrement"`
	Cost       int       `gorm:"column:cost"`
	Distance   int       `gorm:"column:distance"`
	Created_At time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Delivery) TableName() string {
	return "deliveries"
}
