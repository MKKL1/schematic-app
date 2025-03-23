package ports

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	genproto.UnimplementedPostServiceServer
	app app.Application
}

func NewGrpcServer(app app.Application) *GrpcServer {
	return &GrpcServer{app: app}
}

func (g GrpcServer) GetPostById(ctx context.Context, request *genproto.PostByIdRequest) (*genproto.Post, error) {
	dto, err := g.app.Queries.GetPostById.Handle(ctx, query.GetPostByIdParams{Id: request.GetId()})
	if err != nil {
		return nil, err
	}

	return dtoToProto(dto)
}

func (g GrpcServer) CreatePost(ctx context.Context, request *genproto.CreatePostRequest) (*genproto.CreatePostResponse, error) {
	sub, err := uuid.FromBytes(request.AuthSub)
	if err != nil {
		return nil, err
	}

	var categMetadataList []command.CategoryMetadataParams
	err = sonic.Unmarshal(request.Categories, &categMetadataList)
	if err != nil {
		return nil, err
	}

	tags := make([]string, len(request.Tags))
	for i, tag := range request.Tags {
		tags[i] = tag.Tag
	}

	files := make([]command.CreatePostFileParams, len(request.Files))
	for i, f := range request.Files {
		tId, err := uuid.Parse(f.TempId)
		if err != nil {
			//TODO return invalid argument
			return nil, err
		}
		files[i] = command.CreatePostFileParams{
			TempId: tId,
		}
	}

	createdId, err := g.app.Commands.CreatePost.Handle(ctx, command.CreatePostParams{
		Name:        request.Name,
		Description: request.Description,
		AuthorID:    request.AuthorId,
		Sub:         sub,
		Categories:  categMetadataList,
		Tags:        tags,
		Files:       files,
	})
	if err != nil {
		return nil, err
	}

	return &genproto.CreatePostResponse{Id: createdId}, nil
}

func dtoToProto(dto post.Post) (*genproto.Post, error) {
	p := &genproto.Post{
		Id:    dto.ID,
		Name:  dto.Name,
		Owner: dto.Owner,
	}

	if dto.Description != nil {
		p.Description = proto.String(*dto.Description)
	}
	if dto.AuthorID != nil {
		p.Author = proto.Int64(*dto.AuthorID)
	}

	p.Categories = dto.Categories

	p.Tags = make([]*genproto.Tag, len(dto.Tags))
	for i, t := range dto.Tags {
		p.Tags[i] = &genproto.Tag{Tag: t}
	}
	p.Files = make([]*genproto.PostFile, len(dto.Files))
	for i, f := range dto.Files {
		if f.State == post.FileAvailable {
			//Hash, Downloads, FileSize and UpdatedAt should not be null here
			if f.Hash == nil || f.Downloads == nil || f.FileSize == nil || f.UpdatedAt == nil {
				return nil, fmt.Errorf("invalid file state")
			}
			p.Files[i] = &genproto.PostFile{
				State: &genproto.PostFile_Processed{
					Processed: &genproto.ProcessedPostFile{
						Hash:      *f.Hash,
						Name:      f.Name,
						Downloads: *f.Downloads,
						FileSize:  *f.FileSize,
						UpdatedAt: timestamppb.New(*f.UpdatedAt),
					},
				},
			}
		} else if f.State == post.FilePending {
			p.Files[i] = &genproto.PostFile{
				State: &genproto.PostFile_Pending{
					Pending: &genproto.PendingPostFile{
						Name: f.Name,
					},
				},
			}
		} else {
			return nil, fmt.Errorf("invalid file state")
		}
	}

	return p, nil
}
