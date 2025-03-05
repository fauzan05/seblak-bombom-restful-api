package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func CategoryToResponse(category *entity.Category) *model.CategoryResponse {
	return &model.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.Created_At,
		UpdatedAt:   category.Updated_At,
	}
}

func CategoriesToResponse(categories *[]entity.Category) *[]model.CategoryResponse {
	getCategories := make([]model.CategoryResponse, len(*categories))
	for i, category := range *categories {
		getCategories[i] = *CategoryToResponse(&category)
	}
	return &getCategories
}
