package mappers

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/genproto"
	"github.com/MKKL1/schematic-app/server/internal/services/image-service/app/query"
)

// AppToProtoGetImageSizesResponse maps the application result to the gRPC response.
func AppToProtoGetImageSizesResponse(result query.GetImageResult) (*genproto.GetImageSizesResponse, error) {
	var sizes []*genproto.ImageSize
	for _, size := range result.Sizes {
		sizes = append(sizes, &genproto.ImageSize{
			Url:        size.URL,
			PresetName: size.Preset.Name,
		})
	}
	return &genproto.GetImageSizesResponse{
		Sizes: sizes,
	}, nil
}
