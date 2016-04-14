package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"git.maze.io/maze/go-i18n"

	"gopkg.in/yaml.v2"

	"github.com/flosch/pongo2"
)

var (
	tags   = map[string]bool{}
	ignore = map[string]bool{
		".swp": true,
		".swo": true,
	}
)

func Pongo2Filter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	tags[in.String()] = true
	return pongo2.AsValue("<trans>"), nil
}

type tagTransNode struct {
	wrapper  *pongo2.NodeWrapper
	language string
	arg      pongo2.IEvaluator
}

func (node *tagTransNode) Execute(ctx *pongo2.ExecutionContext, w pongo2.TemplateWriter) *pongo2.Error {
	b := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB
	if err := node.wrapper.Execute(ctx, b); err != nil {
		return err
	}

	if node.arg != nil {
		val, perr := node.arg.Evaluate(ctx)
		if perr != nil {
			return perr
		}
		tags[val.String()] = true
	} else if node.language != "" {
		tags[b.String()] = true
	}
	w.WriteString("<trans>")
	return nil
}

func Pongo2TagParser(doc *pongo2.Parser, start *pongo2.Token, args *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
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

	if args.Remaining() > 0 {
		return nil, args.Error("Malformed trans-tag arguments.", nil)
	}

	return transNode, nil
}

func parse(r io.Reader) (err error) {
	var (
		b []byte
		t *pongo2.Template
	)

	if b, err = ioutil.ReadAll(r); err != nil {
		return
	}
	if t, err = pongo2.FromString(string(b)); err != nil {
		return
	}
	if _, err = t.Execute(pongo2.Context{"lang": "en"}); err != nil {
		return
	}
	return
}

func init() {
	pongo2.RegisterFilter("trans", Pongo2Filter)
	pongo2.RegisterTag("trans", Pongo2TagParser)
}

func main() {
	var err error

	fsBase := flag.String("fs-base", "", "File System Loader base path")
	languageFile := flag.String("languages", "languages.yaml", "Languages file")
	flag.Parse()

	log.SetOutput(os.Stderr)

	if *fsBase == "" {
		log.Fatalln("missing -fs-base")
	}
	if !filepath.IsAbs(*fsBase) {
		*fsBase, err = filepath.Abs(*fsBase)
		if err != nil {
			log.Fatalln(err)
		}
	}

	lang, err := i18n.New(*languageFile)
	if err != nil {
		log.Fatalln(err)
	}

	loader := pongo2.MustNewLocalFileSystemLoader(*fsBase)
	if err = os.Chdir(*fsBase); err != nil {
		log.Fatalln(err)
	}

	err = filepath.Walk(*fsBase, func(path string, info os.FileInfo, err1 error) (err error) {
		if info.IsDir() {
			return nil
		}
		if ignore[filepath.Ext(strings.ToLower(path))] {
			return nil
		}
		log.Println("parse", path)
		var r io.Reader
		if r, err = loader.Get(path); err == nil {
			err = parse(r)
		}
		return
	})
	if err != nil {
		log.Fatalln(err)
	}

	var langs = make(map[string]map[string]string)
	for l := range lang.Supported() {
		langs[l] = lang.Keys(l)
	}

	var tagList = []string{}
	for tag := range tags {
		tagList = append(tagList, tag)
	}
	sort.Sort(sort.StringSlice(tagList))

	for _, tag := range tagList {
		for l := range langs {
			if _, ok := langs[l][tag]; !ok {
				langs[l][tag] = ""
			}
		}
	}

	b, err := yaml.Marshal(struct {
		Lang map[string]map[string]string
	}{langs})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
}
