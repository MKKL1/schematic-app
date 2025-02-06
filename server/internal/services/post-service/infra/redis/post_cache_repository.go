package redis

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/rueidisaside"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	postProto "github.com/MKKL1/schematic-app/server/internal/services/post-service/infra/proto"
	"github.com/golang/protobuf/proto"
	"time"
)

type PostCacheRepository struct {
	baseRepo    post.Repository
	typedClient rueidisaside.TypedCacheAsideClient[post.Entity]
}

func NewPostCacheRepository(baseRepo post.Repository, cacheClient rueidisaside.CacheAsideClient) *PostCacheRepository {
	typedClient := rueidisaside.NewTypedCacheAsideClient(
		cacheClient,
		func(dbPost *post.Entity) (string, error) {
			pbModel := postProto.FromEntity(*dbPost)
			marshal, err := proto.Marshal(&pbModel)
			if err != nil {
				return "", err
			}
			return string(marshal), err
		},
		func(d string) (*post.Entity, error) {
			pbModel := postProto.PostEntity{}
			err := proto.Unmarshal([]byte(d), &pbModel)
			if err != nil {
				return nil, err
			}

			postModel := postProto.ToEntity(&pbModel)
			return &postModel, nil
		},
	)

	return &PostCacheRepository{baseRepo: baseRepo, typedClient: typedClient}
}

func (p PostCacheRepository) FindById(ctx context.Context, id int64) (post.Entity, error) {
	val, err := p.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("post:", id), func(ctx context.Context, key string) (val *post.Entity, err error) {
		postModel, err := p.baseRepo.FindById(ctx, id)
		if err != nil {
			return nil, err
		}

		return &postModel, nil
	})
	if err != nil {
		return post.Entity{}, err
	}
	if val == nil {
		return post.Entity{}, fmt.Errorf("post not found")
	}
	return *val, err
}

func (p PostCacheRepository) Create(ctx context.Context, model post.Entity) error {
	return p.baseRepo.Create(ctx, model)
}
