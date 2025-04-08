package decorator

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

type commandLoggingDecorator[C any, R any] struct {
	base   CommandHandler[C, R]
	logger zerolog.Logger
}

func (d commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (res R, err error) {
	logger := d.logger.With().
		Str("command", generateActionName(cmd)).
		Str("command_body", fmt.Sprintf("%#v", cmd)).
		Logger()

	logger.Debug().Msg("Executing command")
	defer func() {
		if err == nil {
			logger.Info().Msg("Command executed successfully")
		} else {
			logger.Error().Err(err).Msg("Failed to execute command")
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger zerolog.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := d.logger.With().
		Str("query", generateActionName(cmd)).
		Str("query_body", fmt.Sprintf("%#v", cmd)).
		Logger()

	logger.Debug().Msg("Executing query")
	defer func() {
		if err == nil {
			logger.Info().Msg("Query executed successfully")
		} else {
			logger.Error().Err(err).Msg("Failed to execute query")
		}
	}()

	return d.base.Handle(ctx, cmd)
}
