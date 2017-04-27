package viewmiddleware

import (
	"context"
	"html/template"
	"net/http"
)

type key int

const (
	viewModelKey key = 0
	templateKey  key = 1
)

func RequestWithViewModel(r *http.Request, m interface{}) *http.Request {
	// return r.WithContext(context.WithValue(r.Context(), viewModelKey, m))
	*r = *(r.WithContext(context.WithValue(r.Context(), viewModelKey, m)))
	return r
}

func RequestWithTemplate(r *http.Request, t *template.Template) *http.Request {
	// return r.WithContext(context.WithValue(r.Context(), templateKey, t))
	*r = *(r.WithContext(context.WithValue(r.Context(), templateKey, t)))
	return r
}

func GetViewModel(r *http.Request, defaultModel interface{}) interface{} {
	model := r.Context().Value(viewModelKey)

	if model == nil {
		model = defaultModel
	}

	return model
}

func GetTemplate(r *http.Request, defaultTemplate *template.Template) *template.Template {
	t, ok := r.Context().Value(templateKey).(*template.Template)

	if !ok {
		t = defaultTemplate
	}

	return t
}

type Request struct {
	*http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{r}
}

func (r *Request) WithViewModel(m interface{}) *Request {
	return &Request{RequestWithViewModel(r.HTTPRequest(), m)}
}

func (r *Request) WithTemplate(t *template.Template) *Request {
	return &Request{RequestWithTemplate(r.HTTPRequest(), t)}
}

func (r *Request) GetViewModel(defaultValue interface{}) interface{} {
	return GetViewModel(r.HTTPRequest(), defaultValue)
}

func (r *Request) GetTemplate(defaultTemplate *template.Template) *template.Template {
	return GetTemplate(r.HTTPRequest(), defaultTemplate)
}

func (r *Request) HTTPRequest() *http.Request {
	return r.Request
}
