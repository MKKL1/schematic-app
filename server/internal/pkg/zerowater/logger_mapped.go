package zerowater

//All it does it maps info level logs to debug
//Not a very good solution, but least complicated one

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/rs/zerolog"
)

type ZerologLoggerAdapterMapped struct {
	logger zerolog.Logger
}

// Logs an error message.
func (loggerAdapter *ZerologLoggerAdapterMapped) Error(msg string, err error, fields watermill.LogFields) {
	event := loggerAdapter.logger.Err(err)

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// Logs an info message.
func (loggerAdapter *ZerologLoggerAdapterMapped) Info(msg string, fields watermill.LogFields) {
	event := loggerAdapter.logger.Debug()

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// Logs a debug message.
func (loggerAdapter *ZerologLoggerAdapterMapped) Debug(msg string, fields watermill.LogFields) {
	event := loggerAdapter.logger.Debug()

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// Logs a trace.
func (loggerAdapter *ZerologLoggerAdapterMapped) Trace(msg string, fields watermill.LogFields) {
	event := loggerAdapter.logger.Trace()

	if fields != nil {
		addWatermillFieldsData(event, fields)
	}

	event.Msg(msg)
}

// Creates new adapter wiht the input fields as context.
func (loggerAdapter *ZerologLoggerAdapterMapped) With(fields watermill.LogFields) watermill.LoggerAdapter {
	if fields == nil {
		return loggerAdapter
	}

	subLog := loggerAdapter.logger.With()

	for i, v := range fields {
		subLog = subLog.Interface(i, v)
	}

	return &ZerologLoggerAdapterMapped{
		logger: subLog.Logger(),
	}
}

// Gets a new zerolog adapter for use in the watermill context.
func NewZerologLoggerAdapterMapped(logger zerolog.Logger) *ZerologLoggerAdapterMapped {
	return &ZerologLoggerAdapterMapped{
		logger: logger,
	}
}
