package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/cyberwlodarczyk/auth/argon2id"
	"github.com/cyberwlodarczyk/auth/handler"
	"github.com/cyberwlodarczyk/auth/jwt"
	"github.com/cyberwlodarczyk/auth/postgres"
	"github.com/cyberwlodarczyk/auth/smtp"
	"github.com/cyberwlodarczyk/auth/validation"
	"github.com/sirupsen/logrus"
)

type secret []byte

func (s *secret) UnmarshalText(src []byte) error {
	b, err := base64.RawStdEncoding.DecodeString(string(src))
	if err != nil {
		return err
	}
	*s = b
	return nil
}

type httpConfig struct {
	Host    string `env:"HOST"`
	Port    string `env:"PORT"`
	TLSCert string `env:"TLS_CERT"`
	TLSKey  string `env:"TLS_KEY"`
}

type postgresConfig struct {
	URI string `env:"URI"`
}

type jwtConfig struct {
	ConfirmationSecret  secret `env:"CONFIRMATION_SECRET"`
	SessionSecret       secret `env:"SESSION_SECRET"`
	PasswordResetSecret secret `env:"PASSWORD_RESET_SECRET"`
	SudoSecret          secret `env:"SUDO_SECRET"`
}

type smtpConfig struct {
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
	Name     string `env:"NAME"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	From     string `env:"FROM"`
}

type config struct {
	HTTP     httpConfig     `envPrefix:"HTTP_"`
	Postgres postgresConfig `envPrefix:"POSTGRES_"`
	JWT      jwtConfig      `envPrefix:"JWT_"`
	SMTP     smtpConfig     `envPrefix:"SMTP_"`
}

func createUserTokenTmpl(heading, action string) *template.Template {
	text := strings.ReplaceAll(
		fmt.Sprintf(`<h1>%s</h1>
<p>To %s, please use the following token:</p>
<p><strong>{{.}}</strong></p>
<p><em>Security Notice:</em> Please do not share this token with anyone else. It is confidential and should be kept private.</p>`,
			heading,
			strings.ToLower(action),
		),
		"\n",
		"",
	)
	return template.Must(template.New(action).Parse(text))
}

func run(cfg *config) error {
	db, err := postgres.NewService(context.Background(), cfg.Postgres.URI)
	if err != nil {
		return err
	}
	defer db.Close()
	userDB, err := postgres.NewUserService(context.Background(), db)
	if err != nil {
		return err
	}
	errorWriter := logrus.StandardLogger().WriterLevel(logrus.ErrorLevel)
	defer errorWriter.Close()
	mail := smtp.NewService(&smtp.Config{
		Host:      cfg.SMTP.Host,
		Port:      cfg.SMTP.Port,
		Name:      cfg.SMTP.Name,
		Username:  cfg.SMTP.Username,
		Password:  cfg.SMTP.Password,
		From:      cfg.SMTP.From,
		ErrorLog:  log.New(errorWriter, "", 0),
		TLSConfig: &tls.Config{ServerName: cfg.SMTP.Host},
	})
	defer mail.Close()
	user := &handler.User{
		DB:   userDB,
		Mail: mail,
		ConfirmationToken: jwt.NewService[handler.UserConfirmationToken](
			cfg.JWT.ConfirmationSecret,
			15*time.Minute,
		),
		SessionToken: jwt.NewService[handler.UserSessionToken](
			cfg.JWT.SessionSecret,
			7*24*time.Hour,
		),
		SudoToken: jwt.NewService[handler.UserSessionToken](
			cfg.JWT.SudoSecret,
			5*time.Minute,
		),
		PasswordResetToken: jwt.NewService[handler.UserPasswordResetToken](
			cfg.JWT.PasswordResetSecret,
			15*time.Minute,
		),
		ConfirmatonTmpl:    createUserTokenTmpl("Email confirmation", "Confirm your email"),
		SudoTmpl:           createUserTokenTmpl("Performing sensitive action", "Perform sensitive action"),
		PasswordResetTmpl:  createUserTokenTmpl("Password reset", "Reset your password"),
		Password:           argon2id.NewService(argon2id.DefaultParams),
		NameValidation:     validation.NewMinMaxService(1, 1000),
		EmailValidation:    validation.NewEmailService(validation.DefaultEmailPattern),
		PasswordValidation: validation.NewPasswordService(validation.DefaultPasswordConfig),
	}
	server := &http.Server{
		Addr:           net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:        handler.New(user),
		TLSConfig:      &tls.Config{MinVersion: tls.VersionTLS13},
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 12,
		ErrorLog:       log.New(errorWriter, "", 0),
	}
	done := make(chan error, 1)
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		done <- server.Shutdown(ctx)
	}()
	if err = server.ListenAndServeTLS(cfg.HTTP.TLSCert, cfg.HTTP.TLSKey); err != nil && err != http.ErrServerClosed {
		return err
	}
	return <-done
}

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	var cfg config
	if err := env.ParseWithOptions(
		&cfg,
		env.Options{RequiredIfNoDef: true},
	); err != nil {
		logrus.Fatal(err)
	}
	if err := run(&cfg); err != nil {
		logrus.Fatal(err)
	}
}
