package http

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, server *UserController) {
	authMiddleware := auth.GetAuthMiddleware()

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1/users")
	v1Group.GET("/me", server.GetMe, authMiddleware)
	v1Group.GET("/:id", server.GetUserByID)
	v1Group.POST("/", server.CreateUser, authMiddleware)
}
