package http

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/dto"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserController struct {
	service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	if service == nil {
		panic("service must not be nil")
	}
	return &UserController{
		service: service,
	}
}

func (s *UserController) GetMe(c echo.Context) error {
	subjectUUID, err := auth.ExtractOidcSub(c)
	if err != nil {
		return err
	}

	ctx := context.Background()
	user, err := s.service.GetUserByOidcSub(ctx, subjectUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UserToResponse(user))
}

func (s *UserController) GetUserByID(c echo.Context) error {
	id, err := dto.ParseUserID(c.Param("id"))
	if err != nil {
		return err
	}

	ctx := context.Background()
	user, err := s.service.GetUserById(ctx, id)
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

	createdUser, err := s.service.CreateUser(ctx, requestData.Name, subjectUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, createdUser)
}
