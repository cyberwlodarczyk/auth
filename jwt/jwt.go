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
	ErrInvalidFormat      = errors.New("jwt: invalid format")
	ErrInvalidSignature   = errors.New("jwt: invalid signature")
	ErrMissingExpiration  = errors.New("jwt: missing expiration")
	ErrExceededExpiration = errors.New("jwt: exceeded expiration")
)

type claims[T any] struct {
	Data T `json:"data"`
	jwt.RegisteredClaims
}

type Service[T any] interface {
	Sign(T) (string, error)
	Verify(string) (T, error)
}

func NewService[T any](key []byte, age time.Duration) Service[T] {
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
		jwt.WithLeeway(Leeway),
		jwt.WithValidMethods([]string{Algorithm}),
		jwt.WithExpirationRequired(),
		jwt.WithStrictDecoding(),
	)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			err = ErrInvalidFormat
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			err = ErrInvalidSignature
		case errors.Is(err, jwt.ErrTokenRequiredClaimMissing):
			err = ErrMissingExpiration
		case errors.Is(err, jwt.ErrTokenExpired):
			err = ErrExceededExpiration
		}
	} else {
		data = c.Data
	}
	return
}
