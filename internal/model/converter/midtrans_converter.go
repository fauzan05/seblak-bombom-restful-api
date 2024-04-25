package converter

import (
	"seblak-bombom-restful-api/internal/model"

	"github.com/midtrans/midtrans-go/snap"
)

func MidtransToResponse(snapResponse *snap.Response) *model.SnapResponse {
	return &model.SnapResponse{
		Token: snapResponse.Token,
		RedirectUrl: snapResponse.RedirectURL,
	}
}