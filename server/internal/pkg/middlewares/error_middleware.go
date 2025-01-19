package middlewares

import (
	"errors"
	appErr "github.com/MKKL1/schematic-app/server/internal/pkg/error"
	appHttp "github.com/MKKL1/schematic-app/server/internal/pkg/http"
	"github.com/labstack/echo/v4"
	"net/http"
)

func HTTPErrorHandler(errorMapper func(err appErr.Error) error) func(err error, c echo.Context) {
	return func(err error, c echo.Context) {
		var serviceError *appErr.Error
		var errorResponse *appHttp.ErrorResponse
		var echoError *echo.HTTPError

		if errors.As(err, &echoError) {
			mappedError, ok := appHttp.MapEchoError(echoError)
			if ok {
				errorResponse = mappedError
			} else {
				errorResponse = appHttp.NewErrorResponse("ECHO_ERROR", echoError.Code, echoError.Message.(string))
			}
		} else if errors.As(err, &serviceError) {
			err = errorMapper(*serviceError)
		}

		if !errors.As(err, &errorResponse) || errorResponse == nil {
			//Handle generic error
			errorResponse = appHttp.NewErrorResponse("INTERNAL_SERVER_ERROR", http.StatusInternalServerError, "Internal Server Error")
		}

		errorResponse.ID = c.Response().Header().Get(echo.HeaderXRequestID)
		err = c.JSON(errorResponse.Status, errorResponse)
		if err != nil {
			return
		}
	}
}
