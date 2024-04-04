package converter

import (
	"seblak-bombom-restful-api/internal/entity"
	"seblak-bombom-restful-api/internal/model"
)

func ProductReviewToResponse(productReview *entity.ProductReview) *model.ProductReviewResponse {
	return &model.ProductReviewResponse{
		ID:        productReview.ID,
		ProductId: productReview.ProductId,
		UserId:    productReview.UserId,
		Rate:      productReview.Rate,
		Comment:   productReview.Comment,
		CreatedAt: productReview.Created_At.Format("2006-01-02 15:04:05"),
		UpdatedAt: productReview.Updated_At.Format("2006-01-02 15:04:05"),
	}
}

func ProductReviewsToResponse(productReviews *[]entity.ProductReview) *[]model.ProductReviewResponse {
	getProductReviews := make([]model.ProductReviewResponse, len(*productReviews))
	for i, productReview := range *productReviews {
		getProductReviews[i] = *ProductReviewToResponse(&productReview)
	}
	return &getProductReviews
}
