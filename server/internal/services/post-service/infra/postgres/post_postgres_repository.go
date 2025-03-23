package postgres

import (
	"context"
	"errors"
	"fmt"
	errorDB "github.com/MKKL1/schematic-app/server/internal/pkg/db"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/postgres/db"
	"github.com/bytedance/sonic"
	"github.com/jackc/pgx/v5"
	"time"
)

type attachedFileModel struct {
	Hash      *string   `json:"hash"`
	TempID    *string   `json:"temp_id"`
	Name      string    `json:"name"`
	FileSize  int32     `json:"file_size"`
	Downloads int32     `json:"downloads"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostPostgresRepository struct {
	queries *db.Queries
}

func NewPostPostgresRepository(queries *db.Queries) *PostPostgresRepository {
	return &PostPostgresRepository{queries}
}

func (p PostPostgresRepository) FindById(ctx context.Context, id int64) (post.Post, error) {
	row, err := p.queries.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return post.Post{}, errorDB.ErrNoRows
		}
		return post.Post{}, err
	}

	var categoryVars post.PostCategories
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

	var fileModel []attachedFileModel
	if row.Files != nil {
		s, ok := row.Files.(string)
		if !ok {
			return post.Post{}, errors.New("invalid type for Files")
		}
		err := sonic.UnmarshalString(s, fileModel)
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
		}
		files[i] = post.PostFile{
			Hash:      file.Hash,
			Name:      file.Name,
			Downloads: &file.Downloads,
			FileSize:  &file.FileSize,
			UpdatedAt: &file.UpdatedAt,
			State:     post.FileAvailable,
		}
	}

	postEntity := post.Post{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Owner:       row.Owner,
		AuthorID:    row.AuthorID,
		Categories:  categoryVars,
		Tags:        tags,
		Files:       files,
	}

	return postEntity, nil
}

func (p PostPostgresRepository) Create(ctx context.Context, params post.CreatePostParams) error {
	categoriesJSON, err := sonic.Marshal(params.Categories)
	if err != nil {
		return err
	}

	filesJSON, err := sonic.Marshal(params.Files)
	if err != nil {
		return err
	}

	err = p.queries.CreatePost(ctx, db.CreatePostParams{
		ID:       params.ID,
		Name:     params.Name,
		Desc:     params.Description,
		Owner:    params.Owner,
		AuthorID: params.AuthorID,
		Column6:  params.Tags,
		Column7:  categoriesJSON,
		Column8:  filesJSON,
	})
	if err != nil {
		return err
	}
	return nil
}

func (p PostPostgresRepository) GetCountForTag(ctx context.Context, tag string) (int64, error) {
	//TODO implement me
	panic("implement me")
}
