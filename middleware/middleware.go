package middleware

import (
	"net/http"
)

var (
	defaultHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

type MiddlewareFunc func(http.Handler) http.Handler

type MiddlewareReducer func(MiddlewareFunc, MiddlewareFunc) MiddlewareFunc

func (m MiddlewareFunc) Compose(next MiddlewareFunc) MiddlewareFunc {
	return Compose(m, next)
}

func Id() MiddlewareFunc {
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
}

func (m MiddlewareFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m(defaultHandler).ServeHTTP(w, r)
}

func Compose(f MiddlewareFunc, g MiddlewareFunc) MiddlewareFunc {
	return MiddlewareFunc(func(h http.Handler) http.Handler {
		return f(g(h))
	})
}

func ComposeAll(middlewares ...MiddlewareFunc) MiddlewareFunc {
	if len(middlewares) == 0 {
		middlewares = append(middlewares, Id())
	}
	return fold(Compose, middlewares[0], middlewares[1:])
}

func fold(f MiddlewareReducer, x MiddlewareFunc, xs []MiddlewareFunc) MiddlewareFunc {
	for _, m := range xs {
		x = f(x, m)
	}
	return x
}
