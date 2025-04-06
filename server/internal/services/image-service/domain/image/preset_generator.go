package image

import "fmt"

type ResizingType int

const (
	Fit  ResizingType = iota
	Fill ResizingType = iota
	Auto ResizingType = iota
)

type Preset struct {
	Name         string
	Width        int
	Height       int
	MaxBytes     *int
	ResizingType ResizingType
}

func GeneratePreset(imageType string) ([]Preset, error) {
	switch imageType {
	case "gallery":
		return generateGalleryPresets(), nil
	case "avatar":
		return generateUserAvatarPresets(), nil
	}

	return []Preset{}, fmt.Errorf("unsupported image type: %s", imageType)
}

//TODO move to config file

func generateGalleryPresets() []Preset {
	return []Preset{
		{
			Name:         "g_350",
			Width:        350,
			Height:       350,
			MaxBytes:     nil,
			ResizingType: Fill,
		},
	}
}

func generateUserAvatarPresets() []Preset {
	return []Preset{
		{
			Name:         "a_64",
			Width:        64,
			Height:       64,
			MaxBytes:     nil,
			ResizingType: Fit,
		},
		{
			Name:         "a_128",
			Width:        128,
			Height:       128,
			MaxBytes:     nil,
			ResizingType: Fit,
		},
	}
}
