package model

type ApiResponse[T any] struct {
	Code int8 `json:"code"`
	Status string `json:"status"`
	Data T `json:"data"`
}