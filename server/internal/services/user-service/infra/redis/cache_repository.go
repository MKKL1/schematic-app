package redis

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/rueidisaside"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	protoModel "github.com/MKKL1/schematic-app/server/internal/services/user-service/infra/proto"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
	"time"
)

// I don't like how ristretto is used without any abstraction
// It would be better to have some kind of "class" that handles mapping between given values and id
// This solution is very hard to extend
// It also doesn't share data between services using redis

type CacheRepository struct {
	baseRepo    user.Repository
	cacheClient rueidisaside.CacheAsideClient
	typedClient rueidisaside.TypedCacheAsideClient[user.Entity]
	subCache    *ristretto.Cache[string, int64]
	nameCache   *ristretto.Cache[string, int64]
}

func NewCacheRepository(baseRepo user.Repository, cacheClient rueidisaside.CacheAsideClient) CacheRepository {
	subCache, err := ristretto.NewCache(&ristretto.Config[string, int64]{
		NumCounters: 10000,
		MaxCost:     1000,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	defer subCache.Close()

	nameCache, err := ristretto.NewCache(&ristretto.Config[string, int64]{
		NumCounters: 10000,
		MaxCost:     1000,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	defer nameCache.Close()

	typedClient := rueidisaside.NewTypedCacheAsideClient(
		cacheClient,
		func(user *user.Entity) (string, error) {
			pbModel, err := protoModel.FromEntity(user)
			if err != nil {
				return "", err
			}

			marshal, err := proto.Marshal(pbModel)
			if err != nil {
				return "", err
			}
			return string(marshal), err
		},
		func(d string) (*user.Entity, error) {
			pbModel := protoModel.UserEntity{}
			err := proto.Unmarshal([]byte(d), &pbModel)
			if err != nil {
				return nil, err
			}

			return protoModel.ToEntity(&pbModel)
		},
	)

	return CacheRepository{baseRepo, cacheClient, typedClient, subCache, nameCache}
}

func (c CacheRepository) findById(ctx context.Context, id int64) (user.Entity, error) {
	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usr:", id), func(ctx context.Context, key string) (val *user.Entity, err error) {
		userModel, err := c.baseRepo.FindById(ctx, id)
		if err != nil {
			return nil, err
		}

		return &userModel, nil
	})
	if err != nil {
		return user.Entity{}, err
	}
	if val == nil {
		return user.Entity{}, fmt.Errorf("user not found")
	}
	return *val, err
}

func (c CacheRepository) FindById(ctx context.Context, id int64) (user.Entity, error) {
	return c.findById(ctx, id)
}

// NOTE: Right now each query saves a copy of user instead of one user object that is shared across query methods
//TODO find id of user based on sub and name in cache

func (c CacheRepository) FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (user.Entity, error) {
	oidcSubStr := oidcSub.String()

	userId, exists := c.subCache.Get(oidcSubStr)
	if exists {
		return c.findById(ctx, userId)
	}

	//Cache new
	userModel, err := c.baseRepo.FindByOidcSub(ctx, oidcSub)
	if err != nil {
		return user.Entity{}, err
	}
	userId = userModel.ID
	c.subCache.SetWithTTL(oidcSubStr, userId, 1, 30*time.Minute)

	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usr:", userId), func(ctx context.Context, key string) (val *user.Entity, err error) {
		//If user model is not already cached, use the object that is already fetched
		return &userModel, nil
	})

	if err != nil {
		return user.Entity{}, err
	}
	if val == nil {
		return user.Entity{}, fmt.Errorf("user not found")
	}
	return *val, err
}

func (c CacheRepository) FindByName(ctx context.Context, name string) (user.Entity, error) {
	userId, exists := c.nameCache.Get(name)
	if exists {
		return c.findById(ctx, userId)
	}

	//Cache new
	userModel, err := c.baseRepo.FindByName(ctx, name)
	if err != nil {
		return user.Entity{}, err
	}
	userId = userModel.ID
	c.nameCache.SetWithTTL(name, userId, 1, 30*time.Minute)

	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usr:", userId), func(ctx context.Context, key string) (val *user.Entity, err error) {
		//If user model is not already cached, use the object that is already fetched
		return &userModel, nil
	})

	if err != nil {
		return user.Entity{}, err
	}
	if val == nil {
		return user.Entity{}, fmt.Errorf("user not found")
	}
	return *val, err
}

func (c CacheRepository) CreateUser(ctx context.Context, user user.Entity) (int64, error) {
	//Nothing to cache
	return c.baseRepo.CreateUser(ctx, user)
}
