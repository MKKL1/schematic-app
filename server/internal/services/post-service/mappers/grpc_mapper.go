package mappers

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/bytedance/sonic"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreatePostRequestProtoToCmd(request *genproto.CreatePostRequest) (command.CreatePost, error) {
	sub, err := uuid.FromBytes(request.AuthSub)
	if err != nil {
		return command.CreatePost{}, err
	}

	var categMetadataList []command.CreatePostCategory
	err = sonic.Unmarshal(request.Categories, &categMetadataList)
	if err != nil {
		return command.CreatePost{}, err
	}

	tags := make([]string, len(request.Tags))
	for i, tag := range request.Tags {
		tags[i] = tag.Tag
	}

	files := make([]uuid.UUID, len(request.Files))
	for i, f := range request.Files {
		tId, err := uuid.Parse(f.TempId)
		if err != nil {
			//TODO return invalid argument
			return command.CreatePost{}, err
		}
		files[i] = tId
	}

	return command.CreatePost{
		Name:        request.Name,
		Description: request.Description,
		AuthorID:    request.AuthorId,
		Sub:         sub,
		Categories:  categMetadataList,
		Tags:        tags,
		Files:       files,
	}, nil
}

func AppToProto(domainPost post.Post) (*genproto.Post, error) {
	p := &genproto.Post{
		Id:    domainPost.ID,
		Name:  domainPost.Name,
		Owner: domainPost.Owner,
	}

	if domainPost.Description != nil {
		p.Description = proto.String(*domainPost.Description)
	}
	if domainPost.AuthorID != nil {
		p.Author = proto.Int64(*domainPost.AuthorID)
	}

	p.Categories = domainPost.Categories

	p.Tags = make([]*genproto.Tag, len(domainPost.Tags))
	for i, t := range domainPost.Tags {
		p.Tags[i] = &genproto.Tag{Tag: t}
	}
	p.Files = make([]*genproto.PostFile, len(domainPost.Files))
	for i, f := range domainPost.Files {
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

//func ProtoToApp(protoPost *genproto.Post) (post.Post, error) {
//
//}
