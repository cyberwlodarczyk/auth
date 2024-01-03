package argon2id

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"runtime"
	"strings"

	"golang.org/x/crypto/argon2"
)

var Version = argon2.Version

var (
	ErrInvalidHash         = errors.New("argon2id: hash is not in the correct format")
	ErrIncompatibleVariant = errors.New("argon2id: incompatible variant of argon2")
	ErrIncompatibleVersion = errors.New("argon2id: incompatible version of argon2")
)

var DefaultParams = &Params{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

type Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func Key(params *Params, salt, password []byte) []byte {
	return argon2.IDKey(
		password,
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)
}

func Encode(params *Params, salt, key []byte) string {
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(Key(params, salt, key)),
	)
}

func Decode(hash string) (params *Params, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 || vals[0] != "" {
		return nil, nil, nil, ErrInvalidHash
	}
	if vals[1] != "argon2id" {
		return nil, nil, nil, ErrIncompatibleVariant
	}
	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}
	params = &Params{}
	_, err = fmt.Sscanf(
		vals[3],
		"m=%d,t=%d,p=%d",
		&params.Memory,
		&params.Iterations,
		&params.Parallelism,
	)
	if err != nil {
		return nil, nil, nil, err
	}
	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	params.SaltLength = uint32(len(salt))
	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	params.KeyLength = uint32(len(key))
	return params, salt, key, nil
}

func RandomSalt(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

func Hash(password []byte, params *Params) (hash string, err error) {
	salt, err := RandomSalt(params.SaltLength)
	if err != nil {
		return "", err
	}
	return Encode(params, salt, password), nil
}

func Compare(password []byte, hash string) (match bool, params *Params, err error) {
	params, salt, key, err := Decode(hash)
	if err != nil {
		return false, nil, err
	}
	otherKey := Key(params, salt, password)
	if subtle.ConstantTimeEq(int32(len(key)), int32(len(otherKey))) == 0 {
		return false, params, nil
	}
	return subtle.ConstantTimeCompare(key, otherKey) == 1, params, nil
}
