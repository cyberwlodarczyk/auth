package handler

import (
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/awnumar/memguard"
	"github.com/cyberwlodarczyk/auth/argon2id"
	"github.com/cyberwlodarczyk/auth/jwt"
	"github.com/cyberwlodarczyk/auth/postgres"
	"github.com/cyberwlodarczyk/auth/smtp"
	"github.com/cyberwlodarczyk/auth/validation"
)

var (
	errBadUserName        = &operationalError{http.StatusBadRequest, "name is too short or too long"}
	errBadUserEmail       = &operationalError{http.StatusBadRequest, "email is not in the correct format"}
	errBadUserPassword    = &operationalError{http.StatusBadRequest, "password is too weak or too long"}
	errBadUserToken       = &operationalError{http.StatusUnauthorized, "token is invalid or expired"}
	errBadUserSession     = &operationalError{http.StatusUnauthorized, "session is invalid or expired"}
	errMissingUserSession = &operationalError{http.StatusUnauthorized, "session is missing"}
	errUserNotFound       = &operationalError{http.StatusNotFound, "user does not exist"}
	errUserAlreadyExists  = &operationalError{http.StatusConflict, "user already exists"}
)

func checkBadUserToken(err error) error {
	if isJWTErrorOperational(err) {
		return errBadUserToken
	}
	return err
}

func checkUserNotFound(err error) error {
	if errors.Is(err, postgres.ErrNotFound) {
		return errUserNotFound
	}
	return err
}

type UserConfirmationToken struct {
	Email string `json:"email"`
}

type UserSessionToken struct {
	Id int64 `json:"id"`
}

type UserPasswordResetToken struct {
	Id int64 `json:"id"`
}

type UserSudoToken struct {
	Id int64 `json:"id"`
}

type User struct {
	Db                 postgres.UserService
	Mail               smtp.Service
	ConfirmationToken  jwt.Service[UserConfirmationToken]
	SessionToken       jwt.Service[UserSessionToken]
	PasswordResetToken jwt.Service[UserPasswordResetToken]
	SudoToken          jwt.Service[UserSudoToken]
	Password           argon2id.Service
	NameValidation     validation.Service[string]
	EmailValidation    validation.Service[string]
	PasswordValidation validation.Service[[]byte]
}

func (u *User) WithSession() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
			header := strings.Split(r.Header.Get("Authorization"), " ")
			if len(header) != 2 || header[0] != "Bearer" {
				err = errMissingUserSession
				return
			}
			token, err := u.SessionToken.Verify(header[1])
			if err != nil {
				if isJWTErrorOperational(err) {
					err = errBadUserSession
				}
				return
			}
			h.ServeHTTP(w, setUserID(r, token.Id))
			return
		})
	}
}

type CreateUserConfirmationTokenBody struct {
	Email string `json:"email"`
}

func (u *User) CreateConfirmationToken(tmpl *template.Template) http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body CreateUserConfirmationTokenBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.EmailValidation.Check(body.Email) {
			err = errBadUserEmail
			return
		}
		token, err := u.ConfirmationToken.Sign(UserConfirmationToken(body))
		if err != nil {
			return
		}
		go u.Mail.Send(body.Email, tmpl, token)
		res = response{http.StatusNoContent, nil}
		return
	})
}

type CreateUserSessionBody struct {
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type CreateUserSessionReply struct {
	Token string `json:"token"`
}

func (u *User) CreateSessionToken() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body CreateUserSessionBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		user, err := u.Db.GetByEmail(r.Context(), body.Email)
		if err != nil {
			return
		}
		match, _, err := u.Password.Compare(body.Password, user.Password)
		memguard.WipeBytes(body.Password)
		if err != nil {
			return
		}
		if !match {
			return
		}
		token, err := u.SessionToken.Sign(UserSessionToken{user.Id})
		if err != nil {
			return
		}
		res = response{http.StatusCreated, CreateUserSessionReply{token}}
		return
	})
}

type CreateUserPasswordResetTokenBody struct {
	Email string `json:"email"`
}

func (u *User) CreatePasswordResetToken(tmpl *template.Template) http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body CreateUserPasswordResetTokenBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		user, err := u.Db.GetByEmail(r.Context(), body.Email)
		if err != nil {
			return
		}
		token, err := u.PasswordResetToken.Sign(UserPasswordResetToken{user.Id})
		if err != nil {
			return
		}
		go u.Mail.Send(
			body.Email,
			tmpl,
			token,
		)
		res = response{http.StatusNoContent, nil}
		return
	})
}

type CreateUserSudoTokenBody struct {
	Password []byte `json:"password"`
}

