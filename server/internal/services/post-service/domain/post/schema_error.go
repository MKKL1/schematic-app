package post

import (
	"fmt"
	"github.com/MKKL1/schematic-app/server/internal/services/post-service/domain/category/validator"
	"strings"
)

type PostMetadataError struct {
	Errors map[string]validator.ValidationError
}

func NewPostMetadataError() *PostMetadataError {
	return &PostMetadataError{}
}

func (p *PostMetadataError) Error() string {
	keys := make([]string, len(p.Errors))

	i := 0
	for k := range p.Errors {
		keys[i] = k
		i++
	}
	return fmt.Sprintf("errors in post metadata, includes categories: %s", strings.Join(keys, ", "))
}
