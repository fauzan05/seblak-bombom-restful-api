package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func ApplicationToResponse(application *entity.Application) *model.ApplicationResponse {
	return &model.ApplicationResponse{
		ID: application.ID,
		AppName: application.AppName,
		OpeningHours: application.OpeningHours,
		ClosingHours: application.ClosingHours,
		Address: application.Address,
		Longitude: application.Longitude,
		Latitude: application.Latitude,
		GoogleMapLink: application.GoogleMapLink,
		Description: application.Description,
		PhoneNumber: application.PhoneNumber,
		Email: application.Email,
		InstagramName: application.SocialMedia.InstagramName,
		InstagramLink: application.SocialMedia.InstagramLink,
		TwitterName: application.SocialMedia.TwitterName,
		TwitterLink: application.SocialMedia.TwitterLink,
		FacebookName: application.SocialMedia.FacebookName,
		FacebookLink: application.SocialMedia.FacebookLink,
		CreatedAt: application.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: application.Updated_At.Format("2006-01-02 15:04:05"),
	}
	
}
