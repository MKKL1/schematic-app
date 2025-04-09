package common

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/closer"
	"github.com/MKKL1/schematic-app/server/internal/pkg/config"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DefaultShutdownTimeout = 15 * time.Second
	DefaultInitTimeout     = 3 * time.Minute
)

// Runtime provides a common structure for application lifecycle management.
type Runtime[C any] struct {
	rootCtx         context.Context
	errGroup        *errgroup.Group
	groupCtx        context.Context
	Logger          zerolog.Logger
	Closer          *closer.Closer
	Config          *C // Store the loaded config
	ShutdownTimeout time.Duration
	cancelRootCtx   context.CancelFunc // To call on fatal errors before setup finishes
}

// NewRuntime initializes the common application runtime components.
// It takes the service name (for logging), config file path, and a pointer to the config struct type.
func NewRuntime[C any](serviceName string) (*Runtime[C], error) {
	baseCtx := context.Background()
	rootCtx, stop := signal.NotifyContext(baseCtx, os.Interrupt, syscall.SIGTERM)

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	logger := zerolog.New(consoleWriter).With().Timestamp().Str("service", serviceName).Logger()
	logger.Info().Msg("Initializing application runtime...")

	configPath := flag.String("config", "config.yaml", "Path to configuration YAML file")
	flag.Parse()
	cfg, err := config.LoadConfig[C](*configPath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to load configuration")
		stop()
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	logger.Info().Msg("Configuration loaded successfully")

	appCloser := closer.New(logger)
	g, groupCtx := errgroup.WithContext(rootCtx)

	rt := &Runtime[C]{
		rootCtx:         rootCtx,
		errGroup:        g,
		groupCtx:        groupCtx,
		Logger:          logger,
		Closer:          appCloser,
		Config:          cfg,
		ShutdownTimeout: DefaultShutdownTimeout,
		cancelRootCtx:   stop,
	}

	return rt, nil
}

// InitCtx returns a context typically used for initializing dependencies.
func (r *Runtime[C]) InitCtx(timeout time.Duration) (context.Context, context.CancelFunc) {
	if timeout == 0 {
		timeout = DefaultInitTimeout
	}
	return context.WithTimeout(r.rootCtx, timeout)
}

// Go runs a function within the runtime's error group.
func (r *Runtime[C]) Go(f func(ctx context.Context) error) {
	r.errGroup.Go(func() error {
		err := f(r.groupCtx)
		if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, syscall.EPIPE) {
			r.Logger.Error().Err(err).Msg("Error in managed goroutine")
			return err
		}
		return nil
	})
}

// AddCloserFunc adds a custom function to the closer.
func (r *Runtime[C]) AddCloserFunc(f func(ctx context.Context) error) {
	r.Closer.Add(f)
}

// Run blocks until the application is signaled to shut down or a critical error occurs.
func (r *Runtime[C]) Run() {
	r.Logger.Info().Msg("Application runtime starting. Waiting for exit signal or error...")

	err := r.errGroup.Wait()

	r.Logger.Info().Msg("Shutdown signal received or error occurred. Initiating graceful shutdown...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), r.ShutdownTimeout)
	defer shutdownCancel()

	r.Closer.CloseAll(shutdownCtx)

	if err != nil && !errors.Is(err, context.Canceled) {
		r.Logger.Error().Err(err).Msg("Application shutdown finished with error.")
		r.cancelRootCtx()
		os.Exit(1)
	} else {
		r.Logger.Info().Msg("Application shutdown finished gracefully.")
		r.cancelRootCtx()
	}
}
