package decorator

import (
	"context"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/rs/zerolog"
	"strings"
)

type CommandHandler[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func ApplyCommandDecorators[H any, R any](handler CommandHandler[H, R], logger zerolog.Logger, metricsClient metrics.Client) CommandHandler[H, R] {
	return commandLoggingDecorator[H, R]{
		base: commandMetricsDecorator[H, R]{
			base: commandErrorDecorator[H, R]{
				base: handler,
			},
			client: metricsClient,
		},
		logger: logger,
	}
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
