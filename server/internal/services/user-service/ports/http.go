package ports

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	userhttp "github.com/MKKL1/schematic-app/server/internal/services/user-service/http"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RegisterRoutes(e *echo.Echo, server *UserController) {
	authMiddleware := auth.GetAuthMiddleware()

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1/users")
	v1Group.GET("/me", server.GetMe, authMiddleware)
	v1Group.GET("/:id", server.GetUserByID)
	v1Group.POST("/", server.CreateUser, authMiddleware)
}

type UserController struct {
	application app.Application
}

func NewUserController(application app.Application) *UserController {
	return &UserController{application}
}

func (s *UserController) GetMe(c echo.Context) error {
	subjectUUID, err := auth.ExtractOidcSub(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	params := query.GetUserBySubParams{Sub: subjectUUID}
	userDto, err := s.application.Queries.GetUserBySub.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, userhttp.UserToResponse(userDto))
}

func (s *UserController) GetUserByID(c echo.Context) error {
	id, err := user.ParseUserID(c.Param("id"))
	if err != nil {
		return err
	}

	ctx := context.Background()
	params := query.GetUserByIdParams{Id: id}
	userDto, err := s.application.Queries.GetUserById.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, userhttp.UserToResponse(userDto))
}

func (s *UserController) CreateUser(c echo.Context) error {
	subjectUUID, err := auth.ExtractOidcSub(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	requestData := userhttp.UserCreateRequest{}
	err = json.NewDecoder(c.Request().Body).Decode(&requestData)
	if err != nil {
		return err
	}

	params := command.CreateUserParams{
		Username: requestData.Name,
		Sub:      subjectUUID,
	}
	id, err = s.application.Commands.CreateUser.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": id})
}
