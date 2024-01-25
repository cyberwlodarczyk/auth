package db

import "context"

type Service interface {
	Ping(context.Context) error
	NewUserService() UserService
	Close()
}
