package mappers

import (
	"github.com/MKKL1/schematic-app/server/internal/services/file-service/domain/file"
	db "github.com/MKKL1/schematic-app/server/internal/services/file-service/postgres/db"
)

func TmpFileModelToDomain(model db.TmpFile) file.TempFile {
	return file.TempFile{
		Key:         model.StoreKey,
		FileName:    model.FileName,
		ContentType: *model.ContentType, // Assuming ContentType is non-null in DB for temp files now
		Status:      model.Status,
		ErrorReason: model.ErrorReason,
		FinalHash:   model.FinalHash,
		ExpiresAt:   model.ExpiresAt.Time,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
	}
}

func PermFileModelToDomain(model db.File) file.PermFile {
	return file.PermFile{
		Hash:        model.Hash,
		FileSize:    model.FileSize,
		ContentType: model.ContentType,
		CreatedAt:   model.CreatedAt.Time,
		UpdatedAt:   model.UpdatedAt.Time,
	}
}
