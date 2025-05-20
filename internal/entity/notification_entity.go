package entity

import (
	"seblak-bombom-restful-api/internal/helper"
	"time"
)

type Notification struct {
	ID          uint64                  `gorm:"primary_key;column:id;autoIncrement"`
	UserID      uint64                  `gorm:"column:user_id"`
	Title       string                  `gorm:"column:title"`
	IsRead      bool                    `gorm:"column:is_read"`
	Type        helper.NotificationType `gorm:"column:type"`
	BodyContent string                  `gorm:"column:body_content"`
	CreatedAt   time.Time               `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt   time.Time               `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (c *Notification) TableName() string {
	return "notifications"
}
