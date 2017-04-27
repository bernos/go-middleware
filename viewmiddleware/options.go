package viewmiddleware

import (
	"html/template"
	"net/http"

	"github.com/bernos/go-middleware/resourceloadermiddleware"
)

type options struct {
	viewModelProvider func(r *http.Request) interface{}
	templateProvider  func(r *http.Request) *template.Template
	errorHandler      func(error, http.ResponseWriter, *http.Request) bool
}

func defaultOptions(t *template.Template) *options {
	return &options{
		viewModelProvider: defaultViewModelProvider,
		errorHandler:      defaultErrorHandler,
		templateProvider:  defaultTemplateProvider(t),
	}
}

func defaultErrorHandler(err error, w http.ResponseWriter, r *http.Request) bool {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return false
}

func defaultViewModelProvider(r *http.Request) interface{} {
	m, _ := resourceloadermiddleware.FromRequest(r)
	return m
}

func defaultTemplateProvider(t *template.Template) func(*http.Request) *template.Template {
	return func(r *http.Request) *template.Template {
		if t == nil {
			t = defaultTemplate
		}

		return t
	}
}

func WithErrorHandler(h func(error, http.ResponseWriter, *http.Request) bool) func(*options) {
	return func(o *options) {
		o.errorHandler = h
	}
}

func WithDefaultViewModelProvider(p func(r *http.Request) interface{}) func(*options) {
	return func(o *options) {
		o.viewModelProvider = p
	}
}

func WithDefaultTemplateProvider(p func(r *http.Request) *template.Template) func(*options) {
	return func(o *options) {
		o.templateProvider = p
	}
}

func WithDefaultTemplate(t *template.Template) func(*options) {
	return func(o *options) {
		o.templateProvider = func(r *http.Request) *template.Template {
			return t
		}
	}
}
