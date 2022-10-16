# Gflag

`gflag` provide command line options and arguments binding, parse, management.

## GoDoc

Please see https://pkg.go.dev/github.com/gookit/gcli/v2/gflag

## Install

```shell
go get github.com/gookit/gcli/v2/gflag
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/gookit/gcli/v2/gflag"
)

func main() {
	gf := gflag.New("testFlags")
	gf.SetHandle(func(p *gflag.Parser) error {
		fmt.Println(p.Name())
		return nil
	})

	gf.Parse()
}
```
