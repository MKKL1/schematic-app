package post

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"strconv"
)

type PostResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"desc"`
	Owner       string  `json:"owner"`
	AuthorID    *int64  `json:"omitempty,author"`
}

func PostToResponse(post client.Post) PostResponse {
	return PostResponse{
		ID:          strconv.FormatInt(post.ID, 10),
		Name:        post.Name,
		Description: post.Description,
		Owner:       strconv.FormatInt(post.Owner, 10),
		AuthorID:    post.AuthorID,
	}
}
