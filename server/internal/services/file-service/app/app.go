package app

import "github.com/MKKL1/schematic-app/server/internal/services/file-service/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	UploadTempFile     command.UploadTempFileHandler
	DeleteExpiredFiles command.DeleteExpiredFilesHandler
	CommitTempFile     command.CommitTempHandler
	PostCreatedHandler command.PostCreatedHandler // Handles the PostCreated event
	ProcessImage       command.ProcessImageHandler
}

type Queries struct {
}
