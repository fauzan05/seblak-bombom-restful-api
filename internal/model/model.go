package model

type ApiResponse[T any] struct {
	Code   int32  `json:"code"`
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type ApiResponsePagination[T any] struct {
	Code               int32  `json:"code"`
	Status             string `json:"status"`
	Data               T      `json:"data"`
	TotalCurrentDatas  int64  `json:"total_current_datas"`
	TotalActiveDatas   int64  `json:"total_active_datas"`
	TotalInactiveDatas int64  `json:"total_inactive_datas"`
	TotalRealDatas     int64  `json:"total_real_datas"`
	TotalPages         int    `json:"total_pages"`
	CurrentPages       int    `json:"current_pages"`
	DataPerPages       int    `json:"data_per_pages"`
}

type ErrorResponse[message string] struct {
	Error message `json:"errors"`
}
