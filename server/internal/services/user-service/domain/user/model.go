package user

import (
	"github.com/google/uuid"
	"strconv"
)

type User struct {
	ID      ID
	Name    string
	OidcSub uuid.UUID
}

type ID int64

func (u ID) Unwrap() int64 {
	return int64(u)
}

func (u ID) String() string {
	return strconv.FormatInt(int64(u), 10)
}

func ParseUserID(str string) (ID, error) {
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return ID(0), err
	}
	return ID(id), nil
}
