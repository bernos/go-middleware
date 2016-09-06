package loggingmiddleware

import (
	"net/http"
	"time"

	"github.com/bernos/go-middleware/middleware"
	"github.com/bernos/go-middleware/middlewarec"
)

func New(log func(RequestInfo)) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := WrapWriter(w)

			next.ServeHTTP(lw, r)

			log(RequestInfo{
				Request: r,
				Status:  lw.Status(),
				Latency: time.Since(start),
			})
		})
	}
}

func NewC(log func(RequestInfo)) middlewarec.Middleware {
	return middlewarec.FromMiddleware(New(log))
}

type RequestInfo struct {
	Request *http.Request
	Status  int
	Latency time.Duration
}
