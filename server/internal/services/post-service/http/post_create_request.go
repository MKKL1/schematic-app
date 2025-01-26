package http

type PostCreateRequest struct {
	Name        string          `json:"name" validate:"required,min=5,max=32"`
	Description *string         `json:"desc" validate:"omitempty,max=100"`
	Author      *AuthorResponse `json:"author" validate:"omitempty"`
}
