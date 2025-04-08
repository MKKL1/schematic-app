package decorator

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/pkg/metrics"
	"strings"
	"time"
)

type commandMetricsDecorator[C any, R any] struct {
	base   CommandHandler[C, R]
	client metrics.Client
}

func (d commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (res R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd)) // Ensure lower case and label-safe

	// Defer using the interface methods
	defer func() {
		duration := time.Since(start)
		success := err == nil

		// Use the interface method
		d.client.RecordCommandDuration(actionName, success, duration)

		// Optional: Call the increment method if you defined/implemented it
		// d.metricsClient.IncrementCommandTotal(actionName, success)
	}()

	return d.base.Handle(ctx, cmd)
}

type queryMetricsDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	client metrics.Client
}

func (d queryMetricsDecorator[C, R]) Handle(ctx context.Context, query C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(query)) // Ensure lower case and label-safe

	// Defer using the interface methods
	defer func() {
		duration := time.Since(start)
		success := err == nil

		// Use the interface method
		d.client.RecordQueryDuration(actionName, success, duration)

		// Optional: Call the increment method if you defined/implemented it
		// d.metricsClient.IncrementCommandTotal(actionName, success)
	}()

	return d.base.Handle(ctx, query)
}
