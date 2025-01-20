package dto

import (
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/postgres/db"
	"github.com/google/uuid"
	"strconv"
)

type User struct {
	ID      UserID
	Name    string
	OidcSub uuid.UUID
}

type UserID int64

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

func ToDTO(user db.User) (User, error) {
	oidcSub, err := uuid.FromBytes(user.OidcSub.Bytes[:])
	if err != nil {
		return User{}, err
	}

	return User{
		ID:      UserID(user.ID),
		Name:    user.Name,
		OidcSub: oidcSub,
	}, nil
}
