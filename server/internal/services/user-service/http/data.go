package http

import (
	"github.com/MKKL1/schematic-app/server/internal/services/user-service/domain/user"
)

type UserResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func UserToResponse(user user.User) UserResponse {
	return UserResponse{
		ID:   user.ID.String(),
		Name: user.Name,
	}
}

type UserCreateRequest struct {
	Name string `json:"name"`
}
