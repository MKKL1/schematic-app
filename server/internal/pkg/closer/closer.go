package closer

import (
	"context"
	"io"
	"sync"

	"github.com/rs/zerolog"
)

// Closer helps manage closing multiple resources gracefully.
type Closer struct {
	mu      sync.Mutex
	closers []func(ctx context.Context) error
	logger  zerolog.Logger
}

// New creates a new Closer.
func New(logger zerolog.Logger) *Closer {
	return &Closer{
		logger: logger.With().Str("component", "closer").Logger(),
	}
}

// Add registers a function to be called during shutdown.
// Functions are called in Last-In, First-Out (LIFO) order.
func (c *Closer) Add(f func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closers = append(c.closers, f)
}

// AddCloser registers an io.Closer to be closed during shutdown.
func (c *Closer) AddCloser(closer io.Closer, name string) {
	c.Add(func(ctx context.Context) error {
		c.logger.Info().Str("resource", name).Msg("Closing resource")
		err := closer.Close()
		if err != nil {
			c.logger.Error().Err(err).Str("resource", name).Msg("Failed to close resource")
		} else {
			c.logger.Info().Str("resource", name).Msg("Resource closed successfully")
		}
		// We typically don't propagate the error here to allow other closers to run
		return nil // Or return err if stopping on first error is desired
	})
}

// CloseAll calls all registered closing functions in reverse order.
// It waits for all of them to complete or the context to be cancelled.
func (c *Closer) CloseAll(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.logger.Info().Int("count", len(c.closers)).Msg("Starting graceful shutdown...")

	var wg sync.WaitGroup
	// Execute closers in LIFO order
	for i := len(c.closers) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(closeFunc func(ctx context.Context) error) {
			defer wg.Done()
			// We ignore the error here to ensure all closers run
			_ = closeFunc(ctx)
		}(c.closers[i])
	}

	// Wait for all closers to finish or context timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		c.logger.Info().Msg("All resources closed.")
	case <-ctx.Done():
		c.logger.Error().Err(ctx.Err()).Msg("Shutdown timed out.")
	}
}
