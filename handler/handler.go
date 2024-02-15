package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/cyberwlodarczyk/auth/jwt"
	"github.com/google/uuid"
)

type contextKey int

const (
	contextKeyRequestId contextKey = iota
	contextKeyRequestTimer
	contextKeyRequestStatus
	contextKeyRequestError
	contextKeyUserId
)

func createContextHelpers[T any](key contextKey) (getter func(*http.Request) T, setter func(*http.Request, T) *http.Request) {
	getter = func(r *http.Request) T {
		return r.Context().Value(key).(T)
	}
	setter = func(r *http.Request, value T) *http.Request {
		return r.WithContext(context.WithValue(r.Context(), key, value))
	}
	return
}

var (
	getRequestID, setRequestID         = createContextHelpers[string](contextKeyRequestId)
	getRequestTime, setRequestTime     = createContextHelpers[time.Time](contextKeyRequestTimer)
	getRequestStatus, setRequestStatus = createContextHelpers[int](contextKeyRequestStatus)
	getRequestError, setRequestError   = createContextHelpers[error](contextKeyRequestError)
	getUserID, setUserID               = createContextHelpers[int64](contextKeyUserId)
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

var (
	errExceededBodyLimit = &operationalError{http.StatusRequestEntityTooLarge, "request body limit has been exceeded"}
	errMalformedBody     = &operationalError{http.StatusBadRequest, "request body is invalid or malformed"}
	errBadBodyEncoding   = &operationalError{http.StatusUnsupportedMediaType, "request body encoding should be json"}
)

func decodeJSONBody(r *http.Request, v any) error {
	mime := strings.ToLower(
		strings.TrimSpace(
			strings.Split(r.Header.Get("Content-Type"), ";")[0],
		),
	)
	if mime != "application/json" {
		return errBadBodyEncoding
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
			return errMalformedBody
		}
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return errExceededBodyLimit
		}
		return err
	}
	err = d.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errMalformedBody
	}
	return nil
}

func createHandler(f func(w http.ResponseWriter, r *http.Request) (response, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := f(w, r)
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
					message{"something went wrong"},
				}
			}
		}
		if res.status != 0 {
			if res.payload != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(res.status)
				err = errors.Join(err, json.NewEncoder(w).Encode(res.payload))
			} else {
				w.WriteHeader(res.status)
			}
			setRequestStatus(r, res.status)
			setRequestError(r, err)
		}
	})
}

func isJWTErrorOperational(err error) bool {
	return errors.Is(err, jwt.ErrExceededExpiration) ||
		errors.Is(err, jwt.ErrInvalidFormat) ||
		errors.Is(err, jwt.ErrInvalidSignature) ||
		errors.Is(err, jwt.ErrMissingExpiration)
}

type RequestLog struct {
	Addr     string
	Method   string
	Path     string
	Status   int
	Id       string
	Error    error
	Time     time.Time
	Duration time.Duration
}

func WithRequestLogger(f func(RequestLog)) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			log := RequestLog{
				Addr:   r.RemoteAddr,
				Method: r.Method,
				Path:   r.URL.EscapedPath(),
				Status: getRequestStatus(r),
				Id:     getRequestID(r),
				Error:  getRequestError(r),
				Time:   getRequestTime(r),
			}
			log.Duration = time.Since(log.Time)
			f(log)
		})
	}
}

func WithRequestTime() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, setRequestTime(r, time.Now()))
		})
	}
}

func WithRequestID() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := uuid.New().String()
			w.Header().Add("X-Request-ID", id)
			h.ServeHTTP(w, setRequestID(r, id))
		})
	}
}

func WithBodyLimit(bytes int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, bytes)
			h.ServeHTTP(w, r)
		})
	}
}
