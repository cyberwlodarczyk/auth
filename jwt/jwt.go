package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Leeway    = 3 * time.Minute
	Algorithm = "HS256"
)

var (
	ErrExpiredToken          = errors.New("jwt: token is expired")
	ErrMalformedToken        = errors.New("jwt: token is not in the correct format")
	ErrInvalidSignature      = errors.New("jwt: token signature is invalid")
	ErrMissingExpiresAtClaim = errors.New("jwt: token expiration time claim is not present")
)

type claims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

func Sign[T any](data T, key []byte, age time.Duration) (token string, err error) {
	return jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&claims[T]{data, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(age)),
		}},
	).SignedString(key)
}

func Verify[T any](token string, key []byte) (data T, expiresAt time.Time, err error) {
	c := &claims[T]{}
	_, err = jwt.ParseWithClaims(
		token,
		c,
		func(t *jwt.Token) (interface{}, error) { return key, nil },
		jwt.WithLeeway(Leeway),
		jwt.WithValidMethods([]string{Algorithm}),
		jwt.WithExpirationRequired(),
		jwt.WithStrictDecoding(),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			err = ErrMalformedToken
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			err = ErrInvalidSignature
		case errors.Is(err, jwt.ErrTokenRequiredClaimMissing):
			err = ErrMissingExpiresAtClaim
		case errors.Is(err, jwt.ErrTokenExpired):
			err = ErrExpiredToken
		}
	} else {
		data, expiresAt = c.Data, c.ExpiresAt.Time
	}
	return
}
