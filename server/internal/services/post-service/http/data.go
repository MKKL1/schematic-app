package http

import (
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"strconv"
)

type PostResponse struct {
	ID          string          `json:"id"`
	Description *string         `json:"desc"`
	Owner       string          `json:"owner"`
	Author      *AuthorResponse `json:"author"`
}

type AuthorResponse struct {
	Type AuthorType `json:"type"`
	ID   string     `json:"id"`
}

type AuthorType string

var (
	AuthorTypeKnownUser   AuthorType = "known_user"
	AuthorTypeUnknownUser AuthorType = "unknown_user"
)

func PostToResponse(post post.Post) PostResponse {
	return PostResponse{
		ID:          strconv.FormatInt(post.ID, 10),
		Description: post.Description,
		Owner:       strconv.FormatInt(post.Owner, 10),
		Author:      PostToAuthorResponse(post),
	}
}

func PostToAuthorResponse(post post.Post) *AuthorResponse {
	var authorType AuthorType
	var authorID string
	if post.Author == nil {
		return nil
	}

	if post.Author.IsKnown {
		authorType = AuthorTypeKnownUser
		authorID = strconv.FormatInt(post.Author.UserID, 10)
	} else {
		authorType = AuthorTypeUnknownUser
		authorID = post.Author.Name
	}

	return &AuthorResponse{
		Type: authorType,
		ID:   authorID,
	}
}
