package proto

import (
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/google/uuid"
	"google.golang.org/protobuf/runtime/protoimpl"
)

func FromEntity(entity *post.Entity) (*PostEntity, error) {
	return &PostEntity{
		Id:    entity.ID,
		Name:  entity.Name,
		Desc:  entity.Description,
		Owner: entity.Owner,
		AName: entity.AuthorName,
		AId:   entity.AuthorID,
	}, nil
}

func ToEntity(protoEntity *PostEntity) (*post.Entity, error) {
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
