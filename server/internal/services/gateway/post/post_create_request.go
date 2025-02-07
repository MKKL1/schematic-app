package post

type PostCreateRequest struct {
	Name        string               `json:"name" validate:"required,min=5,max=32"`
	Description *string              `json:"desc" validate:"omitempty,max=100"`
	Author      *AuthorCreateRequest `json:"author" validate:"omitempty"`
}

type AuthorCreateRequest struct {
	Name *string `json:"name" validate:"omitempty,alphanum"`
	ID   *string `json:"id" validate:"omitempty,number"`
}
