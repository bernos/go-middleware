package resourceloadermiddleware

import "net/http"

type options struct {
	errorHandler    func(error, http.ResponseWriter, *http.Request) bool
	notFoundHandler func(http.ResponseWriter, *http.Request) bool
}

func defaultOptions() *options {
	o := &options{
		errorHandler:    defaultErrorHandler,
		notFoundHandler: defaultNotFoundHandler,
	}

	return o
}

// WithErrorHandler sets the error handler func. An error handler is a func that
// handles errors returned by a ResourceLoader. It should inspect the error,
// update the ResponseWriter appropriately, and finally return a bool indicating
// whether request processing should continue, or terminate immediately.
// A return value of `true` will result in request processing continuing through
// the middleware chain. Returning `false` will terminate request processing and
// send the reponse immediately
func WithErrorHandler(h func(error, http.ResponseWriter, *http.Request) bool) func(*options) {
	return func(o *options) {
		o.errorHandler = h
	}
}

// WithNotFoundHandler set the not found handler func. A not found handler is a
// func that handles cases when a ResourceLoader cannot find a resource. It
// should update the ResponseWriter appropriately and return a bool indicating
// whether request processing should continue. A return value of `true` will
// result in request processing continuing through the middleware chain.
// Returning `false` will terminate request processing and send the response
// immediately
func WithNotFoundHandler(h func(http.ResponseWriter, *http.Request) bool) func(*options) {
	return func(o *options) {
		o.notFoundHandler = h
	}
}

func defaultErrorHandler(err error, w http.ResponseWriter, r *http.Request) bool {
	http.Error(w, "Failed to load resource", http.StatusInternalServerError)
	return false
}

func defaultNotFoundHandler(w http.ResponseWriter, r *http.Request) bool {
	http.NotFound(w, r)
	return true
}
