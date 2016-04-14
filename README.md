# go-i18n - Internationalization for Go
[![GoDoc](https://godoc.org/git.maze.io/maze/go-i18n?status.svg)](https://godoc.org/git.maze.io/maze/go-i18n)

## Getting started

If you want to use go-i18n in your project, you first have to make one or more
translation maps. There are several formats supported.

### [JSON](http://json.org/)

All `.json` and `.js` files loaded with the `NewTemplateFile` loader are interpreted as a flat list of key-value pairs. Nested items are allowed and will be concatenated with a dot.

For example:

```JSON
{
  "foo": {
    "bar": {
      "testing"
    }
  }
}
```

Here, `foo → bar` is accessible as `foo.bar` translation key.

### [YAML](http://yaml.org/)

All `.yml` and `.yaml` files loaded with the `NewTemplateFile` loader are
interpreted as *Rails Application compatible* language files, see the [Rails
Internationalization][] page for more details.

[Rails Internationalization]: http://guides.rubyonrails.org/i18n.html

## Basic project structure

A basic project structure may look like follows:

```Go
package app

import (
    "git.maze.io/maze/go-i18n"

    "golang.org/net/http"
    "golang.org/x/text/language"
)

func main() {
    var lang = i18n.New(language.English)

    lang.Add(i18n.NewMap(language.English, map[string]string{
        "hello %s": "Hello %s!",
    }))

    lang.Add(i18n.NewMap(language.Dutch, map[string]string{
        "hello %s": "Hallo %s!",
    }))

    lang.Add(i18n.NewMap(language.SimplifiedChinese, map[string]string{
        "hello %s": "%s 你好！",
    }))
        
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        t := lang.Accept(r.Header().Get("Accept-Language"))
        t.Fprintf(w, "hello %s", r.RemoteAddr)
    })
    http.ListenAndServe(":8000", nil)
}
```

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
