package postgres

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
	userSvc, err := NewUserService(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	expected1, err := userSvc.Create(ctx, CreateUserOpts{Email: email1, Name: name1, Password: password1})
	if err != nil {
		t.Fatal(err)
	}
	id1 := expected1.Id
	if _, err = userSvc.Create(ctx, CreateUserOpts{Email: email1, Name: name1, Password: password1}); err != ErrAlreadyExists {
		t.Fatalf("expected error: %v, got: %v", ErrAlreadyExists, err)
	}
	expected2, err := userSvc.Create(ctx, CreateUserOpts{Email: email2, Name: name2, Password: password2})
	if err != nil {
		t.Fatal(err)
	}
	id2 := expected2.Id
	user, err := userSvc.GetById(ctx, id1)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected1 {
		t.Fatalf("expected user: %v, got: %v", expected1, user)
	}
	user, err = userSvc.GetByEmail(ctx, email2)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected2 {
		t.Fatalf("expected user: %v, got: %v", expected2, user)
	}
	if err = userSvc.EditEmail(ctx, id1, email2); err != ErrAlreadyExists {
		t.Fatalf("expected error: %v, got: %v", ErrAlreadyExists, err)
	}
	if err = userSvc.EditName(ctx, id1, name2); err != nil {
		t.Fatal(err)
	}
	expected1.Name = name2
	if err = userSvc.EditPassword(ctx, id2, password1); err != nil {
		t.Fatal(err)
	}
	expected2.Password = password1
	user, err = userSvc.GetByEmail(ctx, email1)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected1 {
		t.Fatalf("expected user: %v, got: %v", expected1, user)
	}
	user, err = userSvc.GetById(ctx, id2)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected2 {
		t.Fatalf("expected user: %v, got: %v", expected2, user)
	}
	if err = userSvc.Delete(ctx, id2); err != nil {
		t.Fatal(err)
	}
	if err = userSvc.Delete(ctx, id2); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = userSvc.EditEmail(ctx, id1, email2); err != nil {
		t.Fatal(err)
	}
	expected1.Email = email2
	user, err = userSvc.GetById(ctx, id1)
	if err != nil {
		t.Fatal(err)
	}
	if user != expected1 {
		t.Fatalf("expected user: %v, got: %v", expected1, user)
	}
	if err = userSvc.Delete(ctx, id1); err != nil {
		t.Fatal(err)
	}
	if _, err = userSvc.GetById(ctx, id1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if _, err = userSvc.GetByEmail(ctx, email1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = userSvc.EditEmail(ctx, id1, email1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = userSvc.EditName(ctx, id1, password1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
	if err = userSvc.EditPassword(ctx, id1, password1); err != ErrNotFound {
		t.Fatalf("expected error: %v, got: %v", ErrNotFound, err)
	}
}
