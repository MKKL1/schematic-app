package user

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func RegisterRoutes(e *echo.Echo, server *Controller) {
	authMiddleware := auth.GetAuthMiddleware()

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1/users")
	v1Group.GET("/me", server.GetMe, authMiddleware)
	v1Group.GET("/:id", server.GetUserByID)
	v1Group.POST("/", server.CreateUser, authMiddleware)
}

type Controller struct {
	userApp client.PostApplication
}

func NewController(userApp client.PostApplication) *Controller {
	return &Controller{userApp}
}

func (s *Controller) GetMe(c echo.Context) error {
	subjectUUID, err := auth.ExtractOidcSub(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	userDto, err := s.userApp.Query.GetUserBySub(ctx, subjectUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ToResponse(userDto))
}

func (s *Controller) GetUserByID(c echo.Context) error {
	id, err := user.ParseUserID(c.Param("id"))
	if err != nil {
		return err
	}

	ctx := context.Background()
	userDto, err := s.userApp.Query.GetUserById(ctx, id.Unwrap())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ToResponse(userDto))
}

func (s *Controller) CreateUser(c echo.Context) error {
	subjectUUID, err := auth.ExtractOidcSub(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	requestData := CreateRequest{}
	err = json.NewDecoder(c.Request().Body).Decode(&requestData)
	if err != nil {
		return err
	}

	id, err := s.userApp.Command.CreateUser(ctx, requestData.Name, subjectUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": strconv.FormatInt(id, 10)})
}
