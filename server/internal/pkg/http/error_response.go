package http

type ErrorResponse struct {
	Status  int       `json:"status"`
	ID      string    `json:"id"`
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
}

type ErrorCode string

func (e *ErrorResponse) Error() string {
	return e.Message
}

func NewErrorResponse(code ErrorCode, status int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
		Status:  status,
	}
}
