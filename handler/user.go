package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/awnumar/memguard"
	"github.com/cyberwlodarczyk/auth/argon2id"
	"github.com/cyberwlodarczyk/auth/jwt"
	"github.com/cyberwlodarczyk/auth/postgres"
	"github.com/cyberwlodarczyk/auth/ratelimit"
	"github.com/cyberwlodarczyk/auth/smtp"
	"github.com/cyberwlodarczyk/auth/validation"
)

var (
	errBadUserName             = &operationalError{http.StatusBadRequest, "name is too short or too long"}
	errBadUserEmail            = &operationalError{http.StatusBadRequest, "email is not in the correct format"}
	errBadUserPassword         = &operationalError{http.StatusBadRequest, "password is too weak or too long"}
	errBadUserPasswordEncoding = &operationalError{http.StatusBadRequest, "password is not encoded with standard raw, unpadded base64"}
	errBadUserToken            = &operationalError{http.StatusUnauthorized, "token is invalid or expired"}
	errBadUserSession          = &operationalError{http.StatusUnauthorized, "session is invalid or expired"}
	errMissingUserSession      = &operationalError{http.StatusUnauthorized, "session is missing"}
	errInvalidUserCredentials  = &operationalError{http.StatusUnauthorized, "credentials are invalid"}
	errInvalidUserPassword     = &operationalError{http.StatusUnauthorized, "password is invalid"}
	errUserNotFound            = &operationalError{http.StatusNotFound, "user does not exist"}
	errUserAlreadyExists       = &operationalError{http.StatusConflict, "user already exists"}
)

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
		return errBadUserPasswordEncoding
	}
	*p = b
	return nil
}

func isBadUserToken(err error) error {
	if isJWTErrorOperational(err) {
		return errBadUserToken
	}
	return err
}

func isUserNotFound(err error) error {
	if errors.Is(err, postgres.ErrNotFound) {
		return errUserNotFound
	}
	return err
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

type User struct {
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

func (u *User) WithSession(svc jwt.Service[UserSessionToken], limiter ratelimit.Limiter) func(http.Handler) http.Handler {
	return createMiddleware(func(h http.Handler, w http.ResponseWriter, r *http.Request) error {
		header := strings.Split(r.Header.Get("Authorization"), " ")
		if len(header) != 2 || header[0] != "Bearer" {
			return errMissingUserSession
		}
		token, err := svc.Verify(header[1])
		if err != nil {
			if isJWTErrorOperational(err) {
				return errBadUserSession
			}
			return err
		}
		if !limiter.Allow(strconv.FormatInt(token.Id, 16)) {
			return errTooManyRequests
		}
		h.ServeHTTP(w, setUserID(r, token.Id))
		return nil
	})
}

func (u *User) CreateConfirmationToken(tmpl *template.Template, limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Email string `json:"email"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.EmailValidation.Check(body.Email) {
			err = errBadUserEmail
			return
		}
		if !limiter.Allow(body.Email) {
			err = errTooManyRequests
			return
		}
		token, err := u.ConfirmationToken.Sign(UserConfirmationToken(body))
		if err != nil {
			return
		}
		u.Mail.Send(body.Email, tmpl, token)
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (u *User) CreateSessionToken(ipLimiter ratelimit.Limiter, emailLimiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Email    string       `json:"email"`
		Password userPassword `json:"password"`
	}
	type payload struct {
		Token string `json:"token"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		if !u.EmailValidation.Check(body.Email) {
			err = errBadUserEmail
			return
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return
		}
		if !ipLimiter.Allow(ip) || !emailLimiter.Allow(body.Email) {
			err = errTooManyRequests
			return
		}
		user, err := u.DB.GetByEmail(r.Context(), body.Email)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				err = errInvalidUserCredentials
			}
			return
		}
		match, _, err := u.Password.Compare(body.Password, user.Password)
		if err != nil {
			return
		}
		if !match {
			err = errInvalidUserCredentials
			return
		}
		token, err := u.SessionToken.Sign(UserSessionToken{user.Id})
		if err != nil {
			return
		}
		res = response{http.StatusCreated, payload{token}}
		return
	})
}

