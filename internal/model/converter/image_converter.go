package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

// ubah entitiy image menjadi response image
func ImageToResponse(image *entity.Image) *model.ImageResponse {
	return &model.ImageResponse{
		ID:        image.ID,
		ProductId: image.ProductId,
		FileName:  image.FileName,
		Type:      image.Type,
		Position:  image.Position,
		CreatedAt: helper_others.TimeRFC3339(image.CreatedAt),
		UpdatedAt: helper_others.TimeRFC3339(image.UpdatedAt),
	}
}

func ImagesToResponse(images *[]entity.Image) *[]model.ImageResponse {
	getImages := make([]model.ImageResponse, len(*images))
	for i, image := range *images {
		getImages[i] = *ImageToResponse(&image)
	}
	return &getImages
}
