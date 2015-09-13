package middleware

import (
	"net/http"
)

type MiddlewareFunc func(http.Handler) http.Handler

func (m MiddlewareFunc) Compose(next MiddlewareFunc) MiddlewareFunc {
	return compose(m, next)
}

func Compose(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return reduce(compose, middlewares)
}

func compose(f MiddlewareFunc, g MiddlewareFunc) MiddlewareFunc {
	return MiddlewareFunc(func(h http.Handler) http.Handler {
		return f(g(h))
	})
}

func fold(f func(MiddlewareFunc, MiddlewareFunc) MiddlewareFunc, x MiddlewareFunc, xs []MiddlewareFunc) MiddlewareFunc {
	for _, m := range xs {
		x = f(x, m)
	}
	return x
}

func reduce(f func(MiddlewareFunc, MiddlewareFunc) MiddlewareFunc, xs []MiddlewareFunc) MiddlewareFunc {
	// Use the unit MiddlewareFunc if the list is empty
	if len(xs) == 0 {
		xs = append(xs, unit())
	}

	// xs[1:] will return an empty slice in the case of a single element
	// MiddlewareFunc slice.
	return fold(f, xs[0], xs[1:])
}

func unit() MiddlewareFunc {
	return MiddlewareFunc(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
}
