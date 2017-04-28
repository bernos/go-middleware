package errormiddleware

import (
	"net/http"

	"github.com/bernos/go-middleware/middleware"
)

type Error struct {
	error
	status int
}

func (err *Error) Status() int {
	return err.status
}

func NewError(err error, status int) *Error {
	return &Error{err, status}
}

type HTTPError interface {
	error
	Status() int
}

func HandleErrors(options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions()

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := FromRequest(r)

			if err != nil {
				w.Header().Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(getErrorStatus(err))

				cfg.errorHandler(err, w, r)
			}

			if cfg.continueOnError {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func HandlerFunc(fn func(http.ResponseWriter, *http.Request) error) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			err := fn(w, r)

			if err != nil {
				r = UpdateRequest(r, err, getErrorStatus(err))
			}

			next.ServeHTTP(w, r)
		})
	}
}

func getErrorStatus(err error) int {
	httpErr, ok := err.(HTTPError)

	if ok {
		return httpErr.Status()
	}

	return http.StatusInternalServerError
}
