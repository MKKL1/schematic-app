package server

import (
	"context"
	"fmt"
	grpc2 "github.com/MKKL1/schematic-app/server/internal/pkg/grpc"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			l.Trace().Msg(msg)
		case logging.LevelInfo:
			l.Trace().Msg(msg)
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

func RunGRPCServer(ctx context.Context, addr string, errorMapper func(err error) error, registerServer func(server *grpc.Server)) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.DateTime,
	}

	logger := zerolog.New(consoleWriter).With().Timestamp().Logger()
	grpcLogger := logger.With().Str("component", "grpc-server").Logger()

	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add any other option (check functions starting with logging.With).
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc2.ErrorMappingUnaryInterceptor(errorMapper),
			logging.UnaryServerInterceptor(InterceptorLogger(grpcLogger), opts...),
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(grpcLogger), opts...),
		),
	)
	registerServer(grpcServer)

	go func() {
		listen, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Err(err)
		}
		grpcLogger.Info().Msg("Starting gRPC server")
		grpcLogger.Fatal().Err(grpcServer.Serve(listen))
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				grpcLogger.Info().Msg("Shutting down gRPC server")
				grpcServer.GracefulStop()
				grpcLogger.Info().Msg("gRPC server shut down")
				return
			}
		}
	}()
}
