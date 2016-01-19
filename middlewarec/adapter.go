package middlewarec

import (
	"golang.org/x/net/context"
	"net/http"
)

type Adapter struct {
	ctx     context.Context
	handler Handler
}

func (a *Adapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTPC(a.ctx, w, r)
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
