package file

import (
	"time"
)

type TempFileCreated struct {
	Key        string
	Expiration time.Duration
	Url        string
}
