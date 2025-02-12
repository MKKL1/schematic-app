package post

type Post struct {
	ID           int64
	Name         string
	Description  *string
	Owner        int64
	AuthorID     *int64
	CategoryVars []CategoryVars
	Tags         []string
}

type CategoryVars struct {
	CategoryName string
	Values       CategoryMetadata
}

func ToDTO(postModel Entity) Post {
	return Post{
		ID:          postModel.ID,
		Name:        postModel.Name,
		Description: postModel.Description,
		Owner:       postModel.Owner,
		AuthorID:    postModel.AuthorID,
	}
}
