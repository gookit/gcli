# Gflag

`gflag` provide command line options and arguments parse, binding, management.

## GoDoc

Please see https://pkg.go.dev/github.com/gookit/gcli/v3

## Install

```shell
go get github.com/gookit/gcli/v3/gflag
```

## Flags

- `options`   - start with `-` or `--`, and the first character must be a letter
- `arguments` - not start with `-` or `--`, and after options.

### Flag Options

- Support long option. eg: `--long` `--long value`
- Support short option. eg: `-s -a value`
- Support define array option
    - eg: `--tag php --tag go` will get `tag: [php, go]`

### Flag Arguments

- Support binding named argument
- Support define array argument

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

## Binding methods

### Binding cli options

```go
Bool(name, shorts string, defVal bool, desc string) *bool
BoolOpt(ptr *bool, name, shorts string, defVal bool, desc string)
BoolVar(ptr *bool, opt *CliOpt)
Float64Opt(p *float64, name, shorts string, defVal float64, desc string)
Float64Var(ptr *float64, opt *CliOpt)

Int(name, shorts string, defValue int, desc string) *int
Int64(name, shorts string, defValue int64, desc string) *int64
Int64Opt(ptr *int64, name, shorts string, defValue int64, desc string)
Int64Var(ptr *int64, opt *CliOpt)
IntOpt(ptr *int, name, shorts string, defValue int, desc string)
IntVar(ptr *int, opt *CliOpt)

Str(name, shorts string, defValue, desc string) *string
StrOpt(p *string, name, shorts, defValue, desc string)
StrVar(p *string, opt *CliOpt)

Uint(name, shorts string, defVal uint, desc string) *uint
Uint64(name, shorts string, defVal uint64, desc string) *uint64
Uint64Opt(ptr *uint64, name, shorts string, defVal uint64, desc string)
Uint64Var(ptr *uint64, opt *CliOpt)
UintOpt(ptr *uint, name, shorts string, defValue uint, desc string)
UintVar(ptr *uint, opt *CliOpt)

Var(ptr flag.Value, opt *CliOpt)
VarOpt(v flag.Value, name, shorts, desc string)
```

### Binding cli arguments

```go
AddArg(name, desc string, requiredAndArrayed ...bool) *CliArg
AddArgByRule(name, rule string) *CliArg
AddArgument(arg *CliArg) *CliArg
BindArg(arg *CliArg) *CliArg
```
