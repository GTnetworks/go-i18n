package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"

	"github.com/naoina/toml"
	"gopkg.in/yaml.v2"
)

var (
	ErrUnknownFormat            = errors.New("i18n: unknown file format")
	ErrNoTranslationFile        = errors.New("i18n: no translation found")
	ErrMultipleTranslationFiles = errors.New("i18n: multiple translations found")
)

type TranslationFile struct {
	keys map[string]string
	tag  language.Tag
}

func NewTranslationFile(name string) (t *TranslationFile, err error) {
	var (
		f   *os.File
		b   []byte
		k   = make(map[string]map[string]string)
		tag string
	)
	if f, err = os.Open(name); err != nil {
		return
	}
	defer f.Close()
	if b, err = ioutil.ReadAll(f); err != nil {
		return
	}

	switch strings.ToLower(filepath.Ext(name)) {
	case ".js", ".json":
		if err = json.Unmarshal(b, &k); err != nil {
			return
		}
	case ".tml", ".toml", ".conf":
		if err = toml.Unmarshal(b, &k); err != nil {
			return
		}
	case ".yml", ".yaml":
		if err = yaml.Unmarshal(b, &k); err != nil {
			return
		}
	default:
		return nil, ErrUnknownFormat
	}

	if len(k) == 0 {
		return nil, ErrNoTranslationFile
	}

	if len(k) > 1 {
		return nil, ErrMultipleTranslationFiles
	}

	for tag = range k {
	}
	if t.tag, err = language.Parse(tag); err != nil {
		return
	}
	t.keys = k[tag]

	return
}

func (t *TranslationFile) Format(key string, args ...interface{}) string {
	return fmt.Sprintf(t.Get(key), args...)
}

func (t *TranslationFile) Formats(keys []string, args ...interface{}) string {
	if keys == nil {
		return ""
	}

	var out = make([]string, len(keys))
	for i, key := range keys {
		out[i] = t.Get(key)
	}

	return fmt.Sprintf(strings.Join(out, " "), args...)
}

func (t *TranslationFile) Get(key string) string {
	if out, ok := t.keys[key]; ok {
		return out
	}
	return key
}

func (t *TranslationFile) Has(key string) bool {
	return t.keys[key] != ""
}

func (t *TranslationFile) Len() int {
	return len(t.keys)
}

func (t *TranslationFile) Tag() language.Tag {
	return t.tag
}

var _ Translation = (*TranslationFile)(nil)
