package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/lestrrat-go/jwx/jwk"
	"sync"
)

var (
	jwkSet jwk.Set
	once   sync.Once
)

func fetchJWKS() (jwk.Set, error) {
	var err error
	once.Do(func() {
		jwkSet, err = jwk.Fetch(context.Background(), "http://localhost:8082/realms/test-realm/protocol/openid-connect/certs")
	})
	return jwkSet, err
}

func GetAuthMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		KeyFunc: func(token *jwt.Token) (interface{}, error) {
			jwks, err := fetchJWKS()
			if err != nil {
				return nil, err
			}

			keyID := token.Header["kid"].(string)
			// Find the key by 'kid'
			key, ok := jwks.LookupKeyID(keyID)
			if !ok {
				return nil, fmt.Errorf("unable to find key %q", keyID)
			}

			var pubkey interface{}
			if err := key.Raw(&pubkey); err != nil {
				return nil, fmt.Errorf("unable to get the public key. Error: %s", err.Error())
			}

			return pubkey, nil
		},
	})
}
