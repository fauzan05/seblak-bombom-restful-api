package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func AddressToResponse(address *entity.Address) *model.AddressResponse {
	response := &model.AddressResponse{
		ID:              address.ID,
		CompleteAddress: address.CompleteAddress,
		GoogleMapsLink:  address.GoogleMapsLink,
		IsMain:          address.IsMain,
		CreatedAt:       helper_others.TimeRFC3339(address.CreatedAt),
		UpdatedAt:       helper_others.TimeRFC3339(address.UpdatedAt),
	}

	if address.Delivery != nil {
		response.Delivery = *DeliveryToResponse(address.Delivery)
	}

	return response
}

func AddressesToResponse(addresses *[]model.AddressResponse) *[]model.AddressResponse {
	getAddresses := make([]model.AddressResponse, len(*addresses))
	copy(getAddresses, *addresses)
	return &getAddresses
}
