package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id        int64
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
}

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(ctx context.Context, db *DB) *UserService {
	return &UserService{db.pool}
}

func (s *UserService) GetById(ctx context.Context, id int64) (user User, err error) {
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

func (s *UserService) GetByEmail(ctx context.Context, email string) (user User, err error) {
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

type CreateUserOpts struct {
	Email    string
	Name     string
	Password string
}

func (s *UserService) Create(ctx context.Context, opts CreateUserOpts) (id int64, createdAt time.Time, err error) {
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

func (s *UserService) EditEmail(ctx context.Context, id int64, email string) error {
	return isUnique(isAffected(s.pool.Exec(
		ctx,
		"UPDATE user_ SET email = $2 WHERE id = $1",
		id,
		email,
	)))
}

func (s *UserService) EditName(ctx context.Context, id int64, name string) error {
	return isAffected(s.pool.Exec(
		ctx,
		"UPDATE user_ SET name = $2 WHERE id = $1",
		id,
		name,
	))
}

func (s *UserService) EditPassword(ctx context.Context, id int64, password string) error {
	return isAffected(s.pool.Exec(
		ctx,
		"UPDATE user_ SET password = $2 WHERE id = $1",
		id,
		password,
	))
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	return isAffected(s.pool.Exec(
		ctx,
		"DELETE FROM user_ WHERE id = $1",
		id,
	))
}
