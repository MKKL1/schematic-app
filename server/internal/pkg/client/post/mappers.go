package post

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/bytedance/sonic"
)

func ProtoToDto(post *genproto.Post) (PostDto, error) {
	tags := make([]string, len(post.GetTags()))
	for i, v := range post.Tags {
		tags[i] = v.Tag
	}

	var categories []PostCategoryDto
	err := sonic.Unmarshal(post.Categories, &categories)
	if err != nil {
		return PostDto{}, err
	}

	files := make([]PostFileDto, len(post.Files))
	for i, f := range post.Files {
		switch op := f.State.(type) {
		case *genproto.PostFile_Pending:
			files[i] = PostFileDto{
				Name:  op.Pending.GetName(),
				State: FilePending,
			}
		case *genproto.PostFile_Processed:
			upd := op.Processed.UpdatedAt.AsTime()
			files[i] = PostFileDto{
				Hash:      &op.Processed.Hash,
				Name:      op.Processed.Name,
				Downloads: &op.Processed.Downloads,
				FileSize:  &op.Processed.FileSize,
				UpdatedAt: &upd,
				State:     FileAvailable,
			}
		}
	}
	return PostDto{
		ID:          post.Id,
		Name:        post.Name,
		Description: post.Description,
		Owner:       post.Owner,
		AuthorID:    post.Author,
		Categories:  categories,
		Tags:        tags,
		File:        files,
	}, nil
}

func CreatePostRequestDtoToProto(dto CreatePostRequest) (*genproto.CreatePostRequest, error) {
	subBytes, err := dto.Sub.MarshalBinary()
	if err != nil {
		return nil, err
	}

	categBytes, err := sonic.Marshal(dto.Categories)
	if err != nil {
		return nil, err
	}

	tags := make([]*genproto.Tag, len(dto.Tags))
	for i, tag := range dto.Tags {
		tags[i] = &genproto.Tag{
			Tag: tag,
		}
	}

	files := make([]*genproto.File, len(dto.Files))
	for i, f := range dto.Files {
		files[i] = &genproto.File{
			TempId: f.String(),
		}
	}

	return &genproto.CreatePostRequest{
		Name:        dto.Name,
		Description: dto.Description,
		AuthorId:    dto.AuthorID,
		AuthSub:     subBytes,
		Categories:  categBytes,
		Tags:        tags,
		Files:       files,
	}, nil
}
