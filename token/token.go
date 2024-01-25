package token

import (
	"errors"
	"time"
)

const Leeway = 3 * time.Minute

var (
	ErrInvalidFormat      = errors.New("token: invalid format")
	ErrInvalidSignature   = errors.New("token: invalid signature")
	ErrMissingExpiration  = errors.New("token: missing expiration")
	ErrExceededExpiration = errors.New("token: exceeded expiration")
)

type Service[T any] interface {
	Sign(T) (string, error)
	Verify(string) (T, error)
}
