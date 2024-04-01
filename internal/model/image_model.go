package model

type ImageResponse struct {
	ID        uint64 `json:"id,omitempty"`
	ProductId uint64 `json:"product_id,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	Type      string `json:"type,omitempty"`
	Position  int    `json:"position,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

type AddImagesRequest struct {
	Images []ImageRequest `json:"-" validate:"required"`
}

type ImageRequest struct {
	ProductId uint64 `json:"product_id" validate:"required"`
	FileName  string `json:"file_name" validate:"required,max=100"`
	Type      string `json:"type" validate:"required"`
	Position  int    `json:"position" validate:"required"`
}
