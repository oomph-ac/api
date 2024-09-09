package jwt

import (
	_ "embed"
	"time"

	"github.com/golang-jwt/jwt"
)

//go:embed JWT_MASTER_KEY
var secretKey []byte

// NewAuthToken creates a new JWT token that other endpoints in the API may require. It includes
// claims such as the authentication key, the IP address that the authentication request originated
// from, a time where the Oomph client should attempt to refresh the token, and a time the token will
// no longer be valid.
func NewAuthToken(key, addr string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"addr": addr,
		"key":  key,
		"exp":  time.Now().Add(time.Hour),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
