package i18n

import (
	"html/template"

	"golang.org/x/text/language"
)

// recovery will silently swallow all unexpected panics.
func recovery() {
	recover()
}

// TemplateFuncs generates a new template function map.
func (i *I18N) TemplateFuncs() template.FuncMap {
	// Translation used in the template
	var t = i.fallback

	return template.FuncMap{
		"setlang": func(lang string) {
			if n, err := language.Parse(lang); err == nil {
				t = n
			}
		},
		"translate": func(key string) string {
			defer recovery()
			return i.translate(t, key)
		},
		"yesno": func(value bool) string {
			defer recovery()
			if value {
				return i.translate(t, "yes")
			}
			return i.translate(t, "no")
		},
	}
}

// TemplateNew is a wrapper for template.New that injects the I18N template functions.
func (i *I18N) TemplateNew(name string) *template.Template {
	return template.New(name).Funcs(i.TemplateFuncs())
}

// TemplateInject injects I18N template functions into the passed template function map.
func (i *I18N) TemplateInject(f map[string]interface{}) {
	for k, v := range i.TemplateFuncs() {
		f[k] = v
	}
}
