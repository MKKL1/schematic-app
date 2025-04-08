package metrics

import (
	"github.com/rs/zerolog"
	"time"
)

type LogClient struct {
	Logger zerolog.Logger
}

func (l LogClient) RecordCommandDuration(commandName string, success bool, duration time.Duration) {
	l.Logger.Info().Str("component", "metrics-log").Str("command", commandName).Bool("success", success).Str("duration", duration.String()).Msg("Metrics log")
}

func (l LogClient) RecordQueryDuration(queryName string, success bool, duration time.Duration) {
	l.Logger.Info().Str("component", "metrics-log").Str("query", queryName).Bool("success", success).Str("duration", duration.String()).Msg("Metrics log")
}

var _ Client = (*LogClient)(nil)
