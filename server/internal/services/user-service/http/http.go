package http

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/app/query"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/labstack/echo/v4"
	"net/http"
)

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
	user, err := s.application.Queries.GetUserBySub.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UserToResponse(user))
}

func (s *UserController) GetUserByID(c echo.Context) error {
	id, err := user.ParseUserID(c.Param("id"))
	if err != nil {
		return err
	}

	ctx := context.Background()
	params := query.GetUserByIdParams{Id: id}
	user, err := s.application.Queries.GetUserById.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UserToResponse(user))
}

func (s *UserController) CreateUser(c echo.Context) error {
	subjectUUID, err := auth.ExtractOidcSub(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	requestData := UserCreateRequest{}
	err = json.NewDecoder(c.Request().Body).Decode(&requestData)
	if err != nil {
		return err
	}

	params := command.CreateUserParams{
		Username: requestData.Name,
		Sub:      subjectUUID,
	}
	err = s.application.Commands.CreateUser.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
