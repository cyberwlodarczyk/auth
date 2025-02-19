package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserService interface {
	GetById(context.Context, int64) (User, error)
	GetByEmail(context.Context, string) (User, error)
	Create(context.Context, CreateUserOpts) (User, error)
	EditEmail(context.Context, int64, string) error
	EditName(context.Context, int64, string) error
	EditPassword(context.Context, int64, string) error
	Delete(context.Context, int64) error
}

func NewUserService(ctx context.Context, svc Service) (UserService, error) {
	pool := svc.(*service).pool
	if _, err := pool.Exec(
		ctx,
		`
			CREATE TABLE IF NOT EXISTS user_ (
				id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
				email TEXT NOT NULL UNIQUE,
				name TEXT NOT NULL,
				password TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT NOW()
			)
		`,
	); err != nil {
		return nil, err
	}
	return &userService{pool}, nil
}

type userService struct {
	pool *pgxpool.Pool
}

func (s *userService) GetById(ctx context.Context, id int64) (user User, err error) {
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

func (s *userService) GetByEmail(ctx context.Context, email string) (user User, err error) {
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

func (s *userService) Create(ctx context.Context, opts CreateUserOpts) (user User, err error) {
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
	).Scan(&user.Id, &user.CreatedAt))
	if err != nil {
		return
	}
	user.Email = opts.Email
	user.Name = opts.Name
	user.Password = opts.Password
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
