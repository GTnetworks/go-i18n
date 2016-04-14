package i18n

import (
	"bytes"
	"fmt"
	"io"

	"golang.org/x/text/language"
)

type TranslationMap struct {
	keys map[string]string
	lang language.Tag
}

func NewMap(lang language.Tag, keys map[string]string) TranslationMap {
	return TranslationMap{keys, lang}
}

func (t TranslationMap) Get(key string) string {
	if out, ok := t.keys[key]; ok {
		return out
	}
	return key
}

func (t TranslationMap) Has(key string) bool {
	return t.keys[key] != ""
}

func (t TranslationMap) Len() int {
	return len(t.keys)
}

func (t TranslationMap) Tag() language.Tag {
	return t.lang
}

func (t TranslationMap) Fprint(w io.Writer, a ...interface{}) (int, error) {
	var out = make([]interface{}, len(a))

	for i, arg := range a {
		if key, ok := arg.(string); ok {
			out[i] = t.Get(key)
		} else {
			out[i] = arg
		}
	}

	return fmt.Fprint(w, out...)
}

func (t TranslationMap) Fprintf(w io.Writer, key string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, t.Get(key), a...)
}

func (t TranslationMap) Fprintln(w io.Writer, a ...interface{}) (int, error) {
	var out = make([]interface{}, len(a))

	for i, arg := range a {
		if key, ok := arg.(string); ok {
			out[i] = t.Get(key)
		} else {
			out[i] = arg
		}
	}

	return fmt.Fprintln(w, out...)
}

func (t TranslationMap) Sprint(a ...interface{}) string {
	b := new(bytes.Buffer)
	t.Fprint(b, a...)
	return b.String()
}

func (t TranslationMap) Sprintf(key string, a ...interface{}) string {
	b := new(bytes.Buffer)
	t.Fprintf(b, key, a...)
	return b.String()
}

func (t TranslationMap) Sprintln(a ...interface{}) string {
	b := new(bytes.Buffer)
	t.Fprintln(b, a...)
	return b.String()
}

var _ Translation = (*TranslationMap)(nil)
