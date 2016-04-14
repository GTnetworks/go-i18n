package i18n

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

type TranslationMap struct {
	keys map[string]string
	lang language.Tag
}

func NewMap(lang language.Tag, keys map[string]string) TranslationMap {
	return TranslationMap{keys, lang}
}

func (t TranslationMap) Format(key string, args ...interface{}) string {
	return fmt.Sprintf(key, args...)
}

func (t TranslationMap) Formats(keys []string, args ...interface{}) string {
	if keys == nil {
		return ""
	}

	var out = make([]string, len(keys))
	for i, key := range keys {
		out[i] = t.Get(key)
	}

	return fmt.Sprintf(strings.Join(out, " "), args...)
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

var _ Translation = (*TranslationMap)(nil)
