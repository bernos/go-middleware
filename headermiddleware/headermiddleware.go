package headermiddleware

import (
	"net/http"

	"github.com/bernos/go-middleware/middleware"
)

func New(headers map[string]string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for k, v := range headers {
				w.Header().Set(k, v)
			}
			next.ServeHTTP(w, r)
		})
	}
}
