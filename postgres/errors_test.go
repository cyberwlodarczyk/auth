package postgres

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var err = errors.New("something went wrong")

func TestIsFound(t *testing.T) {
	tests := []struct {
		given    error
		expected error
	}{
		{nil, nil},
		{err, err},
		{pgx.ErrNoRows, ErrNotFound},
	}
	for _, test := range tests {
		if err = isFound(test.given); err != test.expected {
			t.Fatalf("expected error: %v, got: %v", test.expected, err)
		}
	}
}

func TestIsUnique(t *testing.T) {
	tests := []struct {
		given    error
		expected error
	}{
		{nil, nil},
		{err, err},
		{&pgconn.PgError{Code: "23505"}, ErrAlreadyExists},
	}
	for _, test := range tests {
		if err = isUnique(test.given); err != test.expected {
			t.Fatalf("expected error: %v, got: %v", test.expected, err)
		}
	}
}

func TestIsAffected(t *testing.T) {
	tests := []struct {
		tag      pgconn.CommandTag
		given    error
		expected error
	}{
		{pgconn.NewCommandTag("DELETE 1"), nil, nil},
		{pgconn.NewCommandTag("DELETE 0"), err, err},
		{pgconn.NewCommandTag("UPDATE 1"), err, err},
		{pgconn.NewCommandTag("UPDATE 0"), nil, ErrNotFound},
	}
	for _, test := range tests {
		if err = isAffected(test.tag, test.given); err != test.expected {
			t.Fatalf("expected error: %v, got: %v", test.expected, err)
		}
	}
}
