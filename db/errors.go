package db

import "errors"

var (
	ErrNotFound      = errors.New("db: record not found in the database")
	ErrAlreadyExists = errors.New("db: record already exists in the database")
)
