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

	"github.com/cyberwlodarczyk/auth/api/argon2id"
	"github.com/cyberwlodarczyk/auth/api/config"
	"github.com/cyberwlodarczyk/auth/api/handler"
	"github.com/cyberwlodarczyk/auth/api/jwt"
	"github.com/cyberwlodarczyk/auth/api/postgres"
	"github.com/cyberwlodarczyk/auth/api/ratelimit"
	"github.com/cyberwlodarczyk/auth/api/smtp"
	"github.com/cyberwlodarczyk/auth/api/validation"
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
	root := handler.NewService(&handler.Config{Errors: cfg.Errors.Root})
	userSessionToken := jwt.NewService[handler.UserSessionToken](cfg.JWT.User.Session)
	userSudoToken := jwt.NewService[handler.UserSessionToken](cfg.JWT.User.Sudo)
	user := handler.NewUserService(&handler.UserConfig{
		Errors:             cfg.Errors.User,
		Root:               root,
		DB:                 userDB,
		Mail:               mail,
		ConfirmationToken:  jwt.NewService[handler.UserConfirmationToken](cfg.JWT.User.Confirmation),
		SessionToken:       userSessionToken,
		SudoToken:          userSudoToken,
		PasswordResetToken: jwt.NewService[handler.UserPasswordResetToken](cfg.JWT.User.PasswordReset),
		Password:           argon2id.NewService(argon2id.DefaultParams),
		NameValidation:     validation.NewMinMaxService(cfg.Validation.User.Name),
		EmailValidation:    validation.NewEmailService(validation.DefaultEmailPattern),
		PasswordValidation: validation.NewPasswordService(validation.DefaultPasswordConfig),
	})
	rl := ratelimit.NewService(
		cfg.RateLimit.CleanupInterval,
		cfg.RateLimit.IdleTimeout,
	)
	defer rl.Close()
	r := chi.NewRouter()
	r.Use(root.WithRequestID)
	r.Use(root.WithRequestID)
	r.Use(root.WithRateLimit(rl.NewLimiter(cfg.RateLimit.IP)))
	r.Use(root.WithBodyLimit(int64(cfg.HTTP.BodyLimit)))
	r.NotFound(root.NotFound())
	r.MethodNotAllowed(root.MethodNotAllowed())
	r.Route(cfg.Routes.User.Prefix, func(r chi.Router) {
		session := user.WithSession(userSessionToken, rl.NewLimiter(cfg.RateLimit.User.Session))
		sudo := user.WithSession(userSudoToken, rl.NewLimiter(cfg.RateLimit.User.Sudo))
		r.Post(cfg.Routes.User.Create, user.Create(rl.NewLimiter(cfg.RateLimit.User.Create)))
		r.Post(cfg.Routes.User.ResetPassword, user.ResetPassword(rl.NewLimiter(cfg.RateLimit.User.ResetPassword)))
		r.Group(func(r chi.Router) {
			r.Use(session)
			r.Get(cfg.Routes.User.Get, user.Get())
			r.Put(cfg.Routes.User.EditName, user.EditName())
			r.Put(cfg.Routes.User.EditPassword, user.EditPassword())
		})
		r.Group(func(r chi.Router) {
			r.Use(sudo)
			r.Put(cfg.Routes.User.EditEmail, user.EditEmail())
			r.Delete(cfg.Routes.User.Delete, user.Delete())
		})
		r.Route(cfg.Routes.User.Token.Prefix, func(r chi.Router) {
			r.Post(
				cfg.Routes.User.Token.CreateConfirmation,
				user.CreateConfirmationToken(
					cfg.Mail.User.Confirmation,
					rl.NewLimiter(cfg.RateLimit.User.CreateConfirmationToken),
				),
			)
			r.Post(cfg.Routes.User.Token.CreateSession, user.CreateSessionToken(
				rl.NewLimiter(cfg.RateLimit.User.CreateSessionToken.IP),
				rl.NewLimiter(cfg.RateLimit.User.CreateSessionToken.Email),
			))
			r.Post(
				cfg.Routes.User.Token.CreatePasswordReset,
				user.CreatePasswordResetToken(
					cfg.Mail.User.PasswordReset,
					rl.NewLimiter(cfg.RateLimit.User.CreatePasswordResetToken),
				),
			)
			r.With(session).Post(
				cfg.Routes.User.Token.CreateSudo,
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
