package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	Ping(context.Context) error
	Close()
}

func NewService(ctx context.Context, uri string) (Service, error) {
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

func (s *service) Close() {
	s.pool.Close()
}
