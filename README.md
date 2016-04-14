# go-i18n - Internationalisation for Go
[![GoDoc](https://godoc.org/git.maze.io/maze/go-i18n?status.svg)](https://godoc.org/git.maze.io/maze/go-i18n)

## Integration

You can use i18n with many web frameworks, some of which are documented below.

### [Beego](http://beego.me/) integration

Add these three lines before "beego.Run()" in your main() function.

```Go
for k, fn := range lang.TemplateFuncs() {
    beego.AddFuncMap(k, fn)
}
```

### [Revel](http://revel.github.io/) integration

```Go
package app

import (
    "git.maze.io/maze/go-i18n"
    "github.com/revel/revel"
    "golang.org/x/text/language"
)

var lang *i18n.I18N

func init() {
    var err error
    if lang, err = i18n.New(language.English); err != nil {
        panic(err)
    }

    lang.TemplateInject(revel.TemplateFuncs)
}
```

### [Pongo2](https://github.com/flosch/pongo2) integration

```Go
package app

import (
    "git.maze.io/maze/go-i18n"
    "github.com/flosch/pongo2"
    "golang.org/x/text/language"
)

var lang *i18n.I18N

func init() {
    var err error
    if lang, err = i18n.New(language.English); err != nil {
        panic(err)
    }

    pongo2.RegisterFilter("translate",
        func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
            if param != nil {
                return pongo2.AsValue(lang.Translate(param.String(), in.String())), nil
            }
            return pongo2.AsValue(lang.Translate(lang.Default().String(), in.String())), nil
        },
    )
}
```
