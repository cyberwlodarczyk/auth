package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/awnumar/memguard"
	"github.com/cyberwlodarczyk/auth/api/argon2id"
	"github.com/cyberwlodarczyk/auth/api/jwt"
	"github.com/cyberwlodarczyk/auth/api/postgres"
	"github.com/cyberwlodarczyk/auth/api/ratelimit"
	"github.com/cyberwlodarczyk/auth/api/smtp"
	"github.com/cyberwlodarczyk/auth/api/validation"
)

var errUserPassword = errors.New("user password has invalid encoding")

type userPassword []byte

func (p *userPassword) UnmarshalJSON(data []byte) error {
	defer memguard.WipeBytes(data)
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	b, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		memguard.WipeBytes(b)
		return errUserPassword
	}
	*p = b
	return nil
}

type UserConfirmationToken struct {
	Email string `json:"email"`
}

type UserPasswordResetToken struct {
	Id int64 `json:"id"`
}

type UserSessionToken struct {
	Id int64 `json:"id"`
}

type UserTokenMail struct {
	Heading string `yaml:"heading"`
	Action  string `yaml:"action"`
}

func (m UserTokenMail) createTmpl() *template.Template {
	text := strings.ReplaceAll(
		fmt.Sprintf(`<h1>%s</h1>
<p>To %s, please use the following token:</p>
<p><strong>{{.}}</strong></p>
<p><em>Security Notice:</em> Please do not share this token with anyone else. It is confidential and should be kept private.</p>`,
			m.Heading,
			strings.ToLower(m.Action),
		),
		"\n",
		"",
	)
	return template.Must(template.New(m.Action).Parse(text))
}

type UserConfig struct {
	Errors             UserErrors
	Root               *Service
	DB                 postgres.UserService
	Mail               smtp.Service
	ConfirmationToken  jwt.Service[UserConfirmationToken]
	SessionToken       jwt.Service[UserSessionToken]
	SudoToken          jwt.Service[UserSessionToken]
	PasswordResetToken jwt.Service[UserPasswordResetToken]
	Password           argon2id.Service
	NameValidation     validation.Service[string]
	EmailValidation    validation.Service[string]
	PasswordValidation validation.Service[[]byte]
}

type UserErrors struct {
	BadName             string `yaml:"badName"`
	BadEmail            string `yaml:"badEmail"`
	BadPassword         string `yaml:"badPassword"`
	BadPasswordEncoding string `yaml:"badPasswordEncoding"`
	BadToken            string `yaml:"badToken"`
	BadSession          string `yaml:"badSession"`
	MissingSession      string `yaml:"missingSession"`
	InvalidCredentials  string `yaml:"invalidCredentials"`
	InvalidPassword     string `yaml:"invalidPassword"`
	NotFound            string `yaml:"notFound"`
	AlreadyExists       string `yaml:"alreadyExists"`
}

type UserService struct {
	errBadName             error
	errBadEmail            error
	errBadPassword         error
	errBadPasswordEncoding error
	errBadToken            error
	errBadSession          error
	errMissingSession      error
	errInvalidCredentials  error
	errInvalidPassword     error
	errNotFound            error
	errAlreadyExists       error
	root                   *Service
	db                     postgres.UserService
	mail                   smtp.Service
	confirmationToken      jwt.Service[UserConfirmationToken]
	sessionToken           jwt.Service[UserSessionToken]
	sudoToken              jwt.Service[UserSessionToken]
	passwordResetToken     jwt.Service[UserPasswordResetToken]
	password               argon2id.Service
	nameValidation         validation.Service[string]
	emailValidation        validation.Service[string]
	passwordValidation     validation.Service[[]byte]
}

func NewUserService(cfg *UserConfig) *UserService {
	return &UserService{
		errBadName:             &operationalError{http.StatusBadRequest, cfg.Errors.BadName},
		errBadEmail:            &operationalError{http.StatusBadRequest, cfg.Errors.BadEmail},
		errBadPassword:         &operationalError{http.StatusBadRequest, cfg.Errors.BadPassword},
		errBadPasswordEncoding: &operationalError{http.StatusBadRequest, cfg.Errors.BadPasswordEncoding},
		errBadToken:            &operationalError{http.StatusUnauthorized, cfg.Errors.BadToken},
		errBadSession:          &operationalError{http.StatusUnauthorized, cfg.Errors.BadSession},
		errMissingSession:      &operationalError{http.StatusUnauthorized, cfg.Errors.MissingSession},
		errInvalidCredentials:  &operationalError{http.StatusUnauthorized, cfg.Errors.InvalidCredentials},
		errInvalidPassword:     &operationalError{http.StatusUnauthorized, cfg.Errors.InvalidPassword},
		errNotFound:            &operationalError{http.StatusNotFound, cfg.Errors.NotFound},
		errAlreadyExists:       &operationalError{http.StatusConflict, cfg.Errors.AlreadyExists},
		root:                   cfg.Root,
		db:                     cfg.DB,
		mail:                   cfg.Mail,
		confirmationToken:      cfg.ConfirmationToken,
		sessionToken:           cfg.SessionToken,
		sudoToken:              cfg.SudoToken,
		passwordResetToken:     cfg.PasswordResetToken,
		password:               cfg.Password,
		nameValidation:         cfg.NameValidation,
		emailValidation:        cfg.EmailValidation,
		passwordValidation:     cfg.PasswordValidation,
	}
}

