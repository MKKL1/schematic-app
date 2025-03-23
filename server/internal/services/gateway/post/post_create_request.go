package post

type PostCreateRequest struct {
	Name        string                      `json:"name" validate:"required,min=5,max=100"`
	Description *string                     `json:"desc" validate:"omitempty,max=500"`
	Author      *string                     `json:"author_id" validate:"omitempty,number,min=10,max=20"`
	Categories  []PostCategoryCreateRequest `json:"categories" validate:"omitempty,max=5,dive"`
	Tags        []string                    `json:"tags" validate:"omitempty,dive,min=1,max=30"`
	Files       []string                    `json:"files" validate:"omitempty,dive,min=1,max=100"`
}

type PostCategoryCreateRequest struct {
	Name     string                 `json:"name" validate:"required,min=3,max=32"`
	Metadata map[string]interface{} `json:"metadata" validate:"omitempty"`
}
