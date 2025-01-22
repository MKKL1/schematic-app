package redis

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/rueidisaside"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user/data"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"time"
)

type CacheRepository struct {
	baseRepo    user.Repository
	cacheClient rueidisaside.CacheAsideClient
	typedClient rueidisaside.TypedCacheAsideClient[user.Model]
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
		func(user *user.Model) (string, error) {
			sub, err := user.OidcSub.MarshalBinary()
			if err != nil {
				return "", err
			}
			pbModel := data.UserModel{
				Id:      user.ID,
				Name:    user.Name,
				OidcSub: sub,
			}
			marshal, err := proto.Marshal(&pbModel)
			if err != nil {
				return "", err
			}
			return string(marshal), err
		},
		func(d string) (*user.Model, error) {
			pbModel := data.UserModel{}
			err := proto.Unmarshal([]byte(d), &pbModel)
			if err != nil {
				return nil, err
			}

			sub, err := uuid.FromBytes(pbModel.OidcSub)
			if err != nil {
				return nil, err
			}

			userModel := &user.Model{
				ID:      pbModel.Id,
				Name:    pbModel.Name,
				OidcSub: sub,
			}

			return userModel, nil
		},
	)

	return CacheRepository{baseRepo, cacheClient, typedClient, subCache, nameCache}
}

func (c CacheRepository) findById(ctx context.Context, id int64) (user.Model, error) {
	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usr:", id), func(ctx context.Context, key string) (val *user.Model, err error) {
		userModel, err := c.baseRepo.FindById(ctx, id)
		if err != nil {
			return nil, err
		}

		return &userModel, nil
	})
	if val == nil {
		return user.Model{}, fmt.Errorf("user not found")
	}
	return *val, err
}

func (c CacheRepository) FindById(ctx context.Context, id int64) (user.Model, error) {
	return c.findById(ctx, id)
}

// NOTE: Right now each query saves a copy of user instead of one user object that is shared across query methods
//TODO find id of user based on sub and name in cache

func (c CacheRepository) FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (user.Model, error) {
	oidcSubStr := oidcSub.String()

	userId, exists := c.subCache.Get(oidcSubStr)
	if exists {
		return c.findById(ctx, userId)
	}

	//Cache new
	userModel, err := c.baseRepo.FindByOidcSub(ctx, oidcSub)
	if err != nil {
		return user.Model{}, err
	}
	userId = userModel.ID
	c.subCache.SetWithTTL(oidcSubStr, userId, 1, 30*time.Minute)

	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usr:", userId), func(ctx context.Context, key string) (val *user.Model, err error) {
		//If user model is not already cached, use the object that is already fetched
		return &userModel, nil
	})

	if val == nil {
		return user.Model{}, fmt.Errorf("user not found")
	}
	return *val, err
}

func (c CacheRepository) FindByName(ctx context.Context, name string) (user.Model, error) {
	userId, exists := c.nameCache.Get(name)
	if exists {
		return c.findById(ctx, userId)
	}

	//Cache new
	userModel, err := c.baseRepo.FindByName(ctx, name)
	if err != nil {
		return user.Model{}, err
	}
	userId = userModel.ID
	c.nameCache.SetWithTTL(name, userId, 1, 30*time.Minute)

	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usr:", userId), func(ctx context.Context, key string) (val *user.Model, err error) {
		//If user model is not already cached, use the object that is already fetched
		return &userModel, nil
	})

	if val == nil {
		return user.Model{}, fmt.Errorf("user not found")
	}
	return *val, err
}

func (c CacheRepository) CreateUser(ctx context.Context, user user.Model) (int64, error) {
	//Nothing to cache
	return c.baseRepo.CreateUser(ctx, user)
}
