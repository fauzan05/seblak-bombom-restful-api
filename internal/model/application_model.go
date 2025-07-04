package model

import (
	"mime/multipart"
	"seblak-bombom-restful-api/internal/helper/helper_others"
)

type ApplicationResponse struct {
	ID             uint64             `json:"id"`
	AppName        string             `json:"app_name"`
	LogoFilename   string             `json:"logo_filename"`
	OpeningHours   string             `json:"opening_hours"`
	ClosingHours   string             `json:"closing_hours"`
	Address        string             `json:"address"`
	GoogleMapsLink string             `json:"google_maps_link"`
	Description    string             `json:"description"`
	PhoneNumber    string             `json:"phone_number"`
	Email          string             `json:"email"`
	ServiceFee     float32            `json:"service_fee"`
	InstagramName  string             `json:"instagram_name"`
	InstagramLink  string             `json:"instagram_link"`
	TwitterName    string             `json:"twitter_name"`
	TwitterLink    string             `json:"twitter_link"`
	FacebookName   string             `json:"facebook_name"`
	FacebookLink   string             `json:"facebook_link"`
	CreatedAt      helper_others.TimeRFC3339 `json:"created_at"`
	UpdatedAt      helper_others.TimeRFC3339 `json:"updated_at"`
}

type CreateApplicationRequest struct {
	ID             uint64                `json:"id"`
	AppName        string                `json:"app_name" validate:"required,max=100"`
	Logo           *multipart.FileHeader `json:"logo_filename"`
	OpeningHours   string                `json:"opening_hours"`
	ClosingHours   string                `json:"closing_hours"`
	Address        string                `json:"address"`
	GoogleMapsLink string                `json:"google_maps_link"`
	Description    string                `json:"description"`
	PhoneNumber    string                `json:"phone_number"`
	Email          string                `json:"email"`
	InstagramName  string                `json:"instagram_name"`
	InstagramLink  string                `json:"instagram_link"`
	TwitterName    string                `json:"twitter_name"`
	TwitterLink    string                `json:"twitter_link"`
	FacebookName   string                `json:"facebook_name"`
	FacebookLink   string                `json:"facebook_link"`
	ServiceFee     float32               `json:"service_fee"`
}
