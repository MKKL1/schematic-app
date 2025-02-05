package post

import (
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"strconv"
)

type PostResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description *string         `json:"desc"`
	Owner       string          `json:"owner"`
	Author      *AuthorResponse `json:"author"`
}

func PostToResponse(post post.Post) PostResponse {
	return PostResponse{
		ID:          strconv.FormatInt(post.ID, 10),
		Name:        post.Name,
		Description: post.Description,
		Owner:       strconv.FormatInt(post.Owner, 10),
		Author:      PostToAuthorResponse(post),
	}
}
