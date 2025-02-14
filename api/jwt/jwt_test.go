package jwt

import (
	"regexp"
	"testing"
	"time"
)

const a1, a2 = 1 * time.Hour, 30 * time.Minute

var (
	k1, k2     = []byte("mysecr3tk3y"), []byte("ot#ersecr3t")
	d1, d2     = uid{1}, uid{2}
	s1, s2, s3 = NewService[uid](Config{k1, a1}), NewService[uid](Config{k2, a1}), NewService[uid](Config{k1, a2})
)

type uid struct {
	UserId int `json:"user_id"`
}

func TestSign(t *testing.T) {
	t1, err := s1.Sign(d1)
	if err != nil {
		t.Fatal(err)
	}
	t2, err := s1.Sign(d1)
	if err != nil {
		t.Fatal(err)
	}
	if t1 != t2 {
		t.Fatalf("tokens must be equal for the same data, key and age. got: %q and %q", t1, t2)
	}
	t3, err := s1.Sign(d2)
	if err != nil {
		t.Fatal(err)
	}
	if t1 == t3 {
		t.Fatalf("tokens cannot be equal for different data. got: %q", t1)
	}
	t4, err := s2.Sign(d1)
	if err != nil {
		t.Fatal(err)
	}
	if t1 == t4 {
		t.Fatalf("tokens cannot be equal for different keys. got: %q", t1)
	}
	t5, err := s3.Sign(d1)
	if err != nil {
		t.Fatal(err)
	}
	if t1 == t5 {
		t.Fatalf("tokens cannot be equal for different ages. got: %q", t1)
	}
	pattern, err := regexp.Compile(`^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]*$`)
	if err != nil {
		t.Fatal(err)
	}
	if !pattern.MatchString(t1) {
		t.Fatalf("token is not in the correct format: %q", t1)
	}
}

func TestVerify(t *testing.T) {
	s1 := NewService[uid](Config{k1, a1})
	t1, err := s1.Sign(d1)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		token string
		data  uid
	}{
		{t1, d1},
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InVzZXJfaWQiOjF9LCJleHAiOjMwMDAwMDAwMDB9.hcGNFzkw4fSrHmW0FZosWFdvYh33VpRP1dPxI8aYld4", d1},
	}
	for _, test := range tests {
		data, err := s1.Verify(test.token)
		if err != nil {
			t.Fatal(err)
		}
		if data != test.data {
			t.Fatalf("expected data: %v, got: %v", test.data, data)
		}
	}
	errors := []struct {
		token string
		err   error
	}{
		{"", ErrInvalidFormat},
		{"eyvGcOiJIUi4NiIInR5cCI6kpXVCJ9.eyJkYxRhIjp7InVzZXJfaWqiOjIzfSiZXhwjxNAwM0ZwMDB9.jSF1UmIkp7zJqo55o4WU2X7kYW-LNA6Xa0", ErrInvalidFormat},
		{"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InVzZXJfaWQiOjF9LCJleHAiOjMwMDAwMDAwMDB9.eW-eRq0xcPb7yuuugYXNE-SbNy-EgPp2gV0p_LCP7ygj1Y9axqUD2Ng5Oad9Y3-TqwAmlyR8cEuKemkibsYH1Q", ErrInvalidSignature},
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InVzZXJfaWQiOjF9LCJleHAiOjMwMDAwMDAwMDB9.svM08XVNEYi3yoYaoAXjtWs7hEh93IrcO6bjLllVq-U", ErrInvalidSignature},
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InVzZXJfaWQiOjF9fQ.xPviFLZMcK8qNk5zEn0zniWnclL9p0lolm-YDkxP8-s", ErrMissingExpiration},
		{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InVzZXJfaWQiOjF9LCJleHAiOjE3MDAwMDAwMDB9.J6hwK-mNecP2NHYZpZwUgXRqts2Cq1Q9ZUVDs8tc3f0", ErrExceededExpiration},
	}
	for _, test := range errors {
		if _, err = s1.Verify(test.token); err != test.err {
			t.Fatalf("expected error: %v, got: %v", test.err, err)
		}
	}
}
