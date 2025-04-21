package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper"
	"seblak-bombom-restful-api/internal/model"
)

func ApplicationToResponse(application *entity.Application) *model.ApplicationResponse {
	return &model.ApplicationResponse{
		ID:             application.ID,
		AppName:        application.AppName,
		LogoFilename:   application.LogoFilename,
		OpeningHours:   application.OpeningHours,
		ClosingHours:   application.ClosingHours,
		Address:        application.Address,
		GoogleMapsLink: application.GoogleMapsLink,
		Description:    application.Description,
		PhoneNumber:    application.PhoneNumber,
		Email:          application.Email,
		InstagramName:  application.SocialMedia.InstagramName,
		InstagramLink:  application.SocialMedia.InstagramLink,
		TwitterName:    application.SocialMedia.TwitterName,
		TwitterLink:    application.SocialMedia.TwitterLink,
		FacebookName:   application.SocialMedia.FacebookName,
		FacebookLink:   application.SocialMedia.FacebookLink,
		CreatedAt:      helper.TimeRFC3339(application.CreatedAt),
		UpdatedAt:      helper.TimeRFC3339(application.UpdatedAt),
	}

}
