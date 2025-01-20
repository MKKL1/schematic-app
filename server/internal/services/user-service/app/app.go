package app

import (
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateUser command.CreateUserHandler
}

type Queries struct {
	GetUserById  query.GetUserByIdHandler
	GetUserBySub query.GetUserBySubHandler
}
