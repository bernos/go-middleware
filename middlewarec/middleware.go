package middlewarec

import (
	"golang.org/x/net/context"
	"net/http"
)

var (
	defaultHandler Handler = HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {})
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
	AsHttpHandler() http.Handler
	AsHttpHandlerWithContext(context.Context) http.Handler
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r)
}

func (h HandlerFunc) AsHttpHandler() http.Handler {
	return NewAdapter(h)
}

func (h HandlerFunc) AsHttpHandlerWithContext(c context.Context) http.Handler {
	return NewAdapterWithContext(c, h)
}

type Middleware func(Handler) Handler

type MiddlewareReducer func(Middleware, Middleware) Middleware

func (m Middleware) Compose(next Middleware) Middleware {
	return Compose(m, next)
}

func (m Middleware) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) {
	m(defaultHandler).ServeHTTP(c, w, r)
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

type Adapter struct {
	ctx     context.Context
	handler Handler
}

func (a *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(a.ctx, w, r)
}

func NewAdapter(h Handler) *Adapter {
	return NewAdapterWithContext(context.Background(), h)
}

func NewAdapterWithContext(c context.Context, h Handler) *Adapter {
	return &Adapter{
		ctx:     c,
		handler: h,
	}
}
