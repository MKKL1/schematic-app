package http

import (
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"strconv"
)

// AuthorResponse This type is used both as a response and in create request
type AuthorResponse struct {
	Name *string `json:"name" validate:"omitempty,alphanum"`
	ID   *string `json:"id" validate:"omitempty,number"`
}

func PostToAuthorResponse(post post.Post) *AuthorResponse {
	if post.Author == nil {
		return nil
	}
	if post.Author.IsKnown {
		id := strconv.FormatInt(post.Author.UserID, 10)
		return &AuthorResponse{
			Name: nil,
			ID:   &id,
		}
	}
	return &AuthorResponse{
		Name: &post.Author.Name,
		ID:   nil,
	}
}
