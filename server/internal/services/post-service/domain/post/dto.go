package post

type Post struct {
	ID          int64
	Description *string
	Owner       int64
	Author      *Author
}

type Author struct {
	IsKnown bool
	Name    string
	UserID  int64
}

func ToDTO(postModel Model) Post {
	var author *Author
	if postModel.AuthorID == nil && postModel.AuthorName == nil {
		author = nil
	} else if postModel.AuthorID != nil {
		author = &Author{
			IsKnown: true,
			UserID:  *postModel.AuthorID,
		}
	} else {
		author = &Author{
			IsKnown: false,
			Name:    *postModel.AuthorName,
		}
	}

	return Post{
		ID:          postModel.ID,
		Description: postModel.Description,
		Owner:       postModel.Owner,
		Author:      author,
	}
}
