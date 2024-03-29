package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func AddressToResponse(address *entity.Address) *model.AddressResponse {
	return &model.AddressResponse{
		ID: address.ID,
		Regency: address.Regency,
		Subdistrict: address.SubDistrict,
		CompleteAddress: address.CompleteAddress,
		GoogleMapLink: address.GoogleMapLink,
		IsMain: address.IsMain,
		CreatedAt: address.Created_At,
		UpdatedAt: address.Updated_At,
	}
}