package errormiddleware

import (
	"fmt"
	"net/http"
)

type options struct {
	errorHandler func(error, http.ResponseWriter, *http.Request) bool
}

func defaultOptions() *options {
	return &options{
		errorHandler: defaultErrorHandler,
	}
}

func defaultErrorHandler(e error, w http.ResponseWriter, r *http.Request) bool {
	fmt.Fprintln(w, e.Error())
	return false
}

func WithErrorHandler(fn func(error, http.ResponseWriter, *http.Request) bool) func(*options) {
	return func(o *options) {
		o.errorHandler = fn
	}
}

func ContinueOnError() func(*options) {
	return func(o *options) {
		o.errorHandler = func(e error, w http.ResponseWriter, r *http.Request) bool {
			fmt.Fprintln(w, e.Error())
			return true
		}
	}
}
