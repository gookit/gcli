# cliapp

golang下的命令行应用，工具库

**[EN Readme](README.md)**

## 功能特色

- 使用简单方便
- 可以添加多个命令，并且支持命令别名
- 支持单个命令当做独立应用运行
- 支持选项绑定，支持添加短选项
- 支持丰富的颜色输出。同时支持html标签式的颜色渲染
- 自动生成命令帮助信息，并且支持颜色显示

## 获取安装

- 使用 dep 包管理

```bash
dep ensure -add github.com/golangkit/cliapp
```

- 使用 go get

```bash
go get -u github.com/golangkit/cliapp
```

- git 克隆

```bash
git clone https://github.com/golangkit/cliapp
```

## 快速开始

如下，引入当前包就可以快速的编写cli应用了

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
        // allow color tag and {$cmd} will be replace to 'demo'
        Description: "this is a description <info>message</> for command", 
        Fn: func (cmd *cliapp.Command, args []string) int {
            cliapp.Stdout("hello, in the demo command\n")
            return 0
        },
    })

    // .... add more ...

    app.Run()
}
```

## Godoc

[godoc](https://godoc.org/github.com/golangkit/cliapp)

## 使用说明

先使用本项目下的 [demo](demo/) 示例代码构建一个小的cli demo应用

```bash
% go build ./demo/cliapp.go                                                           
```

### 打印版本信息

打印我们在创建cli应用时设置的版本信息

```bash
% ./cliapp --version
This is my cli application

Version: 1.0.3                                                           
```

### 应用帮助信息

使用 `./cliapp` 或者 `./cliapp -h` 来显示应用的帮助信息，包含所有的可用命令和一些全局选项

```bash
% ./cliapp                                                            
This is my cli application
Usage:
  ./cliapp {command} [--option ...] [argument ...]

Options:
  -h, --help        Display this help information
  -V, --version     Display this version information

Commands:
  demo         this is a description message for command(alias: dm)
  help         display help information

Use "./cliapp help [command]" for more information about a command

```

### 运行一个命令

```bash
% ./cliapp example --id 12 -c val ag0 ag1                          
hello, in example command
opts {id:12 c:val dir:}
args is [ag0 ag1]

```

### 显示一个命令的帮助

```bash
% ./cliapp example -h                                                
This is a description message

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

## 编写命令

### 简单使用

```go
app.Add(&cliapp.Command{
    Name: "demo",
    Aliases: []string{"dm"},
    // allow color tag and {$cmd} will be replace to 'demo'
    Description: "this is a description <info>message</> for command", 
    Fn: func (cmd *cliapp.Command, args []string) int {
        cliapp.Stdout("hello, in the demo command\n")
        return 0
    },
})
```

### 使用独立的文件

> the source file at: [example.go](demo/cmd/example.go)

```go
package cmd

import (
	cli "github.com/golangkit/cliapp"
	"fmt"
)

// The string flag list, implemented flag.Value interface
type Names []string

func (ns *Names) String() string {
	return fmt.Sprint(*ns)
}

func (ns *Names) Set(value string) error {
	*ns = append(*ns, value)
	return nil
}

// options for the command
var exampleOpts = ExampleOpts{}
type ExampleOpts struct {
	id  int
	c   string
	dir string
	opt string
	names Names
}

// ExampleCommand command definition
func ExampleCommand() *cli.Command {
	cmd := cli.Command{
		Fn:      exampleExecute,
		Name:    "example",
		Aliases: []string{"exp", "ex"},
		ArgList: map[string]string{
			"arg0": "the first argument",
			"arg1": "the second argument",
		},
		Description: "this is a description message",
		// {$binName} {$cmd} is help vars. '{$cmd}' will replace to 'example'
		Examples: `{$binName} {$cmd} --id 12 -c val ag0 ag1
  <cyan>{$fullCmd} --names tom --names john -n c</> test use special option`,
	}

	// use flag package func
	cmd.Flags.IntVar(&exampleOpts.id, "id", 2, "the id option")
	cmd.Flags.StringVar(&exampleOpts.c, "c", "value", "the short option")

	// use Command provided func
	cmd.StrOpt(&exampleOpts.dir, "dir", "d", "","the dir option")

	// setting option name and short-option name
	cmd.StrOpt(&exampleOpts.opt, "opt", "o", "", "the option message")

	// setting a special option var, it must implement the flag.Value interface
	cmd.VarOpt(&exampleOpts.names, "names", "n", "the option message")

	return &cmd
}

// command running
// example run:
// 	go build cliapp.go && ./cliapp example --id 12 -c val ag0 ag1
func exampleExecute(cmd *cli.Command, args []string) int {
	fmt.Print("hello, in example command\n")

	// fmt.Printf("%+v\n", cmd.Flags)
	fmt.Printf("opts %+v\n", exampleOpts)
	fmt.Printf("args is %v\n", args)

	return 0
}
```

