package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	post2 "github.com/MKKL1/schematic-app/server/internal/pkg/client/post"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/post"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/user"
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
			Port:     "1324",
			BasePath: "/",
			Timeout:  10000,
			Host:     "localhost",
		})
		e.Use(http.EchoErrorMiddleware)

		userService := client.NewUsersClient(ctx, ":8001")
		userController := user.NewController(userService)
		user.RegisterRoutes(e, userController)

		postService := post2.NewPostClient(ctx, ":8002")
		//categService := client.NewCategoryClient(ctx, ":8003")
		postController := post.NewController(postService)
		post.RegisterRoutes(e, postController)

		user.RegisterErrorMappers()
		post.RegisterErrorMappers()
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
