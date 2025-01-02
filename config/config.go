package config

import (
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/cyberwlodarczyk/auth/handler"
	"github.com/cyberwlodarczyk/auth/jwt"
	"github.com/cyberwlodarczyk/auth/postgres"
	"github.com/cyberwlodarczyk/auth/ratelimit"
	"github.com/cyberwlodarczyk/auth/smtp"
	"github.com/cyberwlodarczyk/auth/validation"
	"github.com/goccy/go-yaml"
)

type Config struct {
	HTTP struct {
		Host         string        `env:"HOST"`
		Port         string        `env:"PORT"`
		TLSCert      string        `env:"TLS_CERT"`
		TLSKey       string        `env:"TLS_KEY"`
		BodyLimit    int           `yaml:"bodyLimit"`
		HeaderLimit  int           `yaml:"headerLimit"`
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
		IdleTimeout  time.Duration `yaml:"idleTimeout"`
	} `yaml:"http" envPrefix:"HTTP_"`
	Routes struct {
		User struct {
			Prefix        string `yaml:"_prefix"`
			Get           string `yaml:"get"`
			Create        string `yaml:"create"`
			Delete        string `yaml:"delete"`
			ResetPassword string `yaml:"resetPassword"`
			EditName      string `yaml:"editName"`
			EditPassword  string `yaml:"editPassword"`
			EditEmail     string `yaml:"editEmail"`
			Token         struct {
				Prefix              string `yaml:"_prefix"`
				CreateConfirmation  string `yaml:"createConfirmation"`
				CreateSession       string `yaml:"createSession"`
				CreatePasswordReset string `yaml:"createPasswordReset"`
				CreateSudo          string `yaml:"createSudo"`
			} `yaml:"token"`
		} `yaml:"user"`
	} `yaml:"routes"`
	Errors struct {
		Root handler.Errors     `yaml:"root"`
		User handler.UserErrors `yaml:"user"`
	} `yaml:"errors"`
	Validation struct {
		User struct {
			Name     validation.Range          `yaml:"name"`
			Password validation.PasswordConfig `yaml:"password"`
		} `yaml:"user"`
	} `yaml:"validation"`
	RateLimit struct {
		CleanupInterval time.Duration    `yaml:"cleanupInterval"`
		IdleTimeout     time.Duration    `yaml:"idleTimeout"`
		IP              ratelimit.Params `yaml:"ip"`
		User            struct {
			Session                  ratelimit.Params `yaml:"session"`
			Sudo                     ratelimit.Params `yaml:"sudo"`
			Create                   ratelimit.Params `yaml:"create"`
			ResetPassword            ratelimit.Params `yaml:"resetPassword"`
			CreateConfirmationToken  ratelimit.Params `yaml:"createConfirmationToken"`
			CreatePasswordResetToken ratelimit.Params `yaml:"createPasswordResetToken"`
			CreateSessionToken       struct {
				IP    ratelimit.Params `yaml:"ip"`
				Email ratelimit.Params `yaml:"email"`
			} `yaml:"createSessionToken"`
			CreateSudoToken ratelimit.Params `yaml:"createSudoToken"`
		} `yaml:"user"`
	} `yaml:"rateLimit"`
	Mail struct {
		User struct {
			Confirmation  handler.UserTokenMail `yaml:"confirmation"`
			PasswordReset handler.UserTokenMail `yaml:"passwordReset"`
			Sudo          handler.UserTokenMail `yaml:"sudo"`
		} `yaml:"user"`
	} `yaml:"mail"`
	JWT struct {
		User struct {
			Confirmation  jwt.Config `yaml:"confirmation" envPrefix:"CONFIRMATION_"`
			Session       jwt.Config `yaml:"session" envPrefix:"SESSION_"`
			PasswordReset jwt.Config `yaml:"passwordReset" envPrefix:"PASSWORD_RESET_"`
			Sudo          jwt.Config `yaml:"sudo" envPrefix:"SUDO_"`
		} `yaml:"user" envPrefix:"USER_"`
	} `yaml:"jwt" envPrefix:"JWT_"`
	SMTP     smtp.Config     `yaml:"smtp" envPrefix:"SMTP_"`
	Postgres postgres.Config `envPrefix:"POSTGRES_"`
}

func New(file string) (*Config, error) {
	c, err := env.ParseAsWithOptions[Config](env.Options{RequiredIfNoDef: true})
	if err != nil {
		return nil, err
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.NewDecoder(f, yaml.Strict()).Decode(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
