package proto

import (
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/google/uuid"
)

func FromEntity(entity *user.Entity) (*UserEntity, error) {
	sub, err := entity.OidcSub.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return &UserEntity{
		Id:      entity.ID,
		Name:    entity.Name,
		OidcSub: sub,
	}, nil
}

func ToEntity(protoEntity *UserEntity) (*user.Entity, error) {
	sub, err := uuid.FromBytes(protoEntity.OidcSub)
	if err != nil {
		return nil, err
	}

	return &user.Entity{
		ID:      protoEntity.Id,
		Name:    protoEntity.Name,
		OidcSub: sub,
	}, nil
}
