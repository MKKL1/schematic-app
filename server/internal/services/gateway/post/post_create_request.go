package post

type PostCreateRequest struct {
	Name        string                      `json:"name" validate:"required,min=5,max=32"`
	Description *string                     `json:"desc" validate:"omitempty,max=100"`
	Author      *int64                      `json:"author" validate:"omitempty"`
	Categories  []PostCategoryCreateRequest `json:"categories" validate:"required"`
	Tags        []string                    `json:"tags" validate:"omitempty"`
}

type PostCategoryCreateRequest struct {
	Name     string                 `json:"name" validate:"required,min=2,max=32"`
	Metadata map[string]interface{} `json:"metadata" validate:"omitempty"`
}
