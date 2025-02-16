package post

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/pkg/auth"
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"github.com/MKKL1/schematic-app/server/internal/services/gateway/grpc"
	gtHttp "github.com/MKKL1/schematic-app/server/internal/services/gateway/http"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
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
	validate *validator.Validate
	postApp  client.PostApplication
}

func NewController(postApp client.PostApplication) *Controller {
	return &Controller{validator.New(validator.WithRequiredStructEnabled()), postApp}
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

	categParams := make([]client.CreateCategoryMetadataParams, len(requestData.Categories))
	for i, c := range requestData.Categories {
		categParams[i] = client.CreateCategoryMetadataParams{
			Name:     c.Name,
			Metadata: c.Metadata,
		}
	}

	params := client.CreatePostParams{
		Name:        requestData.Name,
		Description: requestData.Description,
		AuthorID:    requestData.Author,
		Sub:         subjectUUID,
		Categories:  categParams,
		Tags:        requestData.Tags,
	}

	id, err := pc.postApp.Command.CreatePost(ctx, params)
	if err != nil {
		//TODO move to other function
		st, ok := status.FromError(err)
		if ok {
			errInfo, found := grpc.GetMessage[errdetails.ErrorInfo](st.Details())
			if !found {
				return err
			}

			if errInfo.GetReason() != "POST_METADATA_VALIDATION_ERROR" {
				return err
			}

			badRequest, found := grpc.GetMessage[errdetails.BadRequest](st.Details())
			if !found {
				return err
			}

			var errDetails []gtHttp.ErrorDetail
			for _, v := range badRequest.GetFieldViolations() {
				parameter := mapFieldPath(v.GetField(), requestData)
				errDetails = append(errDetails, gtHttp.ValidationErrorBuilder{
					Parameter: parameter,
					Detail:    v.GetReason(),
					Message:   v.GetDescription(),
				}.Build())
			}

			return &gtHttp.GatewayError{
				HttpCode: http.StatusBadRequest,
				ErrResponse: gtHttp.ErrorResponse{
					Errors: errDetails,
				},
			}
		}

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
