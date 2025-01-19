package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ExtractOidcSub(c echo.Context) (uuid.UUID, error) {
	subjectStr, err := c.Get("user").(*jwt.Token).Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	subjectUUID, err := uuid.Parse(subjectStr)
	if err != nil {
		return uuid.UUID{}, err
	}

	return subjectUUID, nil
}
