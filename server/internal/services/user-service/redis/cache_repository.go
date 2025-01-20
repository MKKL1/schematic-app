package redis

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/google/uuid"
	"github.com/redis/rueidis/rueidisaside"
	"time"
)

type CacheRepository struct {
	baseRepo    user.Repository
	cacheClient *rueidisaside.Client
}

func NewCacheRepository(baseRepo user.Repository, cacheClient *rueidisaside.Client) CacheRepository {
	return CacheRepository{baseRepo, cacheClient}
}

func (r CacheRepository) FindById(ctx context.Context, id user.UserID) (db.User, error) {
	val, err := r.cacheClient.Get(context.Background(), time.Minute, "usr:"+id.String(), func(ctx context.Context, key string) (val string, err error) {

	})
}

func (r CacheRepository) FindByOidcSub(ctx context.Context, oidcSub uuid.UUID) (db.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r CacheRepository) FindByName(ctx context.Context, name string) (db.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r CacheRepository) CreateUser(ctx context.Context, user user.User) (int64, error) {
	//TODO implement me
	panic("implement me")
}
