package validation

import (
	"regexp"
	"unicode"
	"unicode/utf8"
)

var (
	DefaultPasswordConfig = &PasswordConfig{
		Upper:     1,
		Lower:     1,
		Number:    1,
		Special:   1,
		MinLength: 12,
		MaxLength: 64,
	}
	DefaultEmailPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type Service[T any] interface {
	Check(T) bool
}

type PasswordConfig struct {
	Upper     int
	Lower     int
	Number    int
	Special   int
	MinLength int
	MaxLength int
}

func NewPasswordService(cfg *PasswordConfig) Service[[]byte] {
	return &passwordService{cfg}
}

type passwordService struct {
	cfg *PasswordConfig
}

func (s *passwordService) Check(password []byte) bool {
	var (
		r       rune
		size    int
		upper   int
		lower   int
		number  int
		special int
	)
	for i := 0; i < len(password); {
		r, size = utf8.DecodeRune(password[i:])
		switch {
		case unicode.IsUpper(r):
			upper++
		case unicode.IsLower(r):
			lower++
		case unicode.IsNumber(r):
			number++
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			special++
		}
		i += size
	}
	return upper >= s.cfg.Upper &&
		lower >= s.cfg.Lower &&
		number >= s.cfg.Number &&
		special >= s.cfg.Special &&
		len(password) >= s.cfg.MinLength &&
		len(password) <= s.cfg.MaxLength
}

func NewEmailService(pattern *regexp.Regexp) Service[string] {
	return &emailService{pattern}
}

type emailService struct {
	pattern *regexp.Regexp
}

func (s *emailService) Check(email string) bool {
	return s.pattern.MatchString(email)
}

func NewMinMaxService(min, max int) Service[string] {
	return &minMaxService{min, max}
}

type minMaxService struct {
	min int
	max int
}

func (s *minMaxService) Check(value string) bool {
	return len(value) >= s.min && len(value) <= s.max
}
