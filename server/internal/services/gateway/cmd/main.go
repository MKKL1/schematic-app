package main

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/server"
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

	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}
