package errormiddleware

import (
	"context"
	"net/http"
)

type key int

const errorKey key = 1

func NewContext(parent context.Context, err *Error) context.Context {
	return context.WithValue(parent, errorKey, err)
}

func UpdateRequest(r *http.Request, err error, status int) *http.Request {
	return r.WithContext(NewContext(r.Context(), &Error{err, status}))
}

func FromContext(ctx context.Context) *Error {
	err, ok := ctx.Value(errorKey).(*Error)

	if !ok {
		return nil
	}

	return err
}

func FromRequest(r *http.Request) *Error {
	return FromContext(r.Context())
}
