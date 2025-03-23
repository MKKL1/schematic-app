package post

type PostCreated struct {
	Id          int64                    `json:"id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"description"`
	Owner       int64                    `json:"owner"`
	AuthorId    *int64                   `json:"author_id"`
	Categories  PostCategoriesStructured `json:"categories"`
	Tags        []string                 `json:"tags"`
	Files       []PostCreatedFileData    `json:"files"`
}

type PostCreatedFileData struct {
	TempId string `json:"tempId"`
}
