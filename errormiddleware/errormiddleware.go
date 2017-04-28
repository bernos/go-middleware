package errormiddleware

import (
	"fmt"
	"net/http"

	"github.com/bernos/go-middleware/middleware"
)

type Error struct {
	err    error
	status int
}

func (err *Error) Error() string {
	return err.err.Error()
}

func (err *Error) Status() int {
	return err.status
}

func NewError(err error, status int) *Error {
	return &Error{err, status}
}

func HandleErrors() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			err := FromRequest(r)

			fmt.Printf("==============================\n")
			fmt.Printf(" ERROR: %s\n", err)
			fmt.Printf("==============================\n")

			next.ServeHTTP(w, r)
		})
	}
}

func HandlerFunc(fn func(http.ResponseWriter, *http.Request) *Error) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := fn(w, r)

			if err != nil {
				r = UpdateRequest(r, err.err, err.status)
			}

			next.ServeHTTP(w, r)
		})
	}
}
