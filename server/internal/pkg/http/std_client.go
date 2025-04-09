package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Host            string        `koanf:"host"`
	Port            string        `koanf:"port"` // Optional: can combine with host if preferred
	ShutdownTimeout time.Duration `koanf:"shutdown_timeout"`
}

// GetAddr returns the network address string
func (c Config) GetAddr() string {
	// Handle cases where host might already include port
	if net.ParseIP(c.Host) != nil || c.Host == "localhost" || c.Host == "" { // Simple check if it's likely an IP or hostname
		return net.JoinHostPort(c.Host, c.Port)
	}
	return c.Host
}

// Run starts the HTTP server and handles graceful shutdown.
func Run(ctx context.Context, conf Config, logger zerolog.Logger, handler http.Handler) error {
	addr := conf.GetAddr()
	httpLogger := logger.With().Str("component", "http-server").Str("address", addr).Logger()

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrChan := make(chan error, 1)

	go func() {
		httpLogger.Info().Msg("Starting HTTP server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			httpLogger.Error().Err(err).Msg("HTTP server ListenAndServe failed")
			serverErrChan <- err
		} else {
			httpLogger.Info().Msg("HTTP server stopped serving")
			close(serverErrChan)
		}
	}()

	shutdownComplete := make(chan struct{})
	go func() {
		defer close(shutdownComplete)
		<-ctx.Done()
		httpLogger.Info().Msg("Shutting down HTTP server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			httpLogger.Error().Err(err).Msg("HTTP server graceful shutdown failed")
		} else {
			httpLogger.Info().Msg("HTTP server shut down.")
		}
	}()

	select {
	case err := <-serverErrChan:
		<-shutdownComplete
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}
		return nil
	case <-shutdownComplete:
		if err := <-serverErrChan; err != nil {
			return fmt.Errorf("http server shutdown initiated: %w", err)
		}
		return nil
	}
}
