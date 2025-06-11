package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/helper/helper_others"
	"seblak-bombom-restful-api/internal/model"
)

func ProductToResponse(product *entity.Product) *model.ProductResponse {
	if product == nil {
		return &model.ProductResponse{}
	}
	response := &model.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   helper_others.TimeRFC3339(product.CreatedAt),
		UpdatedAt:   helper_others.TimeRFC3339(product.UpdatedAt),
	}

	if product.Category != nil {
		response.Category = *CategoryToResponse(product.Category)
	}

	if product.Images != nil {
		response.Images = *ImagesToResponse(&product.Images)
	}

	if product.Reviews != nil {
		response.Reviews = *ProductReviewsToResponse(&product.Reviews)
	}

	return response
}

func ProductsToResponse(products *[]entity.Product) *[]model.ProductResponse {
	getProducts := make([]model.ProductResponse, len(*products))
	for i, product := range *products {
		getProducts[i] = *ProductToResponse(&product)
	}
	return &getProducts
}
