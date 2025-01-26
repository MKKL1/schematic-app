package app

import (
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreatePost command.CreatePostHandler
}

type Queries struct {
	GetPostById query.GetPostByIdHandler
}
