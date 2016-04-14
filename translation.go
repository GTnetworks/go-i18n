package i18n

import "golang.org/x/text/language"

type Translation interface {
	Format(string, ...interface{}) string
	Formats([]string, ...interface{}) string
	Get(string) string
	Has(string) bool
	Len() int
	Tag() language.Tag
}
