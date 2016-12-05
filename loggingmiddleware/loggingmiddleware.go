package loggingmiddleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bernos/go-middleware/middleware"
	"github.com/bernos/go-middleware/middlewarec"
)

func New(options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions()

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			lw := &responseWriterWrapper{w, false, 0}

			next.ServeHTTP(lw, r)

			cfg.logger(RequestInfo{
				Request: r,
				Status:  lw.Status(),
				Latency: time.Since(start),
			})
		})
	}
}

func NewC(options ...func(*options)) middlewarec.Middleware {
	return middlewarec.FromMiddleware(New(options...))
}

type RequestInfo struct {
	Request *http.Request
	Status  int
	Latency time.Duration
}

func (r RequestInfo) AsMap() map[string]interface{} {
	fields := map[string]interface{}{
		"http_host": r.Request.Host,
		"method":    r.Request.Method,
		"uri":       r.Request.RequestURI,
		"remote":    r.Request.RemoteAddr,
		"status":    r.Status,
		"latency":   r.Latency,
	}

	for k, v := range r.Request.Header {
		name := strings.ToLower(fmt.Sprintf("http_%s", strings.Replace(k, "-", "_", -1)))
		fields[name] = v
	}

	return fields
}

type responseWriterWrapper struct {
	http.ResponseWriter
	wroteHeader bool
	status      int
}

func (w *responseWriterWrapper) Status() int {
	return w.status
}

func (w *responseWriterWrapper) WriteHeader(status int) {
	if !w.wroteHeader {
		w.status = status
		w.wroteHeader = true
	}

	w.ResponseWriter.WriteHeader(status)
}
