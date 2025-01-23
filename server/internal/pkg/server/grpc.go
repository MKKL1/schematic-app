package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"os"
)

func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg)
		case logging.LevelInfo:
			l.Info().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func RunGRPCServer(ctx context.Context, addr string, registerServer func(server *grpc.Server)) {
	logger := zerolog.New(os.Stdout).With().Logger()
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(logger), opts...),
		),
	)
	registerServer(grpcServer)

	go func() {
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			log.Err(err)
		}
		log.Info().Str("addr", addr).Msg("starting server")
		log.Fatal().Err(grpcServer.Serve(listen))
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info().Str("port", "8001").Msg("shutting down gRPC server")
				grpcServer.GracefulStop()
				log.Info().Msg("server shut down")
				return
			}
		}
	}()
}
