package i18n

import (
	"testing"

	"github.com/flosch/pongo2"
)

var languages = []string{"en", "nl"}

func Test(t *testing.T) {
	tran, err := New("testdata/i18n.yaml")
	if err != nil {
		t.Fatal(err)
	}
	for _, lang := range languages {
		out, err := tran.Translate(lang, "hello world")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s: %q translates to %q", lang, "hello world", out)
	}
}

func TestPongo2Filter(t *testing.T) {
	tran, err := New("testdata/i18n.yaml")
	if err != nil {
		t.Fatal(err)
	}

	pongo2.RegisterFilter("trans", tran.Pongo2Filter)

	tpl, err := pongo2.FromString(`{{ "hello world"|trans:lang }}`)
	if err != nil {
		t.Fatal(err)
	}
	for _, lang := range languages {
		out, err := tpl.Execute(pongo2.Context{"lang": lang})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s: %q translates to %q", lang, "hello world", out)
	}
}

func TestPongo2Tag(t *testing.T) {
	tran, err := New("testdata/i18n.yaml")
	if err != nil {
		t.Fatal(err)
	}

	pongo2.RegisterTag("trans", tran.Pongo2TagParser)

	tpl, err := pongo2.FromString(`{% trans lang %}hello world{% endtrans %}`)
	if err != nil {
		t.Fatal(err)
	}
	for _, lang := range languages {
		out, err := tpl.Execute(pongo2.Context{"lang": lang})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%s: %q translates to %q", lang, "hello world", out)
	}
}
