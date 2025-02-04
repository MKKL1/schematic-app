package app

import "github.com/MKKL1/schematic-app/server/internal/services/tag-service/app/command"

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateCategoryVars command.CreateCategoryVarsHandler
}

type Queries struct {
}
