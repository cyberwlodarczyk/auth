package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cyberwlodarczyk/auth/db"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var svc db.Service

func TestMain(m *testing.M) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}
	if err = dockerPool.Client.Ping(); err != nil {
		log.Fatal(err)
	}
	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "16",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=golang",
			"POSTGRES_DB=test",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatal(err)
	}
	resource.Expire(20)
	dockerPool.MaxWait = 20 * time.Second
	ctx := context.Background()
	if err = dockerPool.Retry(func() error {
		svc, err = NewService(ctx, fmt.Sprintf("postgres://golang:secret@%s/test?sslmode=disable", resource.GetHostPort("5432/tcp")))
		if err != nil {
			return err
		}
		return svc.Ping(ctx)
	}); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	svc.Close()
	if err = dockerPool.Purge(resource); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
