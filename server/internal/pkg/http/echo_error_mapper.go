package http

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type echoErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var echoErrMapper = map[error]echoErr{
	jwt.ErrTokenExpired:     {http.StatusUnauthorized, "Token expired"},
	jwt.ErrInvalidKey:       {http.StatusUnauthorized, "Invalid token signature"},
	jwt.ErrTokenNotValidYet: {http.StatusUnauthorized, "Token not valid yet"},
	jwt.ErrTokenMalformed:   {http.StatusUnauthorized, "Malformed token"},
}

func MapEchoError(err error) (*ErrorResponse, bool) {
	for knownErr, mapped := range echoErrMapper {
		if errors.Is(err, knownErr) {
			return NewErrorResponse("JWT_ERROR", mapped.Code, mapped.Message), true
		}
	}

	return nil, false
}
