package i18n

import (
	"io"

	"golang.org/x/text/language"
)

type Translater interface {
	Get(string) string
	Has(string) bool
	Len() int
	Tag() language.Tag
}

// Formatter partially implements the fmt package.
type Formatter interface {
	Fprint(io.Writer, ...interface{}) (int, error)
	Fprintf(io.Writer, string, ...interface{}) (int, error)
	Fprintln(io.Writer, ...interface{}) (int, error)
	Sprint(...interface{}) string
	Sprintf(string, ...interface{}) string
	Sprintln(...interface{}) string
}

type Translation interface {
	Translater
	Formatter
}
