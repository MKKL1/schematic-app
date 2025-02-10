package http

import "google.golang.org/genproto/googleapis/rpc/errdetails"

type ErrorMessageMapper func(*errdetails.ErrorInfo) ErrorDetail

type DefaultErrorMessageMapper struct {
	mapper map[string]ErrorMessageMapper
}

func NewDefaultErrorMessageMapper() *DefaultErrorMessageMapper {
	return &DefaultErrorMessageMapper{mapper: make(map[string]ErrorMessageMapper)}
}

func (m *DefaultErrorMessageMapper) MapError(errInfo *errdetails.ErrorInfo) ErrorDetail {
	return m.mapper[errInfo.Reason](errInfo)
}

func (m *DefaultErrorMessageMapper) AddMapper(reason string, mapper ErrorMessageMapper) {
	m.mapper[reason] = mapper
}
