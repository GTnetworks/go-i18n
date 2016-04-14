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
	ErrUnknownFormat = errors.New("i18n: unknown file format")
)

type TranslationFile struct {
	keys map[string]string
	tag  language.Tag
}

func NewTranslationFile(name, lang string) (t *TranslationFile, err error) {
	t = new(TranslationFile)
	if t.tag, err = language.Parse(lang); err != nil {
		return
	}

	var (
		f *os.File
		b []byte
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
		err = t.parseJSON(lang, b)
	case ".tml", ".toml", ".conf":
		err = t.parseTOML(lang, b)
	case ".yml", ".yaml":
		err = t.parseYAML(lang, b)
	default:
		return nil, ErrUnknownFormat
	}

	return
}

func (t *TranslationFile) parseJSON(lang string, b []byte) (err error) {
	var k = make(map[string]interface{})

	if err = json.Unmarshal(b, &k); err == nil {
		t.keys, err = flatten(k)
	}

	return
}

func (t *TranslationFile) parseTOML(lang string, b []byte) (err error) {
	var k = make(map[string]interface{})

	if err = toml.Unmarshal(b, &k); err == nil {
		t.keys, err = flatten(k)
	}

	return
}

// Parse Ruby compatible i18n translations file, which is a (tree based)
// key-value structure with the root keys as the language, see
// http://guides.rubyonrails.org/i18n.html
func (t *TranslationFile) parseYAML(lang string, b []byte) (err error) {
	var k = make(map[string]map[string]interface{})

	if err = yaml.Unmarshal(b, &k); err == nil {
		if in, ok := k[lang]; ok {
			t.keys, err = flatten(in)
		} else {
			return fmt.Errorf("i18n: language %q not found", lang)
		}
	}

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
