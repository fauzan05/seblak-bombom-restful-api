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
		Longitude: address.Longitude,
		Latitude: address.Latitude,
		IsMain: address.IsMain,
		CreatedAt: address.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: address.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func AddressesToResponse(addresses *[]model.AddressResponse) *[]model.AddressResponse {
	getAddresses := make([]model.AddressResponse, len(*addresses))
	copy(getAddresses, *addresses)
	return &getAddresses
}
