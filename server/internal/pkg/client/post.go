package client

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
	Categories  []PostCategories
	Tags        []string
	File        []PostFile
}

type PostFile struct {
	Hash      *string
	Name      string
	Downloads *int32
	FileSize  *int32
	UpdatedAt *time.Time
	State     string
}

type PostCategories struct {
	Name     string
	Metadata map[string]interface{}
}

type CreatePostParams struct {
	Name        string
	Description *string
	AuthorID    *int64
	Sub         uuid.UUID
	Categories  []CreateCategoryMetadataParams
	Tags        []string
	Files       []CreatePostFileParams
}

type CreatePostFileParams struct {
	TempId uuid.UUID
}

type CreateCategoryMetadataParams struct {
	Name     string
	Metadata map[string]interface{}
}

type PostApplication struct {
	Command PostCommandService
	Query   PostQueryService
}

type PostCommandService interface {
	CreatePost(ctx context.Context, params CreatePostParams) (int64, error)
}

type PostQueryService interface {
	GetPostById(ctx context.Context, id int64) (Post, error)
}

type PostCommandGrpcService struct {
	grpcClient genproto.PostServiceClient
}

func NewPostCommandGrpcService(grpcClient genproto.PostServiceClient) *PostCommandGrpcService {
	return &PostCommandGrpcService{grpcClient: grpcClient}
}

func (p PostCommandGrpcService) CreatePost(ctx context.Context, params CreatePostParams) (int64, error) {
	subBytes, err := params.Sub.MarshalBinary()
	if err != nil {
		return 0, err
	}

	categBytes, err := sonic.Marshal(params.Categories)
	if err != nil {
		return 0, err
	}

	tags := make([]*genproto.Tag, len(params.Tags))
	for i, tag := range params.Tags {
		tags[i] = &genproto.Tag{
			Tag: tag,
		}
	}

	files := make([]*genproto.File, len(params.Files))
	for i, f := range params.Files {
		files[i] = &genproto.File{
			TempId: f.TempId.String(),
		}
	}

	createdId, err := p.grpcClient.CreatePost(ctx, &genproto.CreatePostRequest{
		Name:        params.Name,
		Description: params.Description,
		AuthorId:    params.AuthorID,
		AuthSub:     subBytes,
		Categories:  categBytes,
		Tags:        tags,
		Files:       files,
	})
	if err != nil {
		return 0, err
	}

	return createdId.Id, err
}

type PostQueryGrpcService struct {
	grpcClient genproto.PostServiceClient
}

func NewPostQueryGrpcService(grpcClient genproto.PostServiceClient) *PostQueryGrpcService {
	return &PostQueryGrpcService{grpcClient: grpcClient}
}

func (p PostQueryGrpcService) GetPostById(ctx context.Context, id int64) (Post, error) {
	post, err := p.grpcClient.GetPostById(ctx, &genproto.PostByIdRequest{Id: id})
	if err != nil {
		return Post{}, err
	}

	return postProtoToDto(post)
}

func postProtoToDto(post *genproto.Post) (Post, error) {
	tags := make([]string, len(post.GetTags()))
	for i, v := range post.Tags {
		tags[i] = v.Tag
	}

	var categories []PostCategories
	err := sonic.Unmarshal(post.Categories, &categories)
	if err != nil {
		return Post{}, err
	}

	files := make([]PostFile, len(post.Files))
	for i, f := range post.Files {
		//TODO implement
		files[i] = PostFile{
			Hash:      ,
			Name:      "",
			Downloads: nil,
			FileSize:  nil,
			UpdatedAt: nil,
			State:     "",
		}
	}
	return Post{
		ID:          post.Id,
		Name:        post.Name,
		Description: post.Description,
		Owner:       post.Owner,
		AuthorID:    post.Author,
		Categories:  categories,
		Tags:        tags,
		File:
	}, nil
}

func NewPostClient(ctx context.Context, addr string) PostApplication {
	conn := NewConnection(ctx, addr)

	service := genproto.NewPostServiceClient(conn)
	query := PostQueryGrpcService{
		grpcClient: service,
	}
	command := PostCommandGrpcService{
		grpcClient: service,
	}

	return PostApplication{
		Query:   query,
		Command: command,
	}
}
