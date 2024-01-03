package argon2id

import (
	"regexp"
	"testing"
)

var password = []byte("pa$$word123")

func TestEncode(t *testing.T) {
	pattern, err := regexp.Compile(`^\$argon2id\$v=19\$m=65536,t=1,p=[0-9]{1,4}\$[A-Za-z0-9+/]{22}\$[A-Za-z0-9+/]{43}$`)
	if err != nil {
		t.Fatal(err)
	}
	salt, err := RandomSalt(DefaultParams.SaltLength)
	if err != nil {
		t.Fatal(err)
	}
	hash := Encode(DefaultParams, salt, password)
	if !pattern.MatchString(hash) {
		t.Fatalf("hash is not in correct format: %q", hash)
	}
}

func TestDecode(t *testing.T) {
	hash, err := Hash(password, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	params, _, _, err := Decode(hash)
	if err != nil {
		t.Fatal(err)
	}
	if *params != *DefaultParams {
		t.Fatalf("expected: %#v, got: %#v", *DefaultParams, *params)
	}
}

func TestHash(t *testing.T) {
	h1, err := Hash(password, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	h2, err := Hash(password, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	if h1 == h2 {
		t.Fatal("hashes must be unique")
	}
}

func TestCompare(t *testing.T) {
	hash, err := Hash(password, DefaultParams)
	if err != nil {
		t.Fatal(err)
	}
	match, _, err := Compare(password, hash)
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatal("expected password and hash to match")
	}
	match, _, err = Compare([]byte("ot#er123"), hash)
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatal("expected password and hash to not match")
	}
	match, _, err = Compare(password, "$argon2id$v=19$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$YD/REpqnz0KJpSvOslIyYw")
	if err != nil {
		t.Fatal(err)
	}
	if !match {
		t.Fatal("expected password and hash to match")
	}
	match, _, err = Compare(password, "$argon2id$v=19$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$YD/REpqnz0KJpSvOslIyYg")
	if err != nil {
		t.Fatal(err)
	}
	if match {
		t.Fatal("expected password and hash to not match")
	}
	match, _, err = Compare(password, "")
	if err != ErrInvalidHash {
		t.Fatalf("expected error: %s", ErrInvalidHash)
	}
	match, _, err = Compare(password, "$argon2d$v=19$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$S/F/LLPkcEMhL8XHcEhfIg")
	if err != ErrIncompatibleVariant {
		t.Fatalf("expected error: %s", ErrIncompatibleVariant)
	}
	match, _, err = Compare(password, "$argon2id$v=18$m=16,t=2,p=1$R0ROc2ViRGt2RjliOUVpaA$YD/REpqnz0KJpSvOslIyYg")
	if err != ErrIncompatibleVersion {
		t.Fatalf("expected error: %s", ErrIncompatibleVersion)
	}
}
