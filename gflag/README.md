# Gflag

`gflag` provide command line options and arguments binding, parse, management.

## GoDoc

Please see https://pkg.go.dev/github.com/gookit/gcli/v3/gflag

## Install

```shell
go get github.com/gookit/gcli/v3/gflag
```

## Usage

```go file="demo.go"
package main

import (
	"fmt"
	"os"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil"
)

var name string

func main() {
	gf := gflag.New("testFlags")
	gf.StrOpt(&name, "name", "n", "", "")

	gf.SetHandle(func(p *gflag.Parser) error {
		fmt.Println(p.Name())
		return nil
	})

	goutil.MustOK(gf.Parse(os.Args[1:]))
}
```

### Run

```shell
go run demo.go
```
