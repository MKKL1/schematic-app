package proto

import (
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/post"
)

func FromEntity(entity post.Entity) PostEntity {
	return PostEntity{
		Id:    entity.ID,
		Name:  entity.Name,
		Desc:  entity.Description,
		Owner: entity.Owner,
		AId:   entity.AuthorID,
	}
}

func ToEntity(protoEntity *PostEntity) post.Entity {
	return post.Entity{
		ID:          protoEntity.Id,
		Name:        protoEntity.Name,
		Description: protoEntity.Desc,
		Owner:       protoEntity.Owner,
		AuthorID:    protoEntity.AId,
	}
}
