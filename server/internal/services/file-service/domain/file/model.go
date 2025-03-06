package file

import (
	"time"
)

type TempFileCreated struct {
	Key        string
	Expiration time.Time
}
