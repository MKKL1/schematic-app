package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection" // Optional: For gRPC reflection
)

// Config for gRPC server
type Config struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}

// GetAddr returns the network address string
func (c Config) GetAddr() string {
	// Handle cases where host might already include port
	if net.ParseIP(c.Host) != nil || c.Host == "localhost" || c.Host == "" { // Simple check if it's likely an IP or hostname
		return net.JoinHostPort(c.Host, c.Port)
	}
	return c.Host // Assume host already contains port
}

// InterceptorLogger adapts zerolog to grpc-ecosystem logger
func InterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l = l.With().Fields(map[string]interface{}{"grpc.fields": fields}).Logger() // Adjust field logging if needed

		switch lvl {
		case logging.LevelDebug:
			l.Debug().Msg(msg) // Changed Debug to Debug level
		case logging.LevelInfo:
			l.Info().Msg(msg) // Changed Trace to Info level
		case logging.LevelWarn:
			l.Warn().Msg(msg)
		case logging.LevelError:
			l.Error().Msg(msg)
		default:
			l.Error().Msgf("unknown logging level %v: %s", lvl, msg) // Log as error instead of panic
		}
	})
}

// Run starts the gRPC server and handles graceful shutdown based on the context.
func Run(ctx context.Context, conf Config, logger zerolog.Logger, errorMapper func(err error) error, registerServer func(server *grpc.Server)) error {
	addr := conf.GetAddr()
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Error().Err(err).Str("address", addr).Msg("Failed to listen on gRPC address")
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	grpcLogger := logger.With().Str("component", "grpc-server").Str("address", addr).Logger()

	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
		// Add other options as needed
	}

	// Setup server with interceptors
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(grpcLogger), logOpts...),
			ErrorMappingUnaryInterceptor(errorMapper), // Ensure this exists in pkg/grpc
			// Add recovery interceptor here if desired
		),
		grpc.ChainStreamInterceptor(
			logging.StreamServerInterceptor(InterceptorLogger(grpcLogger), logOpts...),
			// Add stream error mapping/recovery if needed
		),
	)

	// Register application-specific services
	registerServer(grpcServer)
	reflection.Register(grpcServer) // Optional: Enable server reflection

	// Goroutine to start the server
	serverErrChan := make(chan error, 1)
	go func() {
		grpcLogger.Info().Msg("Starting gRPC server")
		if err := grpcServer.Serve(listen); err != nil && err != grpc.ErrServerStopped {
			grpcLogger.Error().Err(err).Msg("gRPC server failed")
			serverErrChan <- err
		} else {
			grpcLogger.Info().Msg("gRPC server stopped serving")
			close(serverErrChan) // Signal clean stop
		}
	}()

	// Goroutine to handle graceful shutdown
	shutdownComplete := make(chan struct{})
	go func() {
		defer close(shutdownComplete)
		<-ctx.Done() // Wait for context cancellation (e.g., SIGINT)
		grpcLogger.Info().Msg("Shutting down gRPC server gracefully...")

		// Perform graceful stop
		grpcServer.GracefulStop()

		grpcLogger.Info().Msg("gRPC server shut down.")
	}()

	// Wait for either server error or shutdown completion
	select {
	case err := <-serverErrChan:
		// Server failed before graceful shutdown was initiated
		// Ensure shutdown goroutine can exit if it's still waiting on ctx.Done()
		// (though it likely won't if serverErrChan received an error)
		<-shutdownComplete // Wait for shutdown attempt (might be immediate if server failed)
		return err         // Return the server error
	case <-shutdownComplete:
		// Server shutdown completed gracefully after ctx.Done()
		// Check if there was a startup error captured just before shutdown
		if err := <-serverErrChan; err != nil {
			return err
		}
		return nil // Graceful shutdown successful
	}
	// NOTE: Removed the internal shutdown goroutine from the original code.
	// Shutdown is now triggered by the parent context cancellation.
}
