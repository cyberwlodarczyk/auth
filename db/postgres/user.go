package postgres

import (
	"context"
	"time"

	"github.com/cyberwlodarczyk/auth/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userService struct {
	pool *pgxpool.Pool
}

func (s *userService) Init(ctx context.Context) error {
	_, err := s.pool.Exec(
		ctx,
		`
			CREATE TABLE user_ (
				id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
				email TEXT NOT NULL UNIQUE,
				name TEXT NOT NULL,
				password TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT NOW()
			)
		`,
	)
	return err
}

func (s *userService) GetById(ctx context.Context, id int64) (user db.User, err error) {
	err = isFound(s.pool.QueryRow(
		ctx,
		`
			SELECT id, email, name, password, created_at
			FROM user_
			WHERE id = $1
		`,
		id,
	).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.CreatedAt))
	return
}

func (s *userService) GetByEmail(ctx context.Context, email string) (user db.User, err error) {
	err = isFound(s.pool.QueryRow(
		ctx,
		`
			SELECT id, email, name, password, created_at
			FROM user_
			WHERE email = $1
		`,
		email,
	).Scan(&user.Id, &user.Email, &user.Name, &user.Password, &user.CreatedAt))
	return
}

func (s *userService) Create(ctx context.Context, opts db.CreateUserOpts) (id int64, createdAt time.Time, err error) {
	err = isUnique(s.pool.QueryRow(
		ctx,
		`
			INSERT INTO user_ (email, name, password)
			VALUES ($1, $2, $3)
			RETURNING id, created_at
		`,
		opts.Email,
		opts.Name,
		opts.Password,
	).Scan(&id, &createdAt))
	return
}

func (s *userService) EditEmail(ctx context.Context, id int64, email string) error {
	return isUnique(isAffected(s.pool.Exec(
		ctx,
		"UPDATE user_ SET email = $2 WHERE id = $1",
		id,
		email,
	)))
}

func (s *userService) EditName(ctx context.Context, id int64, name string) error {
	return isAffected(s.pool.Exec(
		ctx,
		"UPDATE user_ SET name = $2 WHERE id = $1",
		id,
		name,
	))
}

func (s *userService) EditPassword(ctx context.Context, id int64, password string) error {
	return isAffected(s.pool.Exec(
		ctx,
		"UPDATE user_ SET password = $2 WHERE id = $1",
		id,
		password,
	))
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	return isAffected(s.pool.Exec(
		ctx,
		"DELETE FROM user_ WHERE id = $1",
		id,
	))
}
