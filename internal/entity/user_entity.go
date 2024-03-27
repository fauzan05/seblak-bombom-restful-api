package entity

import "time"

// user is a struct that represents a user entity in database table
type User struct {
	ID         uint64    `gorm:"column:id;primaryKey"`
	Name       Name      `gorm:"embedded"`
	Email      string    `gorm:"column:email"`
	Phone      string    `gorm:"column:phone"`
	Password   string    `gorm:"column:password"`
	Created_At time.Time `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}

type Name struct {
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
}
