# cliapp 

simple cliapp for golang

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

## demos

```bash
% go build ./demos/cliapp.go && ./cliapp                                                            
this is my cli application
Usage:
  ./cliapp command [--options ...] [arguments ...]

Options:
  -h, --help        Display this help information
  -V, --version     Display this version information

Commands:
  example      this is a description message(alias: exp,ex)
  git          collect project info by git info(alias: git-info)
  help         display help information

Use "./cliapp help [command]" for more information about a command
```
