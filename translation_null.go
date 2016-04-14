package i18n

import (
	"fmt"
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

// IsNullTranslation checks if the translation is a NullTranslation or NullTranslation pointer
func IsNullTranslation(t Translation) bool {
	_, ok := t.(NullTranslation)
	if !ok {
		_, ok = t.(*NullTranslation)
	}
	return ok
}

var _ Translation = (*NullTranslation)(nil)
