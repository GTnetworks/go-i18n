package i18n

import "golang.org/x/text/language"

type NullTranslation struct {
	TranslationMap
}

func (t NullTranslation) Tag() language.Tag { return language.Und }

// IsNullTranslation checks if the translation is a NullTranslation or NullTranslation pointer
func IsNullTranslation(t Translation) bool {
	_, ok := t.(NullTranslation)
	if !ok {
		_, ok = t.(*NullTranslation)
	}
	return ok
}

var _ Translation = (*NullTranslation)(nil)
