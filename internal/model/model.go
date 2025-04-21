package model

type ApiResponse[T any] struct {
	Code   int32  `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type ApiResponsePagination[T any] struct {
	Code         int32  `json:"code"`
	Status       string `json:"status"`
	Data         T      `json:"data"`
	TotalDatas   int64  `json:"total_datas"`
	TotalPages   int    `json:"total_pages"`
	CurrentPages int    `json:"current_pages"`
	DataPerPages int    `json:"data_per_pages"`
}

type ErrorResponse[message string] struct {
	Error message `json:"errors"`
}
