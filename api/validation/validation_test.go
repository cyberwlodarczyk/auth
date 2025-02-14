package validation

import (
	"testing"
)

func TestPasswordService(t *testing.T) {
	svc := NewPasswordService(DefaultPasswordConfig)
	tests := []struct {
		password []byte
		valid    bool
	}{
		{[]byte("Pa$$word1234"), true},
		{[]byte("Cz3sław!!!:D"), true},
		{[]byte("Θct0pu$+!2217"), true},
		{[]byte("βetter3...:)"), false},
		{[]byte("Titanic1912"), false},
		{[]byte(""), false},
	}
	for _, test := range tests {
		if svc.Check(test.password) != test.valid {
			t.Fatalf("expected %t for %q", test.valid, test.password)
		}
	}
}

func TestEmailService(t *testing.T) {
	svc := NewEmailService(DefaultEmailPattern)
	tests := []struct {
		email string
		valid bool
	}{
		{"a@b.cd", true},
		{"john@example.com", true},
		{"a@b", false},
		{"gmail.com", false},
		{"bob@", false},
		{"", false},
	}
	for _, test := range tests {
		if svc.Check(test.email) != test.valid {
			t.Fatalf("expected %t for %q", test.valid, test.email)
		}
	}
}

func TestMinMaxService(t *testing.T) {
	svc := NewMinMaxService(Range{1, 10})
	tests := []struct {
		value string
		valid bool
	}{
		{"aaa", true},
		{"bbbbbbbbbb", true},
		{"c", true},
		{"", false},
		{"dddddddddddd", false},
	}
	for _, test := range tests {
		if svc.Check(test.value) != test.valid {
			t.Fatalf("expected %t for %q", test.valid, test.value)
		}
	}
}
