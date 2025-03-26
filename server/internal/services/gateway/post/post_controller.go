package post

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client/post"
	gtHttp "github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

func RegisterRoutes(e *echo.Echo, server *Controller) {
	authMiddleware := auth.GetAuthMiddleware()

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1/posts")
	v1Group.GET("/:id", server.GetPost)
	v1Group.POST("/", server.CreatePost, authMiddleware)
}

type Controller struct {
	validate    *validator.Validate
	postService post.Service
}

func NewController(postApp post.Service) *Controller {
	return &Controller{validator.New(validator.WithRequiredStructEnabled()), postApp}
}

func (pc *Controller) GetPost(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	ctx := context.Background()
	postDto, err := pc.postService.GetPostById(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, PostToResponse(postDto))
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
		return gtHttp.MapValidationError(err, requestData)
	}

	categParams := make([]post.CreatePostRequestCategory, len(requestData.Categories))
	for i, cat := range requestData.Categories {
		categParams[i] = post.CreatePostRequestCategory{
			Name:     cat.Name,
			Metadata: cat.Metadata,
		}
	}

	filesParams := make([]uuid.UUID, len(requestData.Files))
	for i, f := range requestData.Files {
		fId, err := uuid.Parse(f)
		if err != nil {
			return err
		}
		filesParams[i] = fId
	}

	var authorId *int64
	if requestData.Author != nil {
		_authorId, err := strconv.ParseInt(*requestData.Author, 10, 64)
		if err != nil {
			return err
		}
		authorId = &_authorId
	}

	params := post.CreatePostRequest{
		Name:        requestData.Name,
		Description: requestData.Description,
		AuthorID:    authorId,
		Sub:         subjectUUID,
		Categories:  categParams,
		Tags:        requestData.Tags,
		Files:       filesParams,
	}

	id, err := pc.postService.CreatePost(ctx, params)
	if err != nil {
		return HandlePostCreateErr(err, requestData)
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": strconv.FormatInt(id, 10)})
}

// mapFieldPath converts a field string in the format "categoryName:fieldName"
// into "categories[i].metadata.fieldName", where i is the index of the category with that name.
func mapFieldPath(field string, req PostCreateRequest) string {
	parts := strings.Split(field, ":")
	if len(parts) != 2 {
		// Fallback to original if format is unexpected.
		return field
	}
	categoryName, fieldName := parts[0], parts[1]
	for i, cat := range req.Categories {
		if cat.Name == categoryName {
			return fmt.Sprintf("categories[%d].metadata.%s", i, fieldName)
		}
	}
	return field
}
