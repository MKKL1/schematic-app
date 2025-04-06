package query

//type GetImageParams struct {
//	ImageID string
//}
//type QualityInfo struct {
//	URL string `json:"url"`
//}
//type GetImageResult struct {
//}
//
//type GetImageHandler decorator.QueryHandler[GetImageParams, GetImageResult]
//
//type getImageHandler struct {
//	repo   file.ImageRepository
//	logger zerolog.Logger
//}
//
//func NewGetImageHandler(repo file.ImageRepository, logger zerolog.Logger /*, cfg *SomeConfig */) GetImageHandler {
//	return decorator.ApplyQueryDecorators[GetImageParams, GetImageResult](
//		getImageHandler{repo: repo, logger: logger /* config: cfg*/},
//	)
//}
//
//func (h getImageHandler) Handle(ctx context.Context, query GetImageParams) (GetImageResult, error) {
//
//	return result, nil
//}
