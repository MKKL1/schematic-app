package image

import "time"

type Model struct {
	FileHash  string
	ImageType string
	CreatedAt time.Time
}
