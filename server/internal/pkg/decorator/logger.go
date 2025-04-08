package decorator

import (
	"fmt"
	"github.com/rs/zerolog"
)

func AddCmdInfo(handler any, logger zerolog.Logger) zerolog.Logger {
	return logger.With().Str("command", generateActionName(handler)).Str("command_body", fmt.Sprintf("%#v", handler)).Logger()
}

func AddQueryInfo(handler any, logger zerolog.Logger) zerolog.Logger {
	return logger.With().Str("query", generateActionName(handler)).Str("query_body", fmt.Sprintf("%#v", handler)).Logger()
}
