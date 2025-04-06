package file

import (
	"time"
)

type Type string

const (
	Image     Type = "image"
	Schematic Type = "minecraft-schematic"
	Unknown   Type = "unknown"
)

type TempFileCreated struct {
	Key        string
	Expiration time.Time
}
