package i18n

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/language"
)

// NullTranslation provides a dummy Translation that just returns the input keys.
type NullTranslation struct{}

func (t NullTranslation) Format(key string, args ...interface{}) string {
	return fmt.Sprintf(key, args...)
}

func (t NullTranslation) Formats(keys []string, args ...interface{}) string {
	return fmt.Sprintf(strings.Join(keys, " "), args...)
}

func (t NullTranslation) Get(key string) string { return key }
func (t NullTranslation) Has(key string) bool   { return false }
func (t NullTranslation) Len() int              { return 0 }
func (t NullTranslation) Tag() language.Tag     { return language.Und }

func (t NullTranslation) Fprint(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprint(w, a...)
}

func (t NullTranslation) Fprintf(w io.Writer, key string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, key, a...)
}

func (t NullTranslation) Fprintln(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprintln(w, a...)
}

func (t NullTranslation) Sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}

func (t NullTranslation) Sprintf(key string, a ...interface{}) string {
	return fmt.Sprintf(key, a...)
}

func (t NullTranslation) Sprintln(a ...interface{}) string { return fmt.Sprintln(a...) }

// IsNullTranslation checks if the translation is a NullTranslation or NullTranslation pointer
func IsNullTranslation(t Translation) bool {
	_, ok := t.(NullTranslation)
	if !ok {
		_, ok = t.(*NullTranslation)
	}
	return ok
}

var _ Translation = (*NullTranslation)(nil)
