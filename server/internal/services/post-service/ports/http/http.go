package http

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/command"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func RegisterRoutes(e *echo.Echo, server *PostController) {
	authMiddleware := auth.GetAuthMiddleware()

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1/posts")
	v1Group.GET("/:id", server.GetPost)
	v1Group.POST("/", server.CreatePost, authMiddleware)
}

type PostController struct {
	application app.Application
	validate    *validator.Validate
}

func NewPostController(application app.Application) *PostController {
	return &PostController{application, validator.New(validator.WithRequiredStructEnabled())}
}

func (pc *PostController) GetPost(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	ctx := context.Background()
	params := query.GetPostByIdParams{Id: id}
	postDto, err := pc.application.Queries.GetPostById.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, PostToResponse(postDto))
}

func (pc *PostController) CreatePost(c echo.Context) error {
	subjectUUID, err := auth.ExtractOidcSub(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	requestData := PostCreateRequest{}
	err = json.NewDecoder(c.Request().Body).Decode(&requestData)
	if err != nil {
		return err
	}

	if err := pc.validate.Struct(requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	params := command.CreatePostParams{
		PostCreateRequest: requestData,
		Sub:               subjectUUID,
	}

	id, err := pc.application.Commands.CreatePost.Handle(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": strconv.FormatInt(id, 10)})
}
