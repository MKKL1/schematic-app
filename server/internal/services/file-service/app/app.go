package app

import "github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	UploadTempFile command.UploadTempFileHandler
}

type Queries struct {
}