func (u *User) CreateSudoToken(tmpl *template.Template) http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body CreateUserSudoTokenBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		id := getUserID(r)
		user, err := u.Db.GetById(r.Context(), id)
		if err != nil {
			return
		}
		match, _, err := u.Password.Compare(body.Password, user.Password)
		memguard.WipeBytes(body.Password)
		if err != nil {
			return
		}
		if !match {
			return
		}
		token, err := u.SudoToken.Sign(UserSudoToken{id})
		if err != nil {
			return
		}
		go u.Mail.Send(user.Email, tmpl, token)
		res = response{http.StatusCreated, nil}
		return
	})
}

type GetUserReply struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) Get() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		id := getUserID(r)
		user, err := u.Db.GetById(r.Context(), id)
		if err != nil {
			err = checkUserNotFound(err)
			return
		}
		res = response{
			http.StatusOK,
			GetUserReply{
				user.Id,
				user.Email,
				user.Name,
				user.CreatedAt,
			},
		}
		return
	})
}

type CreateUserBody struct {
	Token    string `json:"token"`
	Name     string `json:"name"`
	Password []byte `json:"password"`
}

type CreateUserReply struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

func (u *User) Create() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body CreateUserBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		token, err := u.ConfirmationToken.Verify(body.Token)
		if err != nil {
			err = checkBadUserToken(err)
			return
		}
		if !u.PasswordValidation.Check(body.Password) {
			err = errBadUserPassword
			return
		}
		hash, err := u.Password.Hash(body.Password)
		memguard.WipeBytes(body.Password)
		if err != nil {
			return
		}
		id, createdAt, err := u.Db.Create(r.Context(), postgres.CreateUserOpts{
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
		res = response{http.StatusCreated, CreateUserReply{id, createdAt}}
		return
	})
}

type EditUserEmailBody struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func (u *User) EditEmail() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body EditUserEmailBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.EmailValidation.Check(body.Email) {
			err = errBadUserEmail
			return
		}
		data, err := u.SudoToken.Verify(body.Token)
		if err != nil {
			err = checkBadUserToken(err)
			return
		}
		id := getUserID(r)
		if id != data.Id {
			return
		}
		err = u.Db.EditEmail(r.Context(), id, body.Email)
		if err != nil {
			err = checkUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

type EditUserNameBody struct {
	Name string `json:"name"`
}

func (u *User) EditName() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body EditUserNameBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.NameValidation.Check(body.Name) {
			err = errBadUserName
			return
		}
		id := getUserID(r)
		err = u.Db.EditName(r.Context(), id, body.Name)
		if err != nil {
			err = checkUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

type EditUserPasswordBody struct {
	Password    []byte `json:"password"`
	NewPassword []byte `json:"newPassword"`
}

func (u *User) EditPassword() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body EditUserPasswordBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.PasswordValidation.Check(body.NewPassword) {
			err = errBadUserPassword
			return
		}
		id := getUserID(r)
		user, err := u.Db.GetById(r.Context(), id)
		if err != nil {
			return
		}
		match, _, err := u.Password.Compare(body.Password, user.Password)
		memguard.WipeBytes(body.Password)
		if err != nil {
			return
		}
		if !match {
			return
		}
		hash, err := u.Password.Hash(body.NewPassword)
		memguard.WipeBytes(body.NewPassword)
		if err != nil {
			return
		}
		err = u.Db.EditPassword(r.Context(), id, hash)
		if err != nil {
			err = checkUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

type ResetUserPasswordBody struct {
	Token    string `json:"token"`
	Password []byte `json:"password"`
}

func (u *User) ResetPassword() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body ResetUserPasswordBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		if !u.PasswordValidation.Check(body.Password) {
			err = errBadUserPassword
			return
		}
		token, err := u.PasswordResetToken.Verify(body.Token)
		if err != nil {
			err = checkBadUserToken(err)
			return
		}
		hash, err := u.Password.Hash(body.Password)
		memguard.WipeBytes(body.Password)
		if err != nil {
			return
		}
		err = u.Db.EditPassword(r.Context(), token.Id, hash)
		if err != nil {
			err = checkUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}

type DeleteUserBody struct {
	Token string `json:"token"`
}

func (u *User) Delete() http.Handler {
	return createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		var body DeleteUserBody
		if err = decodeJSONBody(r, &body); err != nil {
			return
		}
		data, err := u.SudoToken.Verify(body.Token)
		if err != nil {
			err = checkBadUserToken(err)
			return
		}
		id := getUserID(r)
		if id != data.Id {
			return
		}
		err = u.Db.Delete(r.Context(), id)
		if err != nil {
			err = checkUserNotFound(err)
			return
		}
		res = response{http.StatusNoContent, nil}
		return
	})
}
