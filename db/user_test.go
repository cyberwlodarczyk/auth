package db

import (
	"context"
	"testing"
)

const (
	email1    = "bar@foo.com"
	email2    = "baz@foo.com"
	name1     = "john"
	name2     = "bob"
	password1 = "pa$$word123"
	password2 = "s3cr3t!"
)

func TestUserService(t *testing.T) {
	ctx := context.Background()
	_, err := db.pool.Exec(
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
	if err != nil {
		t.Fatal(err)
	}
	service := NewUserService(ctx, db)
	id1, createdAt1, err := service.Create(ctx, CreateUserOpts{email1, name1, password1})
	if err != nil {
		t.Fatal(err)
	}
	expected1 := User{id1, email1, name1, password1, createdAt1}
	if _, _, err = service.Create(ctx, CreateUserOpts{email1, name1, password1}); err != ErrAlreadyExists {
		t.Fatalf("expected error: %v, got: %v", ErrAlreadyExists, err)
	}
	id2, createdAt2, err := service.Create(ctx, CreateUserOpts{email2, name2, password2})
	if err != nil {
		t.Fatal(err)
	}
	expected2 := User{id2, email2, name2, password2, createdAt2}
	user, err := service.GetById(ctx, id1)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected1 {
		t.Fatalf("expected user: %v, got: %v", expected1, user)
	}
	user, err = service.GetByEmail(ctx, email2)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected2 {
		t.Fatalf("expected user: %v, got: %v", expected2, user)
	}
	if err = service.EditEmail(ctx, id1, email2); err != ErrAlreadyExists {
		t.Fatalf("expected error: %v, got: %v", ErrAlreadyExists, err)
	}
	if err = service.EditName(ctx, id1, name2); err != nil {
		t.Fatal(err)
	}
	expected1.Name = name2
	if err = service.EditPassword(ctx, id2, password1); err != nil {
		t.Fatal(err)
	}
	expected2.Password = password1
	user, err = service.GetByEmail(ctx, email1)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected1 {
		t.Fatalf("expected user: %v, got: %v", expected1, user)
	}
	user, err = service.GetById(ctx, id2)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected2 {
		t.Fatalf("expected user: %v, got: %v", expected2, user)
	}
	if err = service.Delete(ctx, id2); err != nil {
		t.Fatal(err)
	}
	if err = service.Delete(ctx, id2); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = service.EditEmail(ctx, id1, email2); err != nil {
		t.Fatal(err)
	}
	expected1.Email = email2
	user, err = service.GetById(ctx, id1)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected1 {
		t.Fatalf("expected user: %v, got: %v", expected1, user)
	}
	if err = service.Delete(ctx, id1); err != nil {
		t.Fatal(err)
	}
	if _, err = service.GetById(ctx, id1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if _, err = service.GetByEmail(ctx, email1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = service.EditEmail(ctx, id1, email1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = service.EditName(ctx, id1, password1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = service.EditPassword(ctx, id1, password1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
}
