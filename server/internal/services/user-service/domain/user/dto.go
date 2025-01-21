package user

import (
	"github.com/google/uuid"
	"strconv"
)

type User struct {
	ID      UserID
	Name    string
	OidcSub uuid.UUID
}

type UserID int64

func (u UserID) Unwrap() int64 {
	return int64(u)
}

func (u UserID) String() string {
	return strconv.FormatInt(int64(u), 10)
}

func ParseUserID(str string) (UserID, error) {
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return UserID(0), err
	}
	return UserID(id), nil
}

func ToDTO(user Model) User {
	return User{
		ID:      UserID(user.ID),
		Name:    user.Name,
		OidcSub: user.OidcSub,
	}
}
