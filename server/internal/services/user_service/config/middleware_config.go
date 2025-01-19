package config

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/middlewares"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigMiddlewares(e *echo.Echo) {

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = middlewares.HTTPErrorHandler(http.MapAppError)
	e.Use(middleware.RequestID())
}
