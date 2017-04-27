package viewmiddleware

import (
	"html/template"
	"net/http"

	"github.com/bernos/go-middleware/bodyparsermiddleware"
	"github.com/bernos/go-middleware/middleware"
	"github.com/bernos/go-middleware/resourceloadermiddleware"
)

var defaultTemplate = template.Must(template.New("_default").Parse(`This is the default template`))

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

func RenderView(defaultTemplate *template.Template, options ...func(*options)) middleware.Middleware {
	cfg := defaultOptions(defaultTemplate)

	for _, o := range options {
		o(cfg)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := GetTemplate(r, cfg.templateProvider(r))
			m := GetViewModel(r, cfg.viewModelProvider(r))

			shouldContinue := true

			vm := struct {
				Model           interface{}
				Error           error
				ValidationError error
			}{
				Model:           m,
				Error:           nil,
				ValidationError: bodyparsermiddleware.Validate(r),
			}

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

type View struct {
	Template *template.Template
	Model    interface{}
}

func BuildView(fn func(*http.Request) *View) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			view := fn(r)

			if view.Template != nil {
				r = RequestWithTemplate(r, view.Template)
			}

			if view.Model != nil {
				r = RequestWithViewModel(r, view.Model)
			}

			next.ServeHTTP(w, r)
		})
	}
}
