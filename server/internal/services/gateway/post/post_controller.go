package post

import (
	"context"
	"encoding/json"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func RegisterRoutes(e *echo.Echo, server *Controller) {
	authMiddleware := auth.GetAuthMiddleware()

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1/posts")
	v1Group.GET("/:id", server.GetPost)
	v1Group.POST("/", server.CreatePost, authMiddleware)
}

type Controller struct {
	validate *validator.Validate
	postApp  client.PostApplication
	tagApp   client.TagApplication
}

func NewController(postApp client.PostApplication, tagApp client.TagApplication) *Controller {
	return &Controller{validator.New(validator.WithRequiredStructEnabled()), postApp, tagApp}
}

func (pc *Controller) GetPost(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	ctx := context.Background()
	postDto, err := pc.postApp.Query.GetPostById(ctx, id)
	if err != nil {
		return err
	}

	vars, err := pc.tagApp.Query.GetCategVarsByPost(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, PostToResponse(postDto, vars))
}

func (pc *Controller) CreatePost(c echo.Context) error {
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

	if err = pc.validate.Struct(requestData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	//Make sure that either name or id is set
	var authorName *string
	var authorId *int64

	if requestData.Author != nil {
		if requestData.Author.ID != nil {
			_authorId, err2 := strconv.ParseInt(*requestData.Author.ID, 10, 64)
			if err2 != nil {
				return err2
			}

			authorId = &_authorId
		} else {
			authorName = requestData.Author.Name
		}
	}
	params := client.CreatePostParams{
		Name:        requestData.Name,
		Description: requestData.Description,
		AuthorName:  authorName,
		AuthorID:    authorId,
		Sub:         subjectUUID,
	}

	id, err := pc.postApp.Command.CreatePost(ctx, params)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": strconv.FormatInt(id, 10)})
}
