package entity

import "time"

type PasswordReset struct {
	ID               uint64    `gorm:"primary_key;column:id;autoIncrement"`
	UserId           uint64    `gorm:"column:user_id"`
	VerificationCode int       `gorm:"column:verification_code"`
	ExpiresAt        time.Time `gorm:"column:expires_at"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	User             *User     `gorm:"foreignKey:user_id;references:id"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
