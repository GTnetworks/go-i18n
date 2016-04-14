// Package i18n provides internationalization (i18n) for Go projects and various template engines.
package i18n

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

var nullTranslation = NullTranslation{}

type Internationalization struct {
	translation map[language.Tag]Translation
	fallback    language.Tag
}

// New starts a new internationalization instance.
func New(fallback language.Tag) (*Internationalization, error) {
	i := new(Internationalization)
	i.translation = make(map[language.Tag]Translation)
	i.fallback = fallback
	return i, nil
}

// Add a translation
func (i *Internationalization) Add(t Translation) {
	i.translation[t.Tag()] = t
}

// Accept parses a HTTP Accept-Language header and returns a matching translation. May return a
// NullTranslation if there is no matching language and the fallback language is also not available.
func (i *Internationalization) Accept(accept string) Translation {
	var (
		s = i.Supported()
		m = language.NewMatcher(append([]language.Tag{i.fallback}, s...))
	)
	t, _, err := language.ParseAcceptLanguage(accept)
	if err != nil {
		return i.translation[i.fallback]
	}
	// We ignore the error: the fallback language will be selected for t == nil.
	tag, _, _ := m.Match(t...)
	if tr, ok := i.translation[tag]; ok {
		return tr
	}
	return nullTranslation
}

// Default returns the default (fallback) tag.
func (i *Internationalization) Default() language.Tag { return i.fallback }

// Languages returns a slice of supported languages, in their native translation.
func (i *Internationalization) Languages() []string {
	var (
		s = i.Supported()
		l = make([]string, 0)
	)

	for _, t := range s {
		l = append(l, display.Self.Name(t))
	}

	return l
}

// Supported returns a slice of supported language tags.
func (i *Internationalization) Supported() []language.Tag {
	var s = make([]language.Tag, 0)
	for t := range i.translation {
		s = append(s, t)
	}
	return s
}

// Translate a phrase based on the selected language.
func (i *Internationalization) Translate(lang, key string) string {
	t, err := language.Parse(lang)
	if err != nil {
		t = i.fallback
	}
	return i.translate(t, key)
}

func (i *Internationalization) translate(tag language.Tag, key string) string {
	t, ok := i.translation[tag]
	if !ok {
		t, ok = i.translation[i.fallback]
	}
	if !ok {
		t = nullTranslation
	}
	return t.Get(key)
}
