package postgres

import (
	"context"
	"net"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Database string `env:"DATABASE"`
}

type Service interface {
	Ping(context.Context) error
	Close()
}

func NewService(ctx context.Context, cfg Config) (Service, error) {
	u := &url.URL{
		Scheme: "postgres",
		Host:   net.JoinHostPort(cfg.Host, cfg.Port),
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Path:   cfg.Database,
	}
	pool, err := pgxpool.New(ctx, u.String())
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
