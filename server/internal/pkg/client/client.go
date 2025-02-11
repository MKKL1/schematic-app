package client

import (
	"context"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func NewConnection(ctx context.Context, addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  1 * time.Second,
				Multiplier: 1.6,
				MaxDelay:   30 * time.Second,
			},
			MinConnectTimeout: 5 * time.Second,
		}),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Str("addr", addr).Msg("shutting down gRPC server")
				err := conn.Close()
				if err != nil {
					log.Error().Str("addr", addr).Err(err).Msg("failed to close gRPC connection")
					return
				}
				log.Info().Msg("server shut down")
				return
			}
		}
	}()

	return conn
}
