package redis

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/rueidisaside"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	postProto "github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post/proto"
	"github.com/golang/protobuf/proto"
	"time"
)

type PostCacheRepository struct {
	baseRepo    post.Repository
	typedClient rueidisaside.TypedCacheAsideClient[post.Model]
}

func NewPostCacheRepository(baseRepo post.Repository, cacheClient rueidisaside.CacheAsideClient) *PostCacheRepository {
	typedClient := rueidisaside.NewTypedCacheAsideClient(
		cacheClient,
		func(dbPost *post.Model) (string, error) {
			pbModel := postProto.PostModel{
				Id:    dbPost.ID,
				Desc:  dbPost.Description,
				Owner: dbPost.Owner,
				AName: dbPost.AuthorName,
				AId:   dbPost.AuthorID,
			}
			marshal, err := proto.Marshal(&pbModel)
			if err != nil {
				return "", err
			}
			return string(marshal), err
		},
		func(d string) (*post.Model, error) {
			pbModel := postProto.PostModel{}
			err := proto.Unmarshal([]byte(d), &pbModel)
			if err != nil {
				return nil, err
			}

			postModel := &post.Model{
				ID:          pbModel.Id,
				Description: pbModel.Desc,
				Owner:       pbModel.Owner,
				AuthorName:  pbModel.AName,
				AuthorID:    pbModel.AId,
			}

			return postModel, nil
		},
	)

	return &PostCacheRepository{baseRepo: baseRepo, typedClient: typedClient}
}

func (p PostCacheRepository) FindById(ctx context.Context, id int64) (post.Model, error) {
	val, err := p.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("post:", id), func(ctx context.Context, key string) (val *post.Model, err error) {
		postModel, err := p.baseRepo.FindById(ctx, id)
		if err != nil {
			return nil, err
		}

		return &postModel, nil
	})
	if err != nil {
		return post.Model{}, err
	}
	if val == nil {
		return post.Model{}, fmt.Errorf("post not found")
	}
	return *val, err
}