## 使用颜色输出

## 如何使用

```go
package main

import (
    "github.com/golangkit/cliapp/color"
 )

func main() {
	// simple usage
	color.FgCyan.Printf("Simple to use %s\n", "color")

	// custom color
	color.New(color.FgWhite, color.BgBlack).Println("custom color style")

	// use style tag
	color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>\n")

	// set a style tag
	color.Tag("info").Println("info style text")

	// use info style tips
	color.Tips("info").Print("tips style text")

	// use info style blocked tips
	color.BlockTips("info").Print("blocked tips style text")
}
```

- 输出如下:

<img src="demo/colored-out.jpg" style="max-width: 320px;"/>

### 构建风格

```go
// 仅设置前景色
color.FgCyan.Printf("Simple to use %s\n", "color")
// 仅设置背景色
color.BgRed.Printf("Simple to use %s\n", "color")

// 完全自定义 前景色 背景色 选项
style := color.New(color.FgWhite, color.BgBlack, color.OpBold)
style.Println("custom color style")
```

```go
// 设置console颜色
color.Set(color.FgCyan)

// 输出信息
fmt.Print("message")

// 重置console颜色
color.Reset()
```

### 使用内置风格

- 使用颜色标签

使用颜色标签可以非常方便简单的构建自己需要的任何格式

```go
// 使用内置的 color tag
color.Print("<suc>he</><comment>llo</>, <cyan>wel</><red>come</>")
color.Println("<suc>hello</>")
color.Println("<error>hello</>")
color.Println("<warning>hello</>")

// 自定义颜色属性
color.Print("<fg=yellow;bg=black;op=underscore;>hello, welcome</>\n")
```

- 使用 `color.Tag`

给后面输出的文本信息加上给定的颜色风格标签

```go
// set a style tag
color.Tag("info").Print("info style text")
color.Tag("info").Printf("%s style text", "info")
color.Tag("info").Println("info style text")
```

### 内置的标签

这里列出了内置的标签，基本上涵盖了各种风格和颜色搭配。它们都可用着颜色标签，或者 `color.Tag` `color.Tips` 等的参数

```go
// Some internal defined style tags
// format is: "fg;bg;opt"
// usage: <tag>content text</>
var TagColors = map[string]string{
	// basic tags
	"red":     "0;31",
	"blue":    "0;34",
	"cyan":    "0;36",
	"black":   "0;30",
	"green":   "0;32",
	"brown":   "0;33",
	"white":   "1;37",
	"default": "39", // no color
	"normal":  "39", // no color
	"yellow":  "1;33",
	"magenta": "1;35",

	// alert tags, like bootstrap's alert
	"suc":     "1;32", // same "green" and "bold"
	"success": "1;32",
	"info":    "0;32", // same "green",
	"comment": "0;33", // same "brown"
	"note":    "36;1",
	"notice":  "36;4",
	"warn":    "0;30;43",
	"warning": "0;30;43",
	"primary": "0;34",
	"danger":  "0;31", // same "red"
	"err":     "30;41",
	"error":   "30;41",

	// more tags
	"lightRed":      "1;31",
	"light_red":     "1;31",
	"lightGreen":    "1;32",
	"light_green":   "1;32",
	"lightBlue":     "1;34",
	"light_blue":    "1;34",
	"lightCyan":     "1;36",
	"light_cyan":    "1;36",
	"lightDray":     "37",
	"light_gray":    "37",
	"gray":          "90",
	"darkGray":      "90",
	"dark_gray":     "90",
	"lightYellow":   "93",
	"light_yellow":  "93",
	"lightMagenta":  "95",
	"light_magenta": "95",

	// extra
	"lightRedEx":     "91",
	"light_red_ex":   "91",
	"lightGreenEx":   "92",
	"light_green_ex": "92",
	"lightBlueEx":    "94",
	"light_blue_ex":  "94",
	"lightCyanEx":    "96",
	"light_cyan_ex":  "96",
	"whiteEx":        "97",
	"white_ex":       "97",

	// option
	"bold":       "1",
	"underscore": "4",
	"reverse":    "7",
}
```

## 参考项目

- `issue9/term` https://github.com/issue9/term
- `beego/bee` https://github.com/beego/bee
- `inhere/console` https://github/inhere/php-console

## License

MIT
