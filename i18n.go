package i18n

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	golanguage "golang.org/x/text/language"
	"gopkg.in/yaml.v2"

	"github.com/flosch/pongo2"
)

type dictionary map[string]string

type language map[string]dictionary

type Translations struct {
	language language
	configs  []string
}

type Config struct {
	Language map[string]dictionary `yaml:"lang"`
	Include  []string
}

func New(name string) (*Translations, error) {
	t := &Translations{
		language: make(language),
	}
	return t, t.Load(name)
}

func split(s string, c byte) (head, tail string) {
	if i := strings.IndexByte(s, c); i >= 0 {
		return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+1:])
	}
	return strings.TrimSpace(s), ""
}

type acceptLanguage struct {
	t string
	q float32
}

type acceptLanguages []acceptLanguage

func (a acceptLanguages) Len() int           { return len(a) }
func (a acceptLanguages) Less(i, j int) bool { return a[i].q > a[j].q }
func (a acceptLanguages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// AcceptLanguage parses the Accept-Language HTTP headers and returns a sorted
// list of supported languages.
func (t *Translations) AcceptLanguage(s string) []string {
	tags, quals, err := golanguage.ParseAcceptLanguage(s)
	if err != nil {
		return nil
	}

	var accepted = make(acceptLanguages, 0)
	for i, tag := range tags {
		// Inspect the full language of the language tag.
		if _, ok := t.language[tag.String()]; ok {
			accepted = append(accepted, acceptLanguage{tag.String(), quals[i]})
			continue
		}

		// Inspect the base language of the language tag. We use a modifier to
		// weight down the client-indicated qualifier, because there is no
		// exact match.
		base, confidence := tag.Base()
		var modifier = float32(0.6)
		if confidence <= golanguage.Low {
			modifier = float32(0.3)
		}
		if _, ok := t.language[base.String()]; ok {
			accepted = append(accepted, acceptLanguage{base.String(), quals[i] * modifier})
		}
	}

	if accepted.Len() == 0 {
		return nil
	}

	sort.Sort(accepted)
	var languages = make([]string, accepted.Len())
	for i, accept := range accepted {
		languages[i] = accept.t
	}
	return languages
}

func (t *Translations) Keys(lang string) map[string]string {
	if mapped, ok := t.language[lang]; ok {
		return mapped
	}
	return nil
}

func (t *Translations) Load(name string) error {
	todo := []string{name}
	for len(todo) > 0 {
		name, todo = todo[0], todo[1:]
		next, err := t.parse(name)
		if err != nil {
			return err
		}
		todo = append(todo, next...)
	}
	return nil
}

func (t *Translations) Supported() map[string]bool {
	s := make(map[string]bool)
	for lang := range t.language {
		s[lang] = true
	}
	return s
}

func (t *Translations) parse(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c Config
	c.Language = make(map[string]dictionary)
	if err = yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("%s: %v\n", name, err)
	}
	for lang, dict := range c.Language {
		t.language[lang] = dict
	}
	if (c.Include == nil || len(c.Include) == 0) && len(c.Language) == 0 {
		return nil, fmt.Errorf("%s: does not export any languages or includes", name)
	}
	return c.Include, nil
}

func (t *Translations) Translate(lang, in string) (string, error) {
	if _, ok := t.language[lang]; ok {
		if out, ok := t.language[lang][in]; ok {
			return out, nil
		}
	}
	return in, nil
}

func (t *Translations) Pongo2Filter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	out, err := t.Translate(param.String(), in.String())
	if err != nil {
		return in, nil
	}
	return pongo2.AsValue(out), nil
}

type tagTransNode struct {
	wrapper      *pongo2.NodeWrapper
	translations *Translations
	language     string
	arg          pongo2.IEvaluator
}

func (node *tagTransNode) Execute(ctx *pongo2.ExecutionContext, w pongo2.TemplateWriter) *pongo2.Error {
	b := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB
	if err := node.wrapper.Execute(ctx, b); err != nil {
		return err
	}

	var (
		out string
		err error
	)

	if node.arg != nil {
		val, perr := node.arg.Evaluate(ctx)
		if perr != nil {
			return perr
		}
		out, err = node.translations.Translate(node.language, val.String())
	} else if node.language != "" {
		out, err = node.translations.Translate(node.language, b.String())
	}
	if err != nil {
		return ctx.Error(err.Error(), nil)
	}
	w.WriteString(out)
	return nil
}

func (t *Translations) Pongo2TagParser(doc *pongo2.Parser, start *pongo2.Token, args *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	transNode := &tagTransNode{}

	langToken := args.MatchType(pongo2.TokenString)
	if langToken == nil {
		node, err := args.ParseExpression()
		if err != nil {
			return nil, err
		}
		transNode.arg = node
	} else {
		transNode.language = langToken.Val
	}

	wrapper, _, err := doc.WrapUntilTag("endtrans")
	if err != nil {
		return nil, err
	}
	transNode.wrapper = wrapper
	transNode.translations = t

	if args.Remaining() > 0 {
		return nil, args.Error("Malformed trans-tag arguments.", nil)
	}

	return transNode, nil
}
