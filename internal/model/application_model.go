package model

type ApplicationResponse struct {
	ID            uint64  `json:"id,omitempty"`
	AppName       string  `json:"app_name,omitempty"`
	OpeningHours  string  `json:"opening_hours,omitempty"`
	ClosingHours  string  `json:"closing_hours,omitempty"`
	Address       string  `json:"address,omitempty"`
	Longitude     float64 `json:"longitude,omitempty"`
	Latitude      float64 `json:"latitude,omitempty"`
	Description   string  `json:"description,omitempty"`
	PhoneNumber   string  `json:"phone_number,omitempty"`
	Email         string  `json:"email,omitempty"`
	InstagramName string  `json:"instagram_name,omitempty"`
	InstagramLink string  `json:"instagram_link,omitempty"`
	TwitterName   string  `json:"twitter_name,omitempty"`
	TwitterLink   string  `json:"twitter_link,omitempty"`
	FacebookName  string  `json:"facebook_name,omitempty"`
	FacebookLink  string  `json:"facebook_link,omitempty"`
	CreatedAt     string  `json:"created_at,omitempty"`
	UpdatedAt     string  `json:"updated_at,omitempty"`
}

type CreateApplicationRequest struct {
	AppName       string  `json:"app_name" validate:"required,max=100"`
	OpeningHours  string  `json:"opening_hours"`
	ClosingHours  string  `json:"closing_hours"`
	Address       string  `json:"address"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	Description   string  `json:"description"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	InstagramName string  `json:"instagram_name"`
	InstagramLink string  `json:"instagram_link"`
	TwitterName   string  `json:"twitter_name"`
	TwitterLink   string  `json:"twitter_link"`
	FacebookName  string  `json:"facebook_name"`
	FacebookLink  string  `json:"facebook_link"`
}

type UpdateApplicationRequest struct {
	AppName       string  `json:"app_name" validate:"required,max=100"`
	OpeningHours  string  `json:"opening_hours"`
	ClosingHours  string  `json:"closing_hours"`
	Address       string  `json:"address"`
	Longitude     float64 `json:"longitude"`
	Latitude      float64 `json:"latitude"`
	Description   string  `json:"description"`
	PhoneNumber   string  `json:"phone_number"`
	Email         string  `json:"email"`
	InstagramName string  `json:"instagram_name"`
	InstagramLink string  `json:"instagram_link"`
	TwitterName   string  `json:"twitter_name"`
	TwitterLink   string  `json:"twitter_link"`
	FacebookName  string  `json:"facebook_name"`
	FacebookLink  string  `json:"facebook_link"`
}
