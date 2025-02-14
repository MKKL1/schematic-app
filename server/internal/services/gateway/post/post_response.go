package post

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"strconv"
)

type PostResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"desc"`
	Owner       string                   `json:"owner"`
	AuthorID    *string                  `json:"author"`
	Categories  []PostCategoriesResponse `json:"categories"`
	Tags        []string                 `json:"tags"`
}

type PostCategoriesResponse struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

func PostToResponse(post client.Post) PostResponse {
	var authorID *string
	if post.AuthorID != nil {
		aInt := strconv.FormatInt(*post.AuthorID, 10)
		authorID = &aInt
	}

	tags := make([]string, len(post.Tags))
	for i, v := range post.Tags {
		tags[i] = v
	}

	categs := make([]PostCategoriesResponse, len(post.Categories))
	for i, v := range post.Categories {
		categs[i] = PostCategoriesResponse{
			Name:     v.Name,
			Metadata: v.Metadata,
		}
	}

	return PostResponse{
		ID:          strconv.FormatInt(post.ID, 10),
		Name:        post.Name,
		Description: post.Description,
		Owner:       strconv.FormatInt(post.Owner, 10),
		AuthorID:    authorID,
		Categories:  categs,
		Tags:        tags,
	}
}
