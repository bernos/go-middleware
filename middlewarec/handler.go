package middlewarec

import (
	"golang.org/x/net/context"
	"net/http"
)

type Handler interface {
	http.Handler
	ServeHTTPC(context.Context, http.ResponseWriter, *http.Request)
	AsHttpHandler() http.Handler
	AsHttpHandlerWithContext(context.Context) http.Handler
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.AsHttpHandler().ServeHTTP(w, r)
}

func (h HandlerFunc) ServeHTTPC(c context.Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r)
}

func (h HandlerFunc) AsHttpHandler() http.Handler {
	return NewAdapter(h)
}

func (h HandlerFunc) AsHttpHandlerWithContext(c context.Context) http.Handler {
	return NewAdapterWithContext(c, h)
}

func WrapHTTPHandler(h http.Handler) Handler {
	return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}
