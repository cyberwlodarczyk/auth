package password

import "errors"

var (
	ErrInvalidFormat       = errors.New("password: invalid format")
	ErrIncompatibleVariant = errors.New("password: incompatible variant")
	ErrIncompatibleVersion = errors.New("password: incompatible version")
)

type Service interface {
	Hash([]byte) (string, error)
	Compare([]byte, string) (bool, bool, error)
}
