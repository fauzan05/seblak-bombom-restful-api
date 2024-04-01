package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func ProductToResponse(product *entity.Product) *model.ProductResponse {
	return &model.ProductResponse{
		ID:          product.ID,
		Category:    *CategoryToResponse(product.Category),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		CreatedAt:   product.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt:   product.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func ProductsToResponse(products *[]entity.Product) *[]model.ProductResponse {
	getProducts := make([]model.ProductResponse, len(*products))
	for i, product := range *products {
		getProducts[i] = *ProductToResponse(&product)
	}
	return &getProducts
}
