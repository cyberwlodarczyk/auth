package smtp

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Headers struct {
	From string
	To   string
	Date time.Time
}

func Write(w io.Writer, headers Headers, tmpl *template.Template, data any) error {
	var err error
	for _, header := range []struct {
		key   string
		value string
	}{
		{"From", fmt.Sprintf("<%s>", headers.From)},
		{"To", fmt.Sprintf("<%s>", headers.To)},
		{"Subject", tmpl.Name()},
		{"Date", headers.Date.Format(time.RFC1123)},
		{"Mime-Version", "1.0"},
		{"Content-Type", `text/html; charset="utf-8"`},
	} {
		if _, err = fmt.Fprintf(w, "%s: %s\r\n", header.key, header.value); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}
	if err = tmpl.Execute(w, data); err != nil {
		return err
	}
	if _, err = fmt.Fprint(w, "\r\n"); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Host      string `env:"HOST"`
	Port      string `env:"PORT"`
	Username  string `env:"USERNAME"`
	Password  string `env:"PASSWORD"`
	Name      string `yaml:"name"`
	From      string `yaml:"from"`
	ErrorLog  *log.Logger
	TLSConfig *tls.Config
}

type Service interface {
	Ping() error
	Send(to string, tmpl *template.Template, data any)
	Close()
}

func NewService(cfg *Config) Service {
	if cfg.ErrorLog == nil {
		cfg.ErrorLog = log.Default()
	}
	return &service{
		addr:      net.JoinHostPort(cfg.Host, cfg.Port),
		auth:      sasl.NewPlainClient("", cfg.Username, cfg.Password),
		name:      cfg.Name,
		from:      cfg.From,
		errorLog:  cfg.ErrorLog,
		tlsConfig: cfg.TLSConfig,
	}
}

type service struct {
	addr      string
	auth      sasl.Client
	name      string
	from      string
	tlsConfig *tls.Config
	errorLog  *log.Logger
	wg        sync.WaitGroup
}

func (s *service) connect() (*smtp.Client, error) {
	if s.tlsConfig != nil {
		return smtp.DialStartTLS(s.addr, s.tlsConfig)
	} else {
		return smtp.Dial(s.addr)
	}
}

func (s *service) Ping() error {
	c, err := s.connect()
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello(s.name); err != nil {
		return err
	}
	if err = c.Noop(); err != nil {
		return err
	}
	return c.Quit()
}

func (s *service) send(to string, tmpl *template.Template, data any) error {
	c, err := s.connect()
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello(s.name); err != nil {
		return err
	}
	if err = c.Auth(s.auth); err != nil {
		return err
	}
	if err = c.Mail(s.from, nil); err != nil {
		return err
	}
	if err = c.Rcpt(to, nil); err != nil {
		return err
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	if err = Write(
		wc,
		Headers{s.from, to, time.Now()},
		tmpl,
		data,
	); err != nil {
		return err
	}
	if err = wc.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func (s *service) Send(to string, tmpl *template.Template, data any) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.send(to, tmpl, data); err != nil {
			s.errorLog.Print(err)
		}
	}()
}

func (s *service) Close() {
	s.wg.Wait()
}
