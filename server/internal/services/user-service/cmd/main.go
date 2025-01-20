package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/http/middlewares"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/http"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/services"
	"os"
	"os/signal"
	"time"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		e := server.NewEchoServer()
		server.RunHttpServer(ctx, e, &server.EchoConfig{
			Port:     "1323",
			BasePath: "/",
			Timeout:  10000,
			Host:     "localhost",
		})

		e.HTTPErrorHandler = middlewares.HTTPErrorHandler(http.MapAppError)

		dbPool, err := server.NewPostgreSQL(ctx, &server.PostgresConfig{
			Port:     "5432",
			Host:     "localhost",
			Username: "root",
			Password: "root",
			Database: "sh_user",
		})
		if err != nil {
			panic(err)
		}

		queries := db.New(dbPool)

		userRepo := postgres.NewUserPostgresRepository(queries)
		userService := services.NewUserService(userRepo)
		userController := http.NewUserController(userService)
		http.RegisterRoutes(e, userController)
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
