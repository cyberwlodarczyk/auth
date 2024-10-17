package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cyberwlodarczyk/auth/jwt"
	"github.com/cyberwlodarczyk/auth/ratelimit"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type contextKey int

const (
	contextKeyRequestID contextKey = iota
	contextKeyRequestTime
	contextKeyUserID
)

func createContextHelpers[T any](key contextKey) (get func(*http.Request) T, is func(*http.Request) bool, set func(*http.Request, T) *http.Request) {
	get = func(r *http.Request) T {
		return r.Context().Value(key).(T)
	}
	is = func(r *http.Request) bool {
		_, ok := r.Context().Value(key).(T)
		return ok
	}
	set = func(r *http.Request, value T) *http.Request {
		return r.WithContext(context.WithValue(r.Context(), key, value))
	}
	return
}

var (
	getRequestID, isRequestID, setRequestID       = createContextHelpers[string](contextKeyRequestID)
	getRequestTime, isRequestTime, setRequestTime = createContextHelpers[time.Time](contextKeyRequestTime)
	getUserID, isUserID, setUserID                = createContextHelpers[int64](contextKeyUserID)
)

type message struct {
	Message string `json:"message"`
}

type response struct {
	status  int
	payload any
}

type operationalError struct {
	status  int
	message string
}

func (e *operationalError) Error() string {
	return e.message
}

func isJWTErrorOperational(err error) bool {
	return errors.Is(err, jwt.ErrExceededExpiration) ||
		errors.Is(err, jwt.ErrInvalidFormat) ||
		errors.Is(err, jwt.ErrInvalidSignature) ||
		errors.Is(err, jwt.ErrMissingExpiration)
}

type Config struct {
	Errors Errors
}

type Errors struct {
	Internal          string `yaml:"internal"`
	NotFound          string `yaml:"notFound"`
	MethodNotAllowed  string `yaml:"methodNotAllowed"`
	BodyLimitExceeded string `yaml:"bodyLimitExceeded"`
	BodyMalformed     string `yaml:"bodyMalformed"`
	BadBodyEncoding   string `yaml:"badBodyEncoding"`
	TooManyRequests   string `yaml:"tooManyRequests"`
}

type Service struct {
	errInternal          string
	errNotFound          error
	errMethodNotAllowed  error
	errExceededBodyLimit error
	errMalformedBody     error
	errBadBodyEncoding   error
	errTooManyRequests   error
}

func NewService(cfg *Config) *Service {
	return &Service{
		errInternal:          cfg.Errors.Internal,
		errNotFound:          &operationalError{http.StatusNotFound, cfg.Errors.NotFound},
		errMethodNotAllowed:  &operationalError{http.StatusMethodNotAllowed, cfg.Errors.MethodNotAllowed},
		errExceededBodyLimit: &operationalError{http.StatusRequestEntityTooLarge, cfg.Errors.BodyLimitExceeded},
		errMalformedBody:     &operationalError{http.StatusBadRequest, cfg.Errors.BodyMalformed},
		errBadBodyEncoding:   &operationalError{http.StatusUnsupportedMediaType, cfg.Errors.BadBodyEncoding},
		errTooManyRequests:   &operationalError{http.StatusTooManyRequests, cfg.Errors.TooManyRequests},
	}
}

func (s *Service) NotFound() http.HandlerFunc {
	return s.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		err = s.errNotFound
		return
	})
}

func (s *Service) MethodNotAllowed() http.HandlerFunc {
	return s.createHandler(func(w http.ResponseWriter, r *http.Request) (res response, err error) {
		err = s.errMethodNotAllowed
		return
	})
}

func (s *Service) WithRequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Add("X-Request-ID", id)
		h.ServeHTTP(w, setRequestID(r, id))
	})
}

func (s *Service) WithRequestTime(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, setRequestTime(r, time.Now()))
	})
}

func (s *Service) WithBodyLimit(bytes int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, bytes)
			h.ServeHTTP(w, r)
		})
	}
}

func (s *Service) WithRateLimit(limiter ratelimit.Limiter) func(http.Handler) http.Handler {
	return s.createMiddleware(func(h http.Handler, w http.ResponseWriter, r *http.Request) error {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return err
		}
		if !limiter.Allow(ip) {
			return s.errTooManyRequests
		}
		h.ServeHTTP(w, r)
		return nil
	})
}

func (s *Service) decodeJSONBody(r *http.Request, v any) error {
	mime := strings.ToLower(
		strings.TrimSpace(
			strings.Split(r.Header.Get("Content-Type"), ";")[0],
		),
	)
	if mime != "application/json" {
		return s.errBadBodyEncoding
	}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	err := d.Decode(v)
	if err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError
		if errors.As(err, &syntaxErr) ||
			errors.As(err, &unmarshalTypeErr) ||
			errors.Is(err, io.ErrUnexpectedEOF) ||
			errors.Is(err, io.EOF) ||
			strings.HasPrefix(err.Error(), "json: unknown field ") {
			return s.errMalformedBody
		}
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return s.errExceededBodyLimit
		}
		return err
	}
	err = d.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return s.errMalformedBody
	}
	return nil
}

func (s *Service) reply(w http.ResponseWriter, r *http.Request, res response, err error) {
	if err != nil {
		var operationalErr *operationalError
		if errors.As(err, &operationalErr) {
			res = response{
				operationalErr.status,
				message{operationalErr.message},
			}
		} else {
			res = response{
				http.StatusInternalServerError,
				message{s.errInternal},
			}
		}
	}
	level := logrus.InfoLevel
	if res.status >= 500 {
		level = logrus.ErrorLevel
	} else if res.status >= 400 {
		level = logrus.WarnLevel
	}
	if res.payload != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(res.status)
		if encodeErr := json.NewEncoder(w).Encode(res.payload); encodeErr != nil {
			err = errors.Join(err, encodeErr)
			level = logrus.ErrorLevel
		}
	} else {
		w.WriteHeader(res.status)
	}
	fields := logrus.Fields{
		"ip":     strings.Split(r.RemoteAddr, ":")[0],
		"method": r.Method,
		"path":   r.URL.EscapedPath(),
		"status": res.status,
	}
	if isRequestID(r) {
		fields["id"] = getRequestID(r)
	}
	if isRequestTime(r) {
		fields["duration"] = time.Since(getRequestTime(r))
	}
	if isUserID(r) {
		fields["userId"] = getUserID(r)
	}
	logger := logrus.WithFields(fields)
	if level >= logrus.InfoLevel {
		logger.Logger.Out = os.Stdout
		logger.Log(level)
	} else {
		logger.Log(level, err)
	}
}

func (s *Service) createMiddleware(f func(h http.Handler, w http.ResponseWriter, r *http.Request) error) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := f(h, w, r); err != nil {
				s.reply(w, r, response{}, err)
			}
		})
	}
}

func (s *Service) createHandler(f func(w http.ResponseWriter, r *http.Request) (response, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := f(w, r)
		s.reply(w, r, res, err)
	}
}
