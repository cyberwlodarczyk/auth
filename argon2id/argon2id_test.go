package argon2id

import (
	"regexp"
	"testing"
)

var p1, p2 = []byte("pa$$word123"), []byte("ot#er123")

func TestEncode(t *testing.T) {
	pattern, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=1,p=[0-9]{1,4}\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	if err != nil {
		t.Fatal(err)
	}
	salt, err := RandomSalt(DefaultParams.SaltLength)
	if err != nil {
		t.Fatal(err)
	}
	hash := Encode(DefaultParams, salt, p1)
	if !pattern.MatchString(hash) {
		t.Fatalf("hash is not in the correct format: %q", hash)
	}
}

func TestDecode(t *testing.T) {
	hash, err := Hash(p1, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	params, _, _, err := Decode(hash)
	if err != nil {
		t.Fatal(err)
	}
	if *params != *DefaultParams {
		t.Fatalf("expected params: %#v, got: %#v", *DefaultParams, *params)
	}
}

func TestHash(t *testing.T) {
	h1, err := Hash(p1, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := Hash(p1, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Fatalf("hashes must be unique. got: %q", h1)
	}
}

func TestCompare(t *testing.T) {
	h1, err := Hash(p1, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		password []byte
		hash     string
		match    bool
	}{
		{p1, h1, true},
		{p2, h1, false},
		{p1, "$argon2id$v=19$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$YD/REpqnz0KJpSvOslIyYw", true},
		{p1, "$argon2id$v=19$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$YD/REpqnz0KJpSvOslIyYg", false},
	}
	for _, test := range tests {
		match, _, err := Compare(test.password, test.hash)
		if err != nil {
			t.Fatal(err)
		}
		if match != test.match {
			t.Errorf("expected match: %t", test.match)
		}
	}
	errors := []struct {
		hash string
		err  error
	}{
		{"", ErrMalformedHash},
		{"$argonid$v19$=16,t=2,p=1R0Rc2ViGt2RjliOUpaA$YDREpqnz0KpSvOsIyYg", ErrMalformedHash},
		{"$argon2d$v=19$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$S/F/LLPkcEMhL8XHcEhfIg", ErrIncompatibleVariant},
		{"$argon2id$v=18$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$YD/REpqnz0KJpSvOslIyYg", ErrIncompatibleVersion},
	}
	for _, test := range errors {
		_, _, err = Compare(p1, test.hash)
		if err != test.err {
			t.Errorf("expected error: %v, got: %v", test.err, err)
		}
	}
}
