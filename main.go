package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cyberwlodarczyk/auth/argon2id"
	"github.com/cyberwlodarczyk/auth/config"
	"github.com/cyberwlodarczyk/auth/handler"
	"github.com/cyberwlodarczyk/auth/jwt"
	"github.com/cyberwlodarczyk/auth/postgres"
	"github.com/cyberwlodarczyk/auth/ratelimit"
	"github.com/cyberwlodarczyk/auth/smtp"
	"github.com/cyberwlodarczyk/auth/validation"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func run(cfg *config.Config) error {
	db, err := postgres.NewService(context.Background(), cfg.Postgres)
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
	cfg.SMTP.ErrorLog = log.New(errorWriter, "", 0)
	cfg.SMTP.TLSConfig = &tls.Config{ServerName: cfg.SMTP.Host}
	mail := smtp.NewService(&cfg.SMTP)
	defer mail.Close()
	user := &handler.User{
		DB:                 userDB,
		Mail:               mail,
		ConfirmationToken:  jwt.NewService[handler.UserConfirmationToken](cfg.JWT.User.Confirmation),
		SessionToken:       jwt.NewService[handler.UserSessionToken](cfg.JWT.User.Session),
		SudoToken:          jwt.NewService[handler.UserSessionToken](cfg.JWT.User.Sudo),
		PasswordResetToken: jwt.NewService[handler.UserPasswordResetToken](cfg.JWT.User.PasswordReset),
		Password:           argon2id.NewService(argon2id.DefaultParams),
		NameValidation:     validation.NewMinMaxService(cfg.Validation.User.Name),
		EmailValidation:    validation.NewEmailService(validation.DefaultEmailPattern),
		PasswordValidation: validation.NewPasswordService(validation.DefaultPasswordConfig),
	}
	rl := ratelimit.NewService(
		cfg.RateLimit.CleanupInterval,
		cfg.RateLimit.IdleTimeout,
	)
	defer rl.Close()
	r := chi.NewRouter()
	r.Use(handler.WithRequestID)
	r.Use(handler.WithRequestID)
	r.Use(handler.WithRateLimit(rl.NewLimiter(cfg.RateLimit.IP)))
	r.Use(handler.WithBodyLimit(int64(cfg.HTTP.BodyLimit)))
	r.NotFound(handler.NotFound())
	r.MethodNotAllowed(handler.MethodNotAllowed())
	r.Route("/user", func(r chi.Router) {
		session := user.WithSession(user.SessionToken, rl.NewLimiter(cfg.RateLimit.User.Session))
		sudo := user.WithSession(user.SudoToken, rl.NewLimiter(cfg.RateLimit.User.Sudo))
		r.Post("/", user.Create(rl.NewLimiter(cfg.RateLimit.User.Create)))
		r.Post("/password-reset", user.ResetPassword(rl.NewLimiter(cfg.RateLimit.User.ResetPassword)))
		r.Group(func(r chi.Router) {
			r.Use(session)
			r.Get("/", user.Get())
			r.Put("/name", user.EditName())
			r.Put("/password", user.EditPassword())
		})
		r.Group(func(r chi.Router) {
			r.Use(sudo)
			r.Put("/email", user.EditEmail())
			r.Delete("/", user.Delete())
		})
		r.Route("/token", func(r chi.Router) {
			r.Post(
				"/confirmation",
				user.CreateConfirmationToken(
					cfg.Mail.User.Confirmation,
					rl.NewLimiter(cfg.RateLimit.User.CreateConfirmationToken),
				),
			)
			r.Post("/session", user.CreateSessionToken(
				rl.NewLimiter(cfg.RateLimit.User.CreateSessionToken.IP),
				rl.NewLimiter(cfg.RateLimit.User.CreateSessionToken.Email),
			))
			r.Post(
				"/password-reset",
				user.CreatePasswordResetToken(
					cfg.Mail.User.PasswordReset,
					rl.NewLimiter(cfg.RateLimit.User.CreatePasswordResetToken),
				),
			)
			r.With(session).Post(
				"/sudo",
				user.CreateSudoToken(
					cfg.Mail.User.Sudo,
					rl.NewLimiter(cfg.RateLimit.User.CreateSudoToken),
				),
			)
		})
	})
	server := &http.Server{
		Addr:           net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:        r,
		TLSConfig:      &tls.Config{MinVersion: tls.VersionTLS13},
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		IdleTimeout:    cfg.HTTP.IdleTimeout,
		MaxHeaderBytes: cfg.HTTP.HeaderLimit,
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
	if len(os.Args) < 2 {
		logrus.Fatal("no config file provided")
	}
	cfg, err := config.New(os.Args[1])
	if err != nil {
		logrus.Fatal(err)
	}
	if err := run(cfg); err != nil {
		logrus.Fatal(err)
	}
}
