package entity

import "time"

type Delivery struct {
	ID         uint64    `gorm:"primary_key;column:id;autoIncrement"`
	City       string    `gorm:"column:city"`
	District   string    `gorm:"column:district"`
	Village    string    `gorm:"column:village"`
	Hamlet     string    `gorm:"column:hamlet"`
	Cost       float32   `gorm:"column:cost"`
	Created_At time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Delivery) TableName() string {
	return "deliveries"
}
