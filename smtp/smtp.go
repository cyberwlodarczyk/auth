package smtp

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Headers struct {
	From    string
	To      string
	Subject string
}

func Write(w io.Writer, headers Headers, tmpl *template.Template, data any, date time.Time) error {
	var err error
	for _, header := range []struct {
		key   string
		value string
	}{
		{"From", fmt.Sprintf("<%s>", headers.From)},
		{"To", fmt.Sprintf("<%s>", headers.To)},
		{"Subject", headers.Subject},
		{"Date", date.Format(time.RFC1123)},
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
	Addr      string
	Host      string
	TLSConfig *tls.Config
	Username  string
	Password  string
}

func Ping(cfg *Config) error {
	c, err := smtp.Dial(cfg.Addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello(cfg.Host); err != nil {
		return err
	}
	if err = c.Noop(); err != nil {
		return err
	}
	return c.Quit()
}

func Send(cfg *Config, headers Headers, tmpl *template.Template, data any) error {
	c, err := smtp.Dial(cfg.Addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Hello(cfg.Host); err != nil {
		return err
	}
	if cfg.TLSConfig != nil {
		if err = c.StartTLS(cfg.TLSConfig); err != nil {
			return err
		}
	}
	if err = c.Auth(sasl.NewPlainClient("", cfg.Username, cfg.Password)); err != nil {
		return err
	}
	if err = c.Mail(headers.From, nil); err != nil {
		return err
	}
	if err = c.Rcpt(headers.To, nil); err != nil {
		return err
	}
	wc, err := c.Data()
	if err != nil {
		return err
	}
	if err = Write(wc, headers, tmpl, data, time.Now()); err != nil {
		return err
	}
	if err = wc.Close(); err != nil {
		return err
	}
	return c.Quit()
}
