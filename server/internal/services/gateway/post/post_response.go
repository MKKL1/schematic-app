package post

import (
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"strconv"
)

type PostResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description *string              `json:"desc"`
	Owner       string               `json:"owner"`
	AuthorID    *string              `json:"author"`
	Categories  []CategoriesResponse `json:"categories"`
}

type CategoriesResponse struct {
	Name string          `json:"name"`
	Vars json.RawMessage `json:"vars"`
}

func PostToResponse(post client.Post, vars []client.PostCategoryVars) PostResponse {
	var authorID *string
	if post.AuthorID != nil {
		aInt := strconv.FormatInt(*post.AuthorID, 10)
		authorID = &aInt
	}

	categories := make([]CategoriesResponse, len(vars))
	for i, v := range vars {
		categories[i] = CategoryToResponse(v)
	}

	return PostResponse{
		ID:          strconv.FormatInt(post.ID, 10),
		Name:        post.Name,
		Description: post.Description,
		Owner:       strconv.FormatInt(post.Owner, 10),
		AuthorID:    authorID,
		Categories:  categories,
	}
}

func CategoryToResponse(vars client.PostCategoryVars) CategoriesResponse {
	return CategoriesResponse{
		Name: vars.Category,
		Vars: vars.Values,
	}
}
