package client

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/google/uuid"
)

type Post struct {
	ID          int64
	Name        string
	Description *string
	Owner       int64
	AuthorID    *int64
}

type CreatePostParams struct {
	Name        string
	Description *string
	AuthorName  *string
	AuthorID    *int64
	Sub         uuid.UUID
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

	createdId, err := p.grpcClient.CreatePost(ctx, &genproto.CreatePostRequest{
		Name:        params.Name,
		Description: params.Description,
		AuthorName:  params.AuthorName,
		AuthorId:    params.AuthorID,
		AuthSub:     subBytes,
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
	return Post{
		ID:          post.Id,
		Name:        post.Name,
		Description: post.Description,
		Owner:       post.Owner,
		AuthorID:    post.Author,
	}
}
