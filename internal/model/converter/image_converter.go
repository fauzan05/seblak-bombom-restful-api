package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)
// ubah entitiy image menjadi response image
func ImageToResponse(image *entity.Image) *model.ImageResponse {
	return &model.ImageResponse{
		ID: image.ID,
		ProductId: image.ProductId,
		FileName: image.FileName,
		Type: image.Type,
		Position: image.Position,
		CreatedAt:   image.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt:   image.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func ImagesToResponse(images *[]entity.Image) *[]model.ImageResponse {
	getImages := make([]model.ImageResponse, len(*images))
	for i, image := range *images {
		getImages[i] = *ImageToResponse(&image)
	}
	return &getImages
}