package ports

import (
	"context"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/app/query"
	posthttp "github.com/MKKL1/schematic-app/server/internal/services/post-service/http"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func RegisterRoutes(e *echo.Echo, server *PostController) {
	//authMiddleware := auth.GetAuthMiddleware()

	apiGroup := e.Group("/api")
	v1Group := apiGroup.Group("/v1/posts")
	v1Group.GET("/:id", server.GetPost)
}

type PostController struct {
	application app.Application
}

func NewPostController(application app.Application) *PostController {
	return &PostController{application}
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

	return c.JSON(http.StatusOK, posthttp.PostToResponse(postDto))
}
