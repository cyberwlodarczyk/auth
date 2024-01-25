package postgres

import (
	"context"

	"github.com/cyberwlodarczyk/auth/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewService(ctx context.Context, uri string) (db.Service, error) {
	pool, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}
	return &service{pool}, nil
}

type service struct {
	pool *pgxpool.Pool
}

func (s *service) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}

func (s *service) NewUserService() db.UserService {
	return &userService{s.pool}
}

func (s *service) Close() {
	s.pool.Close()
}
