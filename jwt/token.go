package jwt

import (
	_ "embed"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/oomph-ac/api/endpoint/types"
	"github.com/oomph-ac/api/errors"
)

var (
	//go:embed JWT_MASTER_KEY
	secretKey   []byte
	oldJWTError = errors.New(errors.APIUserFault, "old JWT version", nil)
)

type AuthClaims struct {
	IPAddress  string `json:"addr"`
	OomphKey   string `json:"oomph_key"`
	Admin      bool   `json:"admin"`
	Expiration int64  `json:"expiresAt"`
}

func (ac *AuthClaims) Valid() error {
	if ac.IPAddress == "" || ac.OomphKey == "" || ac.Expiration == 0 {
		return fmt.Errorf("old version ")
	}
	return nil
}

// NewAuthToken creates a new JWT token that other endpoints in the API may require. It includes
// claims such as the authentication key, the IP address that the authentication request originated
// from, a time where the Oomph client should attempt to refresh the token, and a time the token will
// no longer be valid.
func NewAuthToken(dat types.DBAuthData, addr string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"addr":      addr,
		"oomph_key": dat.Key,
		"admin":     dat.Admin,
		"expiresAt": time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateAuthToken returns true if the given token has been signed by the API.
func ValidateAuthToken(token, ipAddr string) (*AuthClaims, *errors.APIError) {
	// Parse the JWT token we recieved with our secret key. This will return an error if the token is missing claims.
	tk, jwtErr := jwt.ParseWithClaims(token, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if jwtErr != nil {
		// Here, we could not parse the key because the signature is invalid - bad actors are at work ://
		if jwtErr.Error() == jwt.ErrSignatureInvalid.Error() {
			return nil, errors.New(
				errors.APIUserFaultNeedsLog,
				"invalid JWT token",
				jwtErr,
			)
		} else if jwtErr == oldJWTError {
			// Here, the JWT token has a valid signature, however, it is missing fields and so
			// we can't accept it.
			return nil, oldJWTError
		}

		// This is most likely the result of a type error.
		return nil, errors.New(
			errors.APIServerFault,
			"server unable to validate token",
			jwtErr,
		)
	}

	// Make sure the authentication token is not expired.
	claims := tk.Claims.(*AuthClaims)
	if claims.Expiration <= time.Now().Unix() {
		return nil, errors.New(
			errors.APIUserFault,
			"authentication token expired",
			nil,
		)
	}

	// Match the IP address in the token's claims to the one we have from the client. If
	// they do not match, a token replay is being attempted.
	if claims.IPAddress != ipAddr {
		return nil, errors.New(
			errors.APIUserFaultNeedsLog,
			"detected token replay",
			nil,
		)
	}

	return claims, nil
}