func (u *User) CreatePasswordResetToken(tmpl *template.Template, limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Email string `json:"email"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.EmailValidation.Check(body.Email) {
			err = errBadUserEmail
			return
		}
		if !limiter.Allow(body.Email) {
			err = errTooManyRequests
			return
		}
		user, err := u.DB.GetByEmail(r.Context(), body.Email)
		if err != nil {
			if errors.Is(err, postgres.ErrNotFound) {
				err = nil
				res = response{http.StatusNoContent, nil}
			}
			return
		}
		token, err := u.PasswordResetToken.Sign(UserPasswordResetToken{user.Id})
		if err != nil {
			return
		}
		u.Mail.Send(body.Email, tmpl, token)
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (u *User) CreateSudoToken(tmpl *template.Template, limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Password userPassword `json:"password"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		id := getUserID(r)
		user, err := u.DB.GetById(r.Context(), id)
		if err != nil {
			err = isUserNotFound(err)
			return
		}
		match, _, err := u.Password.Compare(body.Password, user.Password)
		if err != nil {
			return
		}
		if !match {
			err = errInvalidUserPassword
			return
		}
		if !limiter.Allow(strconv.FormatInt(id, 16)) {
			err = errTooManyRequests
			return
		}
		token, err := u.SudoToken.Sign(UserSessionToken{id})
		if err != nil {
			return
		}
		u.Mail.Send(user.Email, tmpl, token)
		res = response{http.StatusCreated, nil}
		return
	})
}

func (u *User) Get() http.HandlerFunc {
	type payload struct {
		Id        int64     `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"createdAt"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		user, err := u.DB.GetById(r.Context(), getUserID(r))
		if err != nil {
			err = isUserNotFound(err)
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

func (u *User) Create(limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Token    string       `json:"token"`
		Name     string       `json:"name"`
		Password userPassword `json:"password"`
	}
	type payload struct {
		Id        int64     `json:"id"`
		CreatedAt time.Time `json:"createdAt"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		token, err := u.ConfirmationToken.Verify(body.Token)
		if err != nil {
			err = isBadUserToken(err)
			return
		}
		if !u.NameValidation.Check(body.Name) {
			err = errBadUserName
			return
		}
		if !u.PasswordValidation.Check(body.Password) {
			err = errBadUserPassword
			return
		}
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return
		}
		if !limiter.Allow(ip) {
			err = errTooManyRequests
			return
		}
		hash, err := u.Password.Hash(body.Password)
		if err != nil {
			return
		}
		id, createdAt, err := u.DB.Create(r.Context(), postgres.CreateUserOpts{
			Email:    token.Email,
			Name:     body.Name,
			Password: hash,
		})
		if errors.Is(err, postgres.ErrAlreadyExists) {
			err = errUserAlreadyExists
			return
		}
		if err != nil {
			return
		}
		res = response{http.StatusCreated, payload{id, createdAt}}
		return
	})
}

func (u *User) EditEmail() http.HandlerFunc {
	type body struct {
		Token string `json:"token"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		token, err := u.ConfirmationToken.Verify(body.Token)
		if err != nil {
			err = isBadUserToken(err)
			return
		}
		err = u.DB.EditEmail(r.Context(), getUserID(r), token.Email)
		if err != nil {
			err = isUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (u *User) EditName() http.HandlerFunc {
	type body struct {
		Name string `json:"name"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.NameValidation.Check(body.Name) {
			err = errBadUserName
			return
		}
		err = u.DB.EditName(r.Context(), getUserID(r), body.Name)
		if err != nil {
			err = isUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (u *User) EditPassword() http.HandlerFunc {
	type body struct {
		Password    userPassword `json:"password"`
		NewPassword userPassword `json:"newPassword"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		defer memguard.WipeBytes(body.NewPassword)
		if !u.PasswordValidation.Check(body.NewPassword) {
			err = errBadUserPassword
			return
		}
		id := getUserID(r)
		user, err := u.DB.GetById(r.Context(), id)
		if err != nil {
			err = isUserNotFound(err)
			return
		}
		match, _, err := u.Password.Compare(body.Password, user.Password)
		if err != nil {
			return
		}
		if !match {
			err = errInvalidUserPassword
			return
		}
		hash, err := u.Password.Hash(body.NewPassword)
		if err != nil {
			return
		}
		if err = u.DB.EditPassword(r.Context(), id, hash); err != nil {
			err = isUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (u *User) ResetPassword(limiter ratelimit.Limiter) http.HandlerFunc {
	type body struct {
		Token    string       `json:"token"`
		Password userPassword `json:"password"`
	}
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body body
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		defer memguard.WipeBytes(body.Password)
		if !u.PasswordValidation.Check(body.Password) {
			err = errBadUserPassword
			return
		}
		token, err := u.PasswordResetToken.Verify(body.Token)
		if err != nil {
			err = isBadUserToken(err)
			return
		}
		if !limiter.Allow(strconv.FormatInt(token.Id, 16)) {
			err = errTooManyRequests
			return
		}
		hash, err := u.Password.Hash(body.Password)
		if err != nil {
			return
		}
		err = u.DB.EditPassword(r.Context(), token.Id, hash)
		if err != nil {
			err = isUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

func (u *User) Delete() http.HandlerFunc {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		if err = u.DB.Delete(r.Context(), getUserID(r)); err != nil {
			err = isUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}
