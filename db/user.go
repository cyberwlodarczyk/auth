package db

import (
	"context"
	"time"
)

type UserService interface {
	Init(context.Context) error
	GetById(context.Context, int64) (User, error)
	GetByEmail(context.Context, string) (User, error)
	Create(context.Context, CreateUserOpts) (int64, time.Time, error)
	EditEmail(context.Context, int64, string) error
	EditName(context.Context, int64, string) error
	EditPassword(context.Context, int64, string) error
	Delete(context.Context, int64) error
}

type User struct {
	Id        int64
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
}

type CreateUserOpts struct {
	Email    string
	Name     string
	Password string
}
