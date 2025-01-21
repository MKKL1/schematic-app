package redis

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/rueidisaside"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user/data"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"time"
)

type CacheRepository struct {
	baseRepo    user.Repository
	cacheClient rueidisaside.CacheAsideClient
	typedClient rueidisaside.TypedCacheAsideClient[user.Model]
}

func NewCacheRepository(baseRepo user.Repository, cacheClient rueidisaside.CacheAsideClient) CacheRepository {
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

	return CacheRepository{baseRepo, cacheClient, typedClient}
}

func (c CacheRepository) FindById(ctx context.Context, id int64) (user.Model, error) {
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

// NOTE: Right now each query saves a copy of user instead of one user object that is shared across query methods
//TODO find id of user based on sub and name in cache

func (c CacheRepository) FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (user.Model, error) {
	//TODO save sub to id mapping
	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usrsub:", oidcSub), func(ctx context.Context, key string) (val *user.Model, err error) {
		userModel, err := c.baseRepo.FindByOidcSub(ctx, oidcSub)
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

func (c CacheRepository) FindByName(ctx context.Context, name string) (user.Model, error) {
	//TODO save name to id mapping
	val, err := c.typedClient.Get(ctx, 10*time.Minute, fmt.Sprint("usrname:", name), func(ctx context.Context, key string) (val *user.Model, err error) {
		userModel, err := c.baseRepo.FindByName(ctx, name)
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

func (c CacheRepository) CreateUser(ctx context.Context, user user.Model) (int64, error) {
	//Do nothing
	return c.baseRepo.CreateUser(ctx, user)
}
