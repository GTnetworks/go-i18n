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

func NewTranslationFile(lang language.Tag, name string) (t TranslationMap, err error) {
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
		return ParseJSON(lang, b)
	case ".tml", ".toml", ".conf":
		return ParseTOML(lang, b)
	case ".yml", ".yaml":
		return ParseYAML(lang, b)
	default:
		err = ErrUnknownFormat
	}

	return
}

// ParseJSON takes a JSON blob and tries to decode it as a TranslationMap.
func ParseJSON(lang language.Tag, b []byte) (t TranslationMap, err error) {
	var (
		k = make(map[string]interface{})
		m map[string]string
	)

	if err = json.Unmarshal(b, &k); err == nil {
		if m, err = flatten(k); err == nil {
			return NewMap(lang, m), nil
		}
	}

	return
}

// ParseTOML takes a TOML blob and tries to decode it as a TranslationMap.
func ParseTOML(lang language.Tag, b []byte) (t TranslationMap, err error) {
	var (
		k = make(map[string]interface{})
		m map[string]string
	)

	if err = toml.Unmarshal(b, &k); err == nil {
		if m, err = flatten(k); err == nil {
			return NewMap(lang, m), nil
		}
	}

	return
}

// Parse Ruby compatible i18n translations file, which is a (tree based)
// key-value structure with the root keys as the language, see
// http://guides.rubyonrails.org/i18n.html
func ParseYAML(lang language.Tag, b []byte) (t TranslationMap, err error) {
	var (
		k = make(map[string]map[string]interface{})
		m map[string]string
	)

	if err = yaml.Unmarshal(b, &k); err == nil {
		// Test if the full language is in the map, ie `en_US`.
		if in, ok := k[lang.String()]; ok {
			if m, err = flatten(in); err == nil {
				return NewMap(lang, m), nil
			}
		}
		// Test if the base language is in the map, ie `en`.
		base, _ := lang.Base()
		if in, ok := k[base.String()]; ok {
			if m, err = flatten(in); err == nil {
				return NewMap(lang, m), nil
			}
		}
		err = fmt.Errorf("i18n: language %q not found", lang)
	}

	return
}
