# go-i18n - Internationalisation for Go
[![GoDoc](https://godoc.org/git.maze.io/maze/go-i18n?status.svg)](https://godoc.org/git.maze.io/maze/go-i18n)

## Integration

You can use i18n with many web frameworks, some of which are documented below.

### [Revel](http://revel.github.io/) integration

```Go
package app

import "github.com/revel/revel"
import "git.maze.io/maze/go-i18n"

var lang *i18n.I18N

func init() {
    var err error
    if lang, err = i18n.New(); err != nil {
        panic(err)
    }

    lang.TemplateInject(revel.TemplateFuncs)
}
```


### [Beego](http://beego.me/) integration

Add these three lines before "beego.Run()" in your main() function.

```Go
for k, fn := range lang.TemplateFuncs() {
    beego.AddFuncMap(k, fn)
}
```
