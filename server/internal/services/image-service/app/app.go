package app

import (
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	ProcessUploadedImage command.ProcessUploadedImage
}

type Queries struct {
	GetImageSizes query.GetImageSizes
}
