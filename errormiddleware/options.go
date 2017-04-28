package errormiddleware

import "net/http"

type options struct {
	errorHandler func(*Error, http.ResponseWriter, *http.Request) bool
}

func defaultOptions() *options {
	return &options{
		errorHandler: defaultErrorHandler,
	}
}

func defaultErrorHandler(e *Error, w http.ResponseWriter, r *http.Request) bool {
	http.Error(w, e.Error(), e.Status())
	return false
}

func WithErrorHandler(fn func(*Error, http.ResponseWriter, *http.Request) bool) func(*options) {
	return func(o *options) {
		o.errorHandler = fn
	}
}

func ContinueOnError() func(*options) {
	return func(o *options) {
		o.errorHandler = func(e *Error, w http.ResponseWriter, r *http.Request) bool {
			http.Error(w, e.Error(), e.Status())
			return true
		}
	}
}
