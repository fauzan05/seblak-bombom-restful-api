package model

type CategoryResponse struct {
	ID          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
}

type GetCategoryRequest struct {
	ID uint64 `json:"-" validate:"required"`
}

type UpdateCategoryRequest struct {
	ID          uint64 `json:"-" validate:"required"`
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description"`
}

type DeleteCategoryRequest struct {
	ID uint64 `json:"-" validate:"required"`
}
