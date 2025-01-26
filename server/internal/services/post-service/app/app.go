package app

import "github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	GetPostById query.GetPostByIdHandler
}
