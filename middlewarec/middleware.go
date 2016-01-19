package middlewarec

import (
	"golang.org/x/net/context"
	"net/http"
)

var (
	defaultHandler Handler = HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {})
)

type Middleware func(Handler) Handler

type MiddlewareReducer func(Middleware, Middleware) Middleware

func (m Middleware) Compose(next Middleware) Middleware {
	return Compose(m, next)
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.AsHttpHandler().ServeHTTP(w, r)
}

func (m Middleware) ServeHTTPC(c context.Context, w http.ResponseWriter, r *http.Request) {
	m(defaultHandler).ServeHTTPC(c, w, r)
}

func (m Middleware) AsHttpHandler() http.Handler {
	return NewAdapter(m)
}

func (m Middleware) AsHttpHandlerWithContext(c context.Context) http.Handler {
	return NewAdapterWithContext(c, m)
}

func Id() Middleware {
	return Middleware(func(next Handler) Handler {
		return next
	})
}

func Compose(f Middleware, g Middleware) Middleware {
	return Middleware(func(h Handler) Handler {
		return f(g(h))
	})
}

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
