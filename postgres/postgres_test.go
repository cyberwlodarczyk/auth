package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var svc Service

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}
	if err = pool.Client.Ping(); err != nil {
		log.Fatal(err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17-alpine",
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
	pool.MaxWait = 20 * time.Second
	ctx := context.Background()
	if err = pool.Retry(func() error {
		svc, err = NewService(ctx, Config{fmt.Sprintf("postgres://golang:secret@%s/test?sslmode=disable", resource.GetHostPort("5432/tcp"))})
		if err != nil {
			return err
		}
		return svc.Ping(ctx)
	}); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	svc.Close()
	if err = pool.Purge(resource); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}
