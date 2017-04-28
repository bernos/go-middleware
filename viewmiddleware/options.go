package viewmiddleware

import "html/template"

type options struct {
	defaultTemplate *template.Template
}

func defaultOptions() *options {
	return &options{}
}

func WithDefaultTemplate(t *template.Template) func(*options) {
	return func(o *options) {
		o.defaultTemplate = t
	}
}
