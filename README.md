# cliapp 

simple cliapp for golang

- auto generate help info
- support add multi commands
- support command alias
- support cli color, simple to use. like `<info>info style</>`

## install

- use dep

```bash
dep ensure -add github.com/golangkit/cliapp
```

- go get

```bash
go get -u github.com/golangkit/cliapp
```

## quick start

```go 
package main

import (
    "runtime"
    "github.com/golangkit/cliapp"
    "github.com/golangkit/cliapp/demo/cmd"
)

// for test run: go build ./demo/cliapp.go && ./cliapp
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    app := cliapp.NewApp()
    app.Version = "1.0.3"
    app.Verbose = cliapp.VerbDebug
    app.Description = "this is my cli application"

    app.Add(cmd.ExampleCommand())
    app.Add(cmd.GitCommand())
    app.Add(&cliapp.Command{
        Name: "demo",
        Aliases: []string{"dm"},
        Description: "this is a description message for demo",
        Execute: func (cmd *cliapp.Command, args []string) int {
            cliapp.Stdout("hello, in the demo command\n")
            return 0
        },
    })

    // .... add more ...

    app.Run()
}
```

## Usage

### build package 

```bash
% go build ./demo/cliapp.go                                                           
```

### display app version

```bash
% ./cliapp --version
this is my cli application

Version: 1.0.3                                                           
```

### display app help

```bash
% ./cliapp                                                            
this is my cli application
Usage:
  ./cliapp command [--option ...] [argument ...]

Options:
  -h, --help        Display this help information
  -V, --version     Display this version information

Commands:
  demo         this is a description message for demo(alias: dm)
  example      this is a description message(alias: exp,ex)
  git          collect project info by git info(alias: git-info)
  help         display help information

Use "./cliapp help [command]" for more information about a command

```

### run a command

```bash
% ./cliapp example --id 12 -c val ag0 ag1                          
hello, in example command
opts {id:12 c:val dir:}
args is [ag0 ag1]

```

### display command help

```bash
% ./cliapp example -h                                                
this is a description message

Name: example(alias: exp,ex)
Usage: ./cliapp example [--option ...] [argument ...]

Global Options:
  -h, --help        Display this help information

Options:
  -c string
        the short option (default value)
  --dir string
        the dir option
  --id int
        the id option (default 2)

Arguments:
  arg0        the first argument
  arg1        the second argument
 
Examples:
  ./cliapp example --id 12 -c val ag0 ag1

```

## License

MIT
