package smtp

import (
	"html/template"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	svc     Service
	headers = Headers{
		From:    "john@example.com",
		To:      "bob@example.com",
		Subject: "Greeting",
		Date:    time.Time{},
	}
	tmpl = template.Must(template.New("email").Parse(`<h1>Hello, {{.Name}}!</h1><p>This is a sample HTML email content.</p>`))
	data = struct{ Name string }{Name: "Bob"}
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal(err)
	}
	if err = pool.Client.Ping(); err != nil {
		log.Fatal(err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "maildev/maildev",
		Tag:        "2.1.0",
		Env: []string{
			"MAILDEV_INCOMING_USER=golang",
			"MAILDEV_INCOMING_PASS=secret",
			"MAILDEV_DISABLE_WEB=true",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatal(err)
	}
	resource.Expire(20)
	svc = NewService(&Config{
		Addr:     resource.GetHostPort("1025/tcp"),
		Host:     "localhost",
		Username: "golang",
		Password: "secret",
		From:     headers.From,
	})
	pool.MaxWait = 20 * time.Second
	if err = pool.Retry(svc.Ping); err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	if err = pool.Purge(resource); err != nil {
		log.Fatal(err)
	}
	os.Exit(code)
}

func TestWrite(t *testing.T) {
	var sb strings.Builder
	if err := Write(&sb, headers, tmpl, data); err != nil {
		t.Fatal(err)
	}
	expected := "From: <john@example.com>\r\nTo: <bob@example.com>\r\nSubject: Greeting\r\nDate: Mon, 01 Jan 0001 00:00:00 UTC\r\nMime-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n<h1>Hello, Bob!</h1><p>This is a sample HTML email content.</p>\r\n"
	got := sb.String()
	if got != expected {
		t.Fatalf("expected: %s, got: %s", expected, got)
	}
}

func TestSend(t *testing.T) {
	if err := svc.Send(headers.To, headers.Subject, tmpl, data); err != nil {
		t.Fatal(err)
	}
}