func (s *UserService) decodeJSONBody(r *http.Request, v any) error {
	err := s.root.decodeJSONBody(r, v)
	if errors.Is(err, errUserPassword) {
		return s.errBadPasswordEncoding
	}
	return err
}

func (s *UserService) isBadToken(err error) error {
	if isJWTErrorOperational(err) {
		return s.errBadToken
	}
	return err
}

func (s *UserService) isNotFound(err error) error {
	if errors.Is(err, postgres.ErrNotFound) {
		return s.errNotFound
	}
	return err
}

func (s *UserService) WithSession(svc jwt.Service[UserSessionToken], limiter ratelimit.Limiter) func(http.Handler) http.Handler {
	return s.root.createMiddleware(func(h http.Handler, w http.ResponseWriter, r *http.Request) error {
		header := strings.Split(r.Header.Get("Authorization"), " ")
		if len(header) != 2 || header[0] != "Bearer" {
			return s.errMissingSession
		}
		token, err := svc.Verify(header[1])
		if err != nil {
			if isJWTErrorOperational(err) {
				return s.errBadSession
			}
			return err
		}
		if !limiter.Allow(strconv.FormatInt(token.Id, 16)) {
			return s.root.errTooManyRequests
		}
		h.ServeHTTP(w, setUserID(r, token.Id))
		return nil
	})
}

func (s *UserService) CreateConfirmationToken(mail UserTokenMail, limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Email string `json:"email"`
	}
	tmpl := mail.createTmpl()
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		if !s.emailValidation.Check(body.Email) {
			err = s.errBadEmail
			return
		}
		if !limiter.Allow(body.Email) {
			err = s.root.errTooManyRequests
			return
		}
		token, err := s.confirmationToken.Sign(UserConfirmationToken(body))
		if err != nil {
			return
		}
		s.mail.Send(body.Email, tmpl, token)
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (s *UserService) CreateSessionToken(ipLimiter ratelimit.Limiter, emailLimiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Email    string       `json:"email"`
		Password userPassword `json:"password"`
	}
	type payload struct {
		Token string `json:"token"`
	}
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		if !s.emailValidation.Check(body.Email) {
			err = s.errBadEmail
			return
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return
		}
		if !ipLimiter.Allow(ip) || !emailLimiter.Allow(body.Email) {
			err = s.root.errTooManyRequests
			return
		}
		user, err := s.db.GetByEmail(r.Context(), body.Email)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				err = s.errInvalidCredentials
			}
			return
		}
		match, _, err := s.password.Compare(body.Password, user.Password)
		if err != nil {
			return
		}
		if !match {
			err = s.errInvalidCredentials
			return
		}
		token, err := s.sessionToken.Sign(UserSessionToken{user.Id})
		if err != nil {
			return
		}
		res = response{http.StatusCreated, payload{token}}
		return
	})
}

func (s *UserService) CreatePasswordResetToken(mail UserTokenMail, limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Email string `json:"email"`
	}
	tmpl := mail.createTmpl()
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		if !s.emailValidation.Check(body.Email) {
			err = s.errBadEmail
			return
		}
		if !limiter.Allow(body.Email) {
			err = s.root.errTooManyRequests
			return
		}
		user, err := s.db.GetByEmail(r.Context(), body.Email)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				err = nil
				res = response{http.StatusNoContent, nil}
			}
			return
		}
		token, err := s.passwordResetToken.Sign(UserPasswordResetToken{user.Id})
		if err != nil {
			return
		}
		s.mail.Send(body.Email, tmpl, token)
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (s *UserService) CreateSudoToken(mail UserTokenMail, limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Password userPassword `json:"password"`
	}
	tmpl := mail.createTmpl()
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		id := getUserID(r)
		user, err := s.db.GetById(r.Context(), id)
		if err != nil {
			err = s.isNotFound(err)
			return
		}
		match, _, err := s.password.Compare(body.Password, user.Password)
		if err != nil {
			return
		}
		if !match {
			err = s.errInvalidPassword
			return
		}
		if !limiter.Allow(strconv.FormatInt(id, 16)) {
			err = s.root.errTooManyRequests
			return
		}
		token, err := s.sudoToken.Sign(UserSessionToken{id})
		if err != nil {
			return
		}
		s.mail.Send(user.Email, tmpl, token)
		res = response{http.StatusCreated, nil}
		return
	})
}

