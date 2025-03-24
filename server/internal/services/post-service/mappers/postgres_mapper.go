package mappers

import (
	"errors"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/postgres/db"
	"github.com/bytedance/sonic"
	"time"
)

// PostFileModel used to deserialize postgres query
type PostFileModel struct {
	Hash      *string   `json:"hash"`
	TempID    *string   `json:"temp_id"`
	Name      string    `json:"name"`
	FileSize  int32     `json:"file_size"`
	Downloads int32     `json:"downloads"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetPostRowToApp(row db.GetPostRow) (post.Post, error) {
	var categoryVars post.PostCategoriesRaw
	if row.CategoryVars != nil {
		s, ok := row.CategoryVars.(string)
		if !ok {
			return post.Post{}, errors.New("invalid type for category var")
		}
		categoryVars = []byte(s)
	}

	var tags []string
	if row.Tags != nil {
		tagArray, ok := row.Tags.([]interface{})
		if !ok {
			return post.Post{}, fmt.Errorf("invalid type for Tags")
		}
		for _, tag := range tagArray {
			if str, ok := tag.(string); ok {
				tags = append(tags, str)
			} else {
				return post.Post{}, fmt.Errorf("invalid tag type: expected string, got %T", tag)
			}
		}
	}

	var fileModel []PostFileModel
	if row.Files != nil {
		s, ok := row.Files.(string)
		if !ok {
			return post.Post{}, errors.New("invalid type for Files")
		}
		err := sonic.UnmarshalString(s, &fileModel)
		if err != nil {
			return post.Post{}, err
		}
	}

	files := make([]post.PostFile, len(fileModel))
	for i, file := range fileModel {
		if file.Hash == nil {
			files[i] = post.PostFile{
				Name:  file.Name,
				State: post.FilePending,
			}
		} else {
			files[i] = post.PostFile{
				Hash:      file.Hash,
				Name:      file.Name,
				Downloads: &file.Downloads,
				FileSize:  &file.FileSize,
				UpdatedAt: &file.UpdatedAt,
				State:     post.FileAvailable,
			}
		}

	}

	return post.Post{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Owner:       row.Owner,
		AuthorID:    row.AuthorID,
		Categories:  categoryVars,
		Tags:        tags,
		Files:       files,
	}, nil
}

func CreatePostParamsToQuery(params post.CreatePostParams) (db.CreatePostParams, error) {
	categoriesJSON, err := sonic.Marshal(params.Categories)
	if err != nil {
		return db.CreatePostParams{}, err
	}

	filesJSON, err := sonic.Marshal(params.Files) //TODO separated model may be needed
	if err != nil {
		return db.CreatePostParams{}, err
	}

	return db.CreatePostParams{
		ID:       params.ID,
		Name:     params.Name,
		Desc:     params.Description,
		Owner:    params.Owner,
		AuthorID: params.AuthorID,
		Column6:  params.Tags,
		Column7:  categoriesJSON,
		Column8:  filesJSON,
	}, nil
}
