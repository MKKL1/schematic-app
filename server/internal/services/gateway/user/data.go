package user

import (
	"github.com/MKKL1/schematic-app/server/internal/pkg/client"
	"strconv"
)

type Response struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ToResponse(user *client.User) Response {
	return Response{
		ID:   strconv.FormatInt(user.ID, 10),
		Name: user.Name,
	}
}

type CreateRequest struct {
	Name string `json:"name"`
}
