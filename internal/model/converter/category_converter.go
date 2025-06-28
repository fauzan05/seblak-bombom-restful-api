package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func CategoryToResponse(category *entity.Category) *model.CategoryResponse {
	response := &model.CategoryResponse{
		ID:            category.ID,
		Name:          category.Name,
		Description:   category.Description,
		ImageFilename: category.ImageFilename,
		IsActive:      true,
		CreatedAt:     helper_others.TimeRFC3339(category.CreatedAt),
		UpdatedAt:     helper_others.TimeRFC3339(category.UpdatedAt),
	}

	if category.DeletedAt.Valid {
		response.IsActive = false
	}

	return response
}

func CategoriesToResponse(categories *[]entity.Category) *[]model.CategoryResponse {
	getCategories := make([]model.CategoryResponse, len(*categories))
	for i, category := range *categories {
		getCategories[i] = *CategoryToResponse(&category)
	}
	return &getCategories
}
