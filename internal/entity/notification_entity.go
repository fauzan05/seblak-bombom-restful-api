package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type Notification struct {
	ID         uint                    `gorm:"primary_key;column:id;autoIncrement"`
	UserID     int                     `gorm:"column:user_id"`
	Title      string                  `gorm:"column:title"`
	Message    string                  `gorm:"column:message"`
	IsRead     bool                    `gorm:"column:is_read"`
	Type       helper.NotificationType `gorm:"column:type"`
	Link       string                  `gorm:"column:link"`
	Created_At time.Time               `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At time.Time               `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Notification) TableName() string {
	return "notifications"
}
