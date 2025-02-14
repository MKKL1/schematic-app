package client

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
)

type Post struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
	Vars        json.RawMessage
	Tags        []string
}

type CreatePostParams struct {
	Name        string
	Description *string
	AuthorID    *int64
	Sub         uuid.UUID
	Categories  []CreateCategoryMetadataParams
	Tags        []string
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

	createdId, err := p.grpcClient.CreatePost(ctx, &genproto.CreatePostRequest{
		Name:        params.Name,
		Description: params.Description,
		AuthorId:    params.AuthorID,
		AuthSub:     subBytes,
		Categories:  categBytes,
		Tags:        tags,
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

	return postProtoToDto(post), nil
}

func postProtoToDto(post *genproto.Post) Post {
	tags := make([]string, len(post.GetTags()))
	for i, v := range post.Tags {
		tags[i] = v.Tag
	}

	return Post{
		ID:          post.Id,
		Name:        post.Name,
		Description: post.Description,
		Owner:       post.Owner,
		AuthorID:    post.Author,
		Vars:        post.Categories,
		Tags:        tags,
	}
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
