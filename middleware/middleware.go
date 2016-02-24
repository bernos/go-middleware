package middleware

import (
	"net/http"
)

var (
	defaultHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

// Middleware is a function that wraps a regular http.Handler.
type Middleware func(http.Handler) http.Handler

// MiddlewareReducer is a function that creates a new Middleware from
// tow other Middleware. The middleware composition function is an
// example
type MiddlewareReducer func(Middleware, Middleware) Middleware

// Compose creates a new middleware by composing the method receiver
// with a 'next' middleware
func (m Middleware) Compose(next Middleware) Middleware {
	return Compose(m, next)
}

// ServeHTTP allows a Middleware to be used as an http.Handler.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m(defaultHandler).ServeHTTP(w, r)
}

// Then is just a proxy to calling the middleware function directly with
// the provided handler. It is just some syntactic sugar that allows for
// method chaining
func (m Middleware) Then(handler http.Handler) http.Handler {
	return m(handler)
}

// Id returns a middleware that will return the provided http.Handler
func Id() Middleware {
	return Middleware(func(next http.Handler) http.Handler {
		return next
	})
}

// Compose chains two Middleware, f and g, into a new Middleware that
// returns the result of passing the output of g to f
func Compose(f Middleware, g Middleware) Middleware {
	return Middleware(func(h http.Handler) http.Handler {
		return f(g(h))
	})
}

// ComposeAll chains multiple Middleware, by folding Compose over a
// list of Middleware
func ComposeAll(middlewares ...Middleware) Middleware {
	if len(middlewares) == 0 {
		middlewares = append(middlewares, Id())
	}
	return fold(Compose, middlewares[0], middlewares[1:])
}

func fold(f MiddlewareReducer, x Middleware, xs []Middleware) Middleware {
	for _, m := range xs {
		x = f(x, m)
	}
	return x
}
