package viewmiddleware

import (
	"html/template"
	"net/http"

	"github.com/bernos/go-middleware/middleware"
	"github.com/bernos/go-middleware/resourceloadermiddleware"
)

type options struct {
	viewModelProvider func(r *http.Request) (interface{}, bool)
	errorHandler      func(error, http.ResponseWriter, *http.Request) bool
}

func defaultOptions() *options {
	return &options{
		viewModelProvider: defaultViewModelProvider,
		errorHandler:      defaultErrorHandler,
	}
}

func defaultErrorHandler(err error, w http.ResponseWriter, r *http.Request) bool {
	http.Error(w, "Failed to render template", http.StatusInternalServerError)
	return false
}

func defaultViewModelProvider(r *http.Request) (interface{}, bool) {
	return resourceloadermiddleware.FromRequest(r)
}

func WithErrorHandler(h func(error, http.ResponseWriter, *http.Request) bool) func(*options) {
	return func(o *options) {
		o.errorHandler = h
	}
}

func WithViewModelProvider(p func(r *http.Request) (interface{}, bool)) func(*options) {
	return func(o *options) {
		o.viewModelProvider = p
	}
}

func View(t *template.Template, options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions()

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			shouldContinue := true

			vm, _ := cfg.viewModelProvider(r)

			err := t.Execute(w, vm)

			if err != nil {
				shouldContinue = cfg.errorHandler(err, w, r)
			}

			if shouldContinue {
				next.ServeHTTP(w, r)
			}
		})
	}
}
