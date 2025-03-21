package entity

import "time"

// token is a struct that represents a token entity in database table
type Token struct {
	ID         uint64    `gorm:"primary_key;column:id;autoIncrement"`
	Token      string    `gorm:"column:token"`
	UserId     uint64    `gorm:"column:user_id"`
	ExpiryDate time.Time `gorm:"column:expiry_date"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	User       *User     `gorm:"foreignKey:user_id;references:id"`
}

func (u *Token) TableName() string {
	return "tokens"
}
