package post

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/client/post"
	"strconv"
	"time"
)

type PostResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description *string                  `json:"desc"`
	Owner       string                   `json:"owner_id"`
	AuthorID    *string                  `json:"author_id"`
	Categories  []PostCategoriesResponse `json:"categories"`
	Tags        []string                 `json:"tags"`
	Files       []PostFilesResponse      `json:"files"`
}

type PostCategoriesResponse struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type PostFilesResponse struct {
	Hash      *string
	Name      string
	Downloads *int32
	FileSize  *int32
	UpdatedAt *time.Time
}

func PostToResponse(dto post.PostDto) PostResponse {
	var authorID *string
	if dto.AuthorID != nil {
		aInt := strconv.FormatInt(*dto.AuthorID, 10)
		authorID = &aInt
	}

	tags := make([]string, len(dto.Tags))
	for i, v := range dto.Tags {
		tags[i] = v
	}

	categs := make([]PostCategoriesResponse, len(dto.Categories))
	for i, v := range dto.Categories {
		categs[i] = PostCategoriesResponse{
			Name:     v.Name,
			Metadata: v.Metadata,
		}
	}

	files := make([]PostFilesResponse, len(dto.File))
	for i, v := range dto.File {
		//If state is pending, only Name won't be nil
		files[i] = PostFilesResponse{
			Hash:      v.Hash,
			Name:      v.Name,
			Downloads: v.Downloads,
			FileSize:  v.FileSize,
			UpdatedAt: v.UpdatedAt,
		}
	}

	return PostResponse{
		ID:          strconv.FormatInt(dto.ID, 10),
		Name:        dto.Name,
		Description: dto.Description,
		Owner:       strconv.FormatInt(dto.Owner, 10),
		AuthorID:    authorID,
		Categories:  categs,
		Tags:        tags,
	}
}
