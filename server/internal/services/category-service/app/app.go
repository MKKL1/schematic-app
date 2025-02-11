package app

import (
	"github.com/MKKL1/schematic-app/server/internal/services/category-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/category-service/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateCategoryVars command.CreateCategoryVarsHandler
}

type Queries struct {
	GetCategVarsByPost query.GetCategVarsByPostHandler
}
