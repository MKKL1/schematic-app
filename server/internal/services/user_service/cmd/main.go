package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/config"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/http"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/services"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	queries := config.ConfigDB(ctx)

	userRepo := postgres.NewUserPostgresRepository(queries)
	userService := services.NewUserService(userRepo)
	httpServer := http.NewHttpServer(userService)

	e := echo.New()

	config.ConfigMiddlewares(e)
	config.ConfigRoutes(e, httpServer)

	e.Logger.Fatal(e.Start(":1323"))

}