func (s *UserService) Get() http.HandlerFunc {
	type payload struct {
		Id        int64     `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
	}
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		user, err := s.db.GetById(r.Context(), getUserID(r))
		if err != nil {
			err = s.isNotFound(err)
			return
		}
		res = response{
			http.StatusOK,
			payload{
				user.Id,
				user.Email,
				user.Name,
				user.CreatedAt,
			},
		}
		return
	})
}

func (s *UserService) Create(limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Token    string       `json:"token"`
		Name     string       `json:"name"`
		Password userPassword `json:"password"`
	}
	type payload struct {
		Id        int64     `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
	}
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		token, err := s.confirmationToken.Verify(body.Token)
		if err != nil {
			err = s.isNotFound(err)
			return
		}
		if !s.nameValidation.Check(body.Name) {
			err = s.errBadName
			return
		}
		if !s.passwordValidation.Check(body.Password) {
			err = s.errBadPassword
			return
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return
		}
		if !limiter.Allow(ip) {
			err = s.root.errTooManyRequests
			return
		}
		hash, err := s.password.Hash(body.Password)
		if err != nil {
			return
		}
		id, createdAt, err := s.db.Create(r.Context(), postgres.CreateUserOpts{
			Email:    token.Email,
			Name:     body.Name,
			Password: hash,
		})
		if errors.Is(err, postgres.ErrAlreadyExists) {
			err = s.errAlreadyExists
			return
		}
		if err != nil {
			return
		}
		res = response{http.StatusCreated, payload{id, createdAt}}
		return
	})
}

func (s *UserService) EditEmail() http.HandlerFunc {
	type body struct {
		Token string `json:"token"`
	}
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		token, err := s.confirmationToken.Verify(body.Token)
		if err != nil {
			err = s.isBadToken(err)
			return
		}
		err = s.db.EditEmail(r.Context(), getUserID(r), token.Email)
		if err != nil {
			err = s.isNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (s *UserService) EditName() http.HandlerFunc {
	type body struct {
		Name string `json:"name"`
	}
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		if !s.nameValidation.Check(body.Name) {
			err = s.errBadName
			return
		}
		err = s.db.EditName(r.Context(), getUserID(r), body.Name)
		if err != nil {
			err = s.isNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (s *UserService) EditPassword() http.HandlerFunc {
	type body struct {
		Password    userPassword `json:"password"`
		NewPassword userPassword `json:"newPassword"`
	}
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		defer memguard.WipeBytes(body.NewPassword)
		if !s.passwordValidation.Check(body.NewPassword) {
			err = s.errBadPassword
			return
		}
		id := getUserID(r)
		user, err := s.db.GetById(r.Context(), id)
		if err != nil {
			err = s.isNotFound(err)
			return
		}
		match, _, err := s.password.Compare(body.Password, user.Password)
		if err != nil {
			return
		}
		if !match {
			err = s.errInvalidPassword
			return
		}
		hash, err := s.password.Hash(body.NewPassword)
		if err != nil {
			return
		}
		if err = s.db.EditPassword(r.Context(), id, hash); err != nil {
			err = s.isNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (s *UserService) ResetPassword(limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Token    string       `json:"token"`
		Password userPassword `json:"password"`
	}
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = s.decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		if !s.passwordValidation.Check(body.Password) {
			err = s.errBadPassword
			return
		}
		token, err := s.passwordResetToken.Verify(body.Token)
		if err != nil {
			err = s.isBadToken(err)
			return
		}
		if !limiter.Allow(strconv.FormatInt(token.Id, 16)) {
			err = s.root.errTooManyRequests
			return
		}
		hash, err := s.password.Hash(body.Password)
		if err != nil {
			return
		}
		err = s.db.EditPassword(r.Context(), token.Id, hash)
		if err != nil {
			err = s.isNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (s *UserService) Delete() http.HandlerFunc {
	return s.root.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		if err = s.db.Delete(r.Context(), getUserID(r)); err != nil {
			err = s.isNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}
