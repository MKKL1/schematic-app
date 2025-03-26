package http

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

type ErrorMessageMapper func(c echo.Context, errInfo *errdetails.ErrorInfo, details []any) (*GatewayError, bool)

type DefaultErrorMessageMapper struct {
	mapper map[string]ErrorMessageMapper
}

func NewDefaultErrorMessageMapper() *DefaultErrorMessageMapper {
	return &DefaultErrorMessageMapper{mapper: make(map[string]ErrorMessageMapper)}
}

func (m *DefaultErrorMessageMapper) MapError(c echo.Context, errInfo *errdetails.ErrorInfo, details []any) (*GatewayError, bool) {
	mapper, ok := m.mapper[errInfo.Reason]
	if !ok {
		return nil, false
	}
	return mapper(c, errInfo, details)
}

func (m *DefaultErrorMessageMapper) AddMapper(reason string, mapper ErrorMessageMapper) {
	m.mapper[reason] = mapper
}
