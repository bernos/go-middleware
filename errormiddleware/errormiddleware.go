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

func HandleErrors() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			err := FromRequest(r)

			fmt.Printf("==============================\n")
			fmt.Printf(" ERROR: %s\n", err)
			fmt.Printf("==============================\n")
		})
	}
}
