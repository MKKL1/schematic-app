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
	vars := make([]CategoryVars, len(postModel.CategoryVars))
	for i, v := range postModel.CategoryVars {
		vars[i] = CategoryVars{
			CategoryName: v.CategoryName,
			Values:       v.Metadata,
		}
	}

	tags := make([]string, len(postModel.Tags))
	for i, v := range postModel.Tags {
		tags[i] = v
	}

	return Post{
		ID:           postModel.ID,
		Name:         postModel.Name,
		Description:  postModel.Description,
		Owner:        postModel.Owner,
		AuthorID:     postModel.AuthorID,
		CategoryVars: vars,
		Tags:         tags,
	}
}
