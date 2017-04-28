package errormiddleware

import (
	"fmt"
	"html/template"
	"net/http"
)

type options struct {
	errorHandler    func(error, http.ResponseWriter, *http.Request)
	continueOnError bool
}

func defaultOptions() *options {
	return &options{
		errorHandler:    defaultErrorHandler,
		continueOnError: false,
	}
}

func defaultErrorHandler(e error, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, e.Error())
}

func WithErrorHandler(fn func(error, http.ResponseWriter, *http.Request)) func(*options) {
	return func(o *options) {
		o.errorHandler = fn
	}
}

func ContinueOnError() func(*options) {
	return func(o *options) {
		o.continueOnError = true
	}
}

func HaltOnError() func(*options) {
	return func(o *options) {
		o.continueOnError = false
	}
}

func WithErrorTemplate(t *template.Template) func(*options) {
	return func(o *options) {
		o.errorHandler = func(e error, w http.ResponseWriter, r *http.Request) {
			t.Execute(w, e)
		}
	}
}
