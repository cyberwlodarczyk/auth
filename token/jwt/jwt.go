package jwt

import (
	"errors"
	"time"

	"github.com/cyberwlodarczyk/auth/token"
	"github.com/golang-jwt/jwt/v5"
)

const Algorithm = "HS256"

type claims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

func NewService[T any](key []byte, age time.Duration) token.Service[T] {
	return &service[T]{key, age}
}

type service[T any] struct {
	key []byte
	age time.Duration
}

func (s *service[T]) Sign(data T) (t string, err error) {
	t, err = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&claims[T]{data, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.age)),
		}},
	).SignedString(s.key)
	return
}

func (s *service[T]) Verify(t string) (data T, err error) {
	c := &claims[T]{}
	_, err = jwt.ParseWithClaims(
		t,
		c,
		func(t *jwt.Token) (interface{}, error) { return s.key, nil },
		jwt.WithLeeway(token.Leeway),
		jwt.WithValidMethods([]string{Algorithm}),
		jwt.WithExpirationRequired(),
		jwt.WithStrictDecoding(),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			err = token.ErrInvalidFormat
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			err = token.ErrInvalidSignature
		case errors.Is(err, jwt.ErrTokenRequiredClaimMissing):
			err = token.ErrMissingExpiration
		case errors.Is(err, jwt.ErrTokenExpired):
			err = token.ErrExceededExpiration
		}
	} else {
		data = c.Data
	}
	return
}
