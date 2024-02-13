package argon2id

import (
	"regexp"
	"testing"
)

var (
	p1, p2 = &Params{
		Memory:      64 * 1024,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}, &Params{
		Memory:      64 * 1024,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
	s1, s2 = NewService(p1), NewService(p2)
	k1, k2 = []byte("pa$$word123"), []byte("ot#er123")
)

func TestEncode(t *testing.T) {
	pattern, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=1,p=[0-9]{1,4}\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	if err != nil {
		t.Fatal(err)
	}
	salt, err := RandomSalt(p1.SaltLength)
	if err != nil {
		t.Fatal(err)
	}
	hash := Encode(p1, salt, k1)
	if !pattern.MatchString(hash) {
		t.Fatalf("hash is not in the correct format: %q", hash)
	}
}

func TestDecode(t *testing.T) {
	hash, err := s1.Hash(k1)
	if err != nil {
		t.Fatal(err)
	}
	params, _, _, err := Decode(hash)
	if err != nil {
		t.Fatal(err)
	}
	if *params != *p1 {
		t.Fatalf("expected params: %#v, got: %#v", *p1, *params)
	}
	errors := []struct {
		hash string
		err  error
	}{
		{"", ErrInvalidFormat},
		{"argon2id$$v=19$m==65536,t=1,p=1$UXBRcnJTa0p0VnExcmxuNA$O55J0XoGcpMf429Y/Lgq+Erwk8t7xuH++rtFw4boAXQ", ErrInvalidFormat},
		{"$argon2id$v=19$m=65536,t=1,p=1$YzA''''WlWZmcU5RaVN1Qw$P1/FafrWnhJ;Us7YE0wJthWPv0YOPO9jr5rAJgGdA", ErrInvalidFormat},
		{"$argon2d$v=19$m=65536,t=1,p=1$Rm5QSDJhTEh5a3diZjRCYQ$H5WBPRQIhoBlSPwtGBZ20OfBStL6S5BVjpfpF/j3waI", ErrIncompatibleVariant},
		{"$argon2id$v=18$m=65536,t=1,p=1$T3pwRTFrNzlZUFFmUnIybg$oebSb8wynxxVeX03ydN1goSOOtf7WOK3P8jGCqNPLgI", ErrIncompatibleVersion},
	}
	for _, test := range errors {
		_, _, _, err = Decode(test.hash)
		if err != test.err {
			t.Errorf("expected error: %v, got: %v", test.err, err)
		}
	}
}

func TestHash(t *testing.T) {
	h1, err := s1.Hash(k1)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := s1.Hash(k1)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Fatalf("hashes must be unique. got: %q", h1)
	}
}

func TestCompare(t *testing.T) {
	h1, err := s1.Hash(k1)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := s2.Hash(k2)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		password []byte
		hash     string
		match    bool
		rotate   bool
	}{
		{k1, h1, true, false},
		{k2, h1, false, false},
		{k2, h2, true, true},
		{k1, "$argon2id$v=19$m=65536,t=1,p=1$UXdlWUphcVNwOFBJSVJodQ$DpuM7N26KOVtgUUL5GUsyMnHEUjdDnzce7i/I93xgRI", true, false},
		{k1, "$argon2id$v=19$m=65536,t=1,p=1$bG84c2J6YWhiZkhBWGhVdw$K+1Vd+1lEBDBB3CeX/ZPSrUdkuuI3PfWynCxuiod/Uo", false, false},
		{k2, "$argon2id$v=19$m=65536,t=2,p=1$V2xpOG5FTlRHWDZUV09VVA$wYb/TTEdEJ34ADS0S0iiWRb7Oqt81lctiUxtkOjKkkg", true, true},
	}
	for _, test := range tests {
		match, rotate, err := s1.Compare(test.password, test.hash)
		if err != nil {
			t.Fatal(err)
		}
		if match != test.match {
			t.Errorf("expected match: %t", test.match)
		}
		if rotate != test.rotate {
			t.Errorf("expected rotate: %t", test.rotate)
		}
	}
}
