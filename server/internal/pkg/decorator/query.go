package decorator

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"github.com/rs/zerolog"
)

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], logger zerolog.Logger, metricsClient metrics.Client) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base: queryErrorDecorator[H, R]{
				base: handler,
			},
			client: metricsClient,
		},
		logger: logger,
	}
}
