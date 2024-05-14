package entity

import "time"

// token is a struct that represents a token entity in database table
type Application struct {
	ID           uint64      `gorm:"primary_key;column:id;autoIncrement"`
	AppName      string      `gorm:"column:app_name"`
	OpeningHours string      `gorm:"column:opening_hours"`
	ClosingHours string      `gorm:"column:closing_hours"`
	Address      string      `gorm:"column:address"`
	Longitude    float64     `gorm:"column:longitude"`
	Latitude     float64     `gorm:"column:latitude"`
	Description  string      `gorm:"column:description"`
	PhoneNumber  string      `gorm:"column:phone_number"`
	Email        string      `gorm:"column:email"`
	SocialMedia  SocialMedia `gorm:"embedded"`
	Created_At   time.Time   `gorm:"column:created_at;autoCreateTime;<-:create"`
	Updated_At   time.Time   `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
}

func (u *Application) TableName() string {
	return "applications"
}

type SocialMedia struct {
	InstagramName string `gorm:"column:instagram_name"`
	InstagramLink string `gorm:"column:instagram_link"`
	TwitterName   string `gorm:"column:twitter_name"`
	TwitterLink   string `gorm:"column:twitter_link"`
	FacebookName  string `gorm:"column:facebook_name"`
	FacebookLink  string `gorm:"column:facebook_link"`
}
