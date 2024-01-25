package postgres

import (
	"context"
	"testing"

	"github.com/cyberwlodarczyk/auth/db"
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
	service := svc.NewUserService()
	ctx := context.Background()
	if err = service.Init(ctx); err != nil {
		t.Fatal(err)
	}
	id1, createdAt1, err := service.Create(ctx, db.CreateUserOpts{Email: email1, Name: name1, Password: password1})
	if err != nil {
		t.Fatal(err)
	}
	expected1 := db.User{Id: id1, Email: email1, Name: name1, Password: password1, CreatedAt: createdAt1}
	if _, _, err = service.Create(ctx, db.CreateUserOpts{Email: email1, Name: name1, Password: password1}); err != db.ErrAlreadyExists {
		t.Fatalf("expected error: %v, got: %v", db.ErrAlreadyExists, err)
	}
	id2, createdAt2, err := service.Create(ctx, db.CreateUserOpts{Email: email2, Name: name2, Password: password2})
	if err != nil {
		t.Fatal(err)
	}
	expected2 := db.User{Id: id2, Email: email2, Name: name2, Password: password2, CreatedAt: createdAt2}
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
	if err = service.EditEmail(ctx, id1, email2); err != db.ErrAlreadyExists {
		t.Fatalf("expected error: %v, got: %v", db.ErrAlreadyExists, err)
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
	if err = service.Delete(ctx, id2); err != db.ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", db.ErrNotFound, err)
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
	if _, err = service.GetById(ctx, id1); err != db.ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", db.ErrNotFound, err)
	}
	if _, err = service.GetByEmail(ctx, email1); err != db.ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", db.ErrNotFound, err)
	}
	if err = service.EditEmail(ctx, id1, email1); err != db.ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", db.ErrNotFound, err)
	}
	if err = service.EditName(ctx, id1, password1); err != db.ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", db.ErrNotFound, err)
	}
	if err = service.EditPassword(ctx, id1, password1); err != db.ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", db.ErrNotFound, err)
	}
}
