# GCli

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/gookit/gcli?style=flat-square)
[![Actions Status](https://github.com/gookit/gcli/workflows/action-tests/badge.svg)](https://github.com/gookit/gcli/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/gookit/gcli)](https://github.com/gookit/gcli)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/4f071e6858fb4117b6c1376c9316d8ef)](https://www.codacy.com/gh/gookit/gcli/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=gookit/gcli&amp;utm_campaign=Badge_Grade)
[![Go Reference](https://pkg.go.dev/badge/github.com/gookit/goutil.svg)](https://pkg.go.dev/github.com/gookit/goutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/gookit/gcli)](https://goreportcard.com/report/github.com/gookit/gcli)
[![Coverage Status](https://coveralls.io/repos/github/gookit/gcli/badge.svg?branch=master)](https://coveralls.io/github/gookit/gcli?branch=master)

Golang 编写的简单易用的命令行应用，工具库。包含运行命令，颜色风格，数据展示，进度显示，交互方法等

> [EN RAEDME](README.md)

![app-cmd-list](_examples/images/cmd-list.png)

## 功能特色

使用简单方便，功能丰富。按能力分组的特性亮点：

**命令**

- 多级（嵌套）命令，每级命令均可绑定自己的选项
- 命令 **别名**；输入错误时提示相似命令（包含别名提示）
- 命令/应用中间件 `Use(handlers ...RunnerFunc)`
- 支持将单个命令当做独立应用运行

**选项绑定**

- 代码式绑定：`BoolOpt / IntOpt / StrOpt / Float64Opt / VarOpt ...`
- 泛型、类型安全绑定：`gflag.Opt[T]` / `gflag.BindVar[T]`——一次调用即覆盖
  `bool/int/uint/float/string`、`time.Duration`、`[]string/[]int/[]bool`、`map[string]string`
- 结构体标签绑定 `FromStruct`：
    - 三种标签规则：`named`(默认) / `simple` / `field`(用字段名做选项名 + 自动展开匿名嵌套结构体)
    - 字段类型：`bool/int/uint/float/string`、原生 `[]string/[]int/[]bool`(可重复)、`time.Duration`、`map[string]string`(可重复 `--meta k=v`)
    - `enum:"a,b,c"` 标签：设置取值候选(补全)并做成员校验
- 每个选项支持 `Required` / `Validator` / `Choices`；选项可设 `Category` 在帮助中分组显示

**三层选项模型**

- 全局（App 级）选项
- 共享(继承)选项 `Command.SharedOpts()`（对标 cobra `PersistentFlags`）：被该命令及其所有子孙命令继承
  （共享同一变量），在帮助信息中归入 `Inherited Options` 分组
- 局部（命令级）选项

**解析增强**

- 写在参数 **之后** 的选项会被自动重排：`cmd arg --name tom` 等同于 `cmd --name tom arg`
  （默认开启，可用 `gflag.WithReorderArgs(false)` 关闭）
- POSIX 短选项合并（`-ab` = `-a -b`），通过 opt-in 的 `EnhanceShort` 开启；`EnhanceShort=2` 额外支持
  取值紧贴写法 `-Ostdout` = `-O stdout`。仅全为 bool 的组合才拆分，默认关闭以保持兼容
- 选项值为空时支持通过 `Question` 声明式交互收集输入

**参数**

- 绑定命名参数，支持 `必须` / 可选 / `数组` 设置
- 运行命令时自动检测并按对应关系收集参数

**工具**

- 生成 `zsh` / `bash` / `pwsh` 命令补全脚本（含动态补全）
- 生成 `markdown` / `man page` 命令文档（`docgen` 包 + builtin `GenDoc` 命令）
- 自动生成、带颜色渲染的命令帮助信息
- 事件钩子系统（`gevent`，提供 `gcli.Evt*` 别名）

**周边**

- 颜色输出、交互输入、进度与数据展示，分别由 [gookit/color](https://github.com/gookit/color) 与 [gookit/cliui](https://github.com/gookit/cliui) 提供

## GoDoc

- [godoc](https://pkg.go.dev/github.com/gookit/gcli/v3)

## 安装

```bash
go get github.com/gookit/gcli/v3
```

> **升级提示：** 事件包 `gcli/v3/events` 已重命名为 `gcli/v3/gevent`。请更新导入路径；
> 或直接使用 `gcli.Evt*` 别名（如 `gcli.EvtCmdRunBefore`）引用事件名，即可无需导入该包。
> 详见 [CHANGELOG](CHANGELOG.md)。

## 快速开始

如下，引入当前包就可以快速的编写cli应用了

```go
package main

import (
    "runtime"
    "github.com/gookit/gcli/v3"
    "github.com/gookit/gcli/v3/_examples/cmd"
)

// 测试运行: go run ./_examples/cliapp.go && ./cliapp
func main() {
    app := gcli.NewApp()
    app.Version = "1.0.3"
    app.Desc = "this is my cli application"
    // app.SetVerbose(gcli.VerbDebug)

    // TIP: 添加绑定app级别的选项设置(与内置的 -h/--help 同等级)
    app.Flags().BoolOpt(...)
    app.Flags().StrOpt(...)

    app.Add(cmd.Example)
    app.Add(&gcli.Command{
        Name: "demo",
        // allow color tag and {$cmd} will be replace to 'demo'
        Desc: "this is a description <info>message</> for command",
        Aliases: []string{"dm"},
        Func: func (cmd *gcli.Command, args []string) error {
            gcli.Println("hello, in the demo command")
            return nil
        },
    })

    // .... add more ...

    app.Run(nil)
}
```

### 使用说明

先使用本项目下的 [demo](_examples/) 示例代码构建一个小的cli demo应用

```bash
% go build ./_examples/cliapp.go
```

#### 打印版本信息

打印我们在创建cli应用时设置的版本信息。如果你还设置了字符LOGO，也会显示出来。

```bash
% ./cliapp --version
# or use -V
% ./cliapp -V
```

![app-version](_examples/images/app-version.jpg)

#### 显示应用帮助信息

使用 `./cliapp` 或者 `./cliapp -h` 来显示应用的帮助信息，包含所有的可用命令和一些全局选项

示例：

```bash
./cliapp
./cliapp -h # can also
./cliapp --help # can also
```

![cmd-list](_examples/images/cmd-list.png)

#### 显示一个命令的帮助

显示一个指定命令的帮助信息

示例：

```bash
./cliapp {command} -h
./cliapp {command} --help
./cliapp help {command}
```

![cmd-help](_examples/images/cmd-help.png)

#### 相似命令提示

输入了错误的命令，但是有名称相似的会提示出来。

![cmd-tips](_examples/images/err-cmd-tips.jpg)

#### 运行一个命令

语法结构：

```text
./cliapp COMMAND [--OPTION VALUE -S VALUE ...] [ARGUMENT0 ARGUMENT1 ...]
```

示例

```bash
./cliapp ex -c some.txt -d ./dir --id 34 -n tom -n john val0 val1 val2 arrVal0 arrVal1 arrVal2
```

可以观察到选项和参数的搜集结果:

![run-example](_examples/images/run-example.png)

## 生成命令补全脚本

```go
import  "github.com/gookit/gcli/v3/builtin"

    // ...
    // 添加内置提供的生成命令
    app.Add(builtin.GenAutoComplete())

```

构建并运行生成命令(_生成成功后可以去掉此命令_)：

```bash
% go build ./_examples/cliapp.go && ./cliapp genac -h // 使用帮助
% go build ./_examples/cliapp.go && ./cliapp genac // 开始生成, 你将会看到类似的信息
INFO:
  {shell:zsh binName:cliapp output:auto-completion.zsh}

Now, will write content to file auto-completion.zsh
Continue? [yes|no](default yes): y

OK, auto-complete file generate successful
```

> 运行后就会在当前目录下生成一个 `auto-completion.{zsh|bash}` 文件， shell 环境名是自动获取的。当然你可以在运行时手动指定

生成的shell script 文件请参看：

- bash 环境 [auto-completion.bash](resource/auto-completion.bash)
- zsh 环境 [auto-completion.zsh](resource/auto-completion.zsh)

预览效果:

![auto-complete-tips](_examples/images/auto-complete-tips.jpg)

## 共享(继承)选项

`Command.SharedOpts()`（对标 cobra `PersistentFlags`）绑定的选项会被该命令 **及其所有子孙命令** 继承，
共享同一变量。可写在子命令段的任意位置，并在帮助信息中归入 `Inherited Options` 分组。

```go
var gitDir string

top := &gcli.Command{Name: "git", Desc: "git tools"}
// 在 SharedOpts 上绑定: 被每个子命令继承
top.SharedOpts().StrOpt(&gitDir, "git-dir", "", ".git", "the git dir path")

top.Add(&gcli.Command{
    Name: "status",
    Func: func(c *gcli.Command, _ []string) error {
        // --git-dir 虽在父命令上声明, 这里也可使用
        gcli.Printf("git dir: %s\n", gitDir)
        return nil
    },
})

// 使用: ./app git status --git-dir /path/to/.git
```

## 生成命令文档

添加内置的 `GenDoc` 命令后，即可为所有命令导出 `markdown` / `man` 文档：

```go
import "github.com/gookit/gcli/v3/builtin"

app.Add(builtin.GenDoc())
// ./cliapp gendoc -f md  -o ./docs   # 导出 markdown(默认)
// ./cliapp gendoc -f man -o ./docs   # 导出 man 文档
```

也可以编程方式调用：

```go
import "github.com/gookit/gcli/v3/docgen"

docgen.MarkdownTree(app, "./docs") // 每个命令一个 .md + index.md
docgen.ManTree(app, "./docs")      // man 文档
```

## 编写命令

### 简单使用

```go
app.Add(&gcli.Command{
    Name: "demo",
    // allow color tag and {$cmd} will be replace to 'demo'
    Desc: "this is a description <info>message</> for command",
    Aliases: []string{"dm"},
    Func: func (cmd *gcli.Command, args []string) error {
        gcli.Print("hello, in the demo command\n")
        return nil
    },
})
```

### 使用独立的文件

> the source file at: [example.go](_examples/cmd/example.go)

```go
package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
)

// options for the command
var exampleOpts = struct {
	id  int
	c   string
	dir string
	opt string
	names gcli.Strings
}{}

// ExampleCommand command definition
var ExampleCommand = &gcli.Command{
	Name:        "example",
	Desc: "this is a description message",
	Aliases:     []string{"exp", "ex"}, // 命令别名
	// {$binName} {$cmd} is help vars. '{$cmd}' will replace to 'example'
	Examples: `{$binName} {$cmd} --id 12 -c val ag0 ag1
<cyan>{$fullCmd} --names tom --names john -n c</> test use special option`,
	Config: func(c *gcli.Command) {
		// 绑定命令选项信息
		c.IntOpt(&exampleOpts.id, "id", "", 2, "the id option")
		c.StrOpt(&exampleOpts.c, "config", "c", "value", "the config option")
		// notice `DIRECTORY` will replace to option value type
		c.StrOpt(&exampleOpts.dir, "dir", "d", "", "the `DIRECTORY` option")
		// 支持设置选项短名称
		c.StrOpt(&exampleOpts.opt, "opt", "o", "", "the option message")
		// 支持绑定自定义变量, 但必须实现 flag.Value 接口
		c.VarOpt(&exampleOpts.names, "names", "n", "the option message")

		// 绑定命令参数信息，按参数位置绑定
		c.AddArg("arg0", "the first argument, is required", true)
		c.AddArg("arg1", "the second argument, is required", true)
		c.AddArg("arg2", "the optional argument, is optional")
		c.AddArg("arrArg", "the array argument, is array", false, true)
	},
	Func:  exampleExecute,
}

// 命令执行主逻辑代码
// example run:
// 	go run ./_examples/cliapp.go ex -c some.txt -d ./dir --id 34 -n tom -n john val0 val1 val2 arrVal0 arrVal1 arrVal2
func exampleExecute(c *gcli.Command, args []string) error {
	fmt.Print("hello, in example command\n")

	magentaln := color.Magenta.Println

	magentaln("All options:")
	fmt.Printf("%+v\n", exampleOpts)
	magentaln("Raw args:")
	fmt.Printf("%v\n", args)

	magentaln("Get arg by name:")
	arr := c.Arg("arrArg")
	fmt.Printf("named array arg '%s', value: %v\n", arr.Name, arr.Value)

	magentaln("All named args:")
	for _, arg := range c.Args() {
		fmt.Printf("named arg '%s': %+v\n", arg.Name, *arg)
	}

	return nil
}
```

- 查看此命令的帮助信息：

```bash
go build ./_examples/cliapp.go && ./cliapp example -h
```

> 漂亮的帮助信息就已经自动生成并展示出来了

![cmd-help](_examples/images/cmd-help.png)

### 添加选项

添加选项可用的方法：

```go
BoolOpt(p *bool, name, shorts string, defValue bool, desc string)
BoolVar(p *bool, meta FlagMeta)
Float64Opt(p *float64, name, shorts string, defValue float64, desc string)
Float64Var(p *float64, meta FlagMeta)
Int64Opt(p *int64, name, shorts string, defValue int64, desc string)
Int64Var(p *int64, meta FlagMeta)
IntOpt(p *int, name, shorts string, defValue int, desc string)
IntVar(p *int, meta FlagMeta)
StrOpt(p *string, name, shorts, defValue, desc string)
StrVar(p *string, meta FlagMeta)
Uint64Opt(p *uint64, name, shorts string, defValue uint64, desc string)
Uint64Var(p *uint64, meta FlagMeta)
UintOpt(p *uint, name, shorts string, defValue uint, desc string)
UintVar(p *uint, meta FlagMeta)
Var(p flag.Value, meta FlagMeta)
VarOpt(p flag.Value, name, shorts, desc string)
```

Usage examples:

```go
var id int
var b bool
var opt, dir string
var f1 float64
var names gcli.Strings

// bind options
cmd.IntOpt(&id, "id", "", 2, "the id option")
cmd.BoolOpt(&b, "bl", "b", false, "the bool option")
// notice `DIRECTORY` will replace to option value type
cmd.StrOpt(&dir, "dir", "d", "", "the `DIRECTORY` option")
// setting option name and short-option name
cmd.StrOpt(&opt, "opt", "o", "", "the option message")
// setting a special option var, it must implement the flag.Value interface
cmd.VarOpt(&names, "names", "n", "the option message")
```

#### 使用泛型绑定

`gflag.Opt[T]` / `gflag.BindVar[T]` 用一次类型安全的调用绑定一个带类型的指针，
覆盖标量、`time.Duration`、切片以及 `map[string]string`：

```go
var name string
var tags []string

gflag.Opt(cmd.Flags(), &name, "name", "n", "tom", "the user name")
// 切片选项，可重复: --tag php --tag go
gflag.Opt(cmd.Flags(), &tags, "tag", "t", nil, "the tags, repeatable")
```

#### 使用 struct 定义选项

```go
package main

import (
	"github.com/gookit/gcli/v3"
)

type userOpts struct {
	Int  int    `flag:"name=int0;shorts=i;required=true;desc=int option message"`
	Bol  bool   `flag:"name=bol;shorts=b;desc=bool option message"`
	Str1 string `flag:"name=str1;shorts=o;required=true;desc=str1 message"`
	// use ptr
	Str2 *string `flag:"name=str2;required=true;desc=str2 message"`
	// custom type and implement flag.Value
	Verb0 gcli.VerbLevel `flag:"name=verb0;shorts=v0;desc=verb0 message"`
	// use ptr
	Verb1 *gcli.VerbLevel `flag:"name=verb1;desc=verb1 message"`
}

// run: go run ./_examples/issues/iss157.go
func main() {
	astr := "xyz"
	verb := gcli.VerbWarn

	cmd := gcli.NewCommand("test", "desc")
	cmd.Config = func(c *gcli.Command) {
		c.MustFromStruct(&userOpts{
			Str2:  &astr,
			Verb1: &verb,
		})
	}

	// disable auto bind global options: verbose,version, progress...
	gcli.GOpts().SetDisable()

	// direct run
	if err := cmd.Run(nil); err != nil {
		colorp.Errorln( err)
	}
}
```

#### 结构体标签规则

`FromStruct` 支持三种标签规则，通过 `c.FromStruct(ptr, ruleType)` 选择：

- `gcli.TagRuleNamed`（默认）：`flag:"name=int0;shorts=i;required=true;desc=message"`
- `gcli.TagRuleSimple`：`flag:"desc;required;default;shorts"`
- `gcli.TagRuleField`：用**字段名**(SnakeCase) 做选项名，元数据从独立 tag 键读取，并**自动展开匿名嵌套结构体**。

```go
type commonOpts struct {
	Verbose bool `flag:"v" desc:"enable verbose output"`
}
type demoOpts struct {
	commonOpts        // 匿名嵌套：展开为 --verbose/-v 选项
	UserName string `flag:"u" desc:"the user name" required:"true"`
	Age      int    `desc:"the user age" default:"18"`
}

c.MustFromStruct(&demoOpts{}, gcli.TagRuleField)
// => 选项: --user-name/-u (必填), --age (默认 18), --verbose/-v
```

#### 通过 Question 交互收集

当选项值为空时，可用一个交互式问题自动收集输入（内置默认收集器）。`Collector` 优先级高于 `Question`。

```go
c.StrOpt2(&token, "token", "the access token",
	gflag.WithQuestion("Please input your access token: "))
// 不带 --token 运行将提示: "Please input your access token: "
```

#### POSIX 短选项增强

组合短选项默认关闭，通过 `Config.EnhanceShort` 开启：

```go
c.ParserCfg().EnhanceShort = gcli.EnhanceShortMerge  // 1: -aux => -a -u -x (全为 bool 才拆)
c.ParserCfg().EnhanceShort = gcli.EnhanceShortAttach // 2: 额外支持 -Ostdout => -O stdout
```

也可以一次性**全局**设置，作用于所有命令——命令自身的设置（若有）仍然优先：

```go
gcli.SetEnhanceShort(gcli.EnhanceShortMerge) // 作用于每个命令
```

仅当组合中**全部**是 bool 短选项时才拆分；混合写法原样保留，避免误伤取值型短选项。

> 可运行示例见 `_examples/cmd`：`struct-flag`(B6)、`short-merge`(B4+B5)、`ask-demo`(B7)。

### 命令/选项分类

命令和选项都支持设置 `Category` 进行分组，用于在帮助信息中分类显示。
未设置分类时输出与原先一致（命令归入 `Available Commands`、选项直接列在 `Options:` 下）。

为命令设置分类：

```go
app.Add(&gcli.Command{Name: "migrate", Desc: "run db migrate", Category: "database"})
app.Add(&gcli.Command{Name: "seed", Desc: "seed db data", Category: "database"})
app.Add(&gcli.Command{Name: "serve", Desc: "start http server"}) // 默认分组
```

为选项设置分类（结构体字面量或 `WithCategory`）：

```go
cmd.StrVar(&dsn, &gcli.CliOpt{Name: "db-dsn", Desc: "database dsn", Category: "database"})
// 或使用配置函数
cmd.StrOpt2(&port, "port", "bind port", gflag.WithCategory("network"))
```

帮助信息中将按分类分组展示（分组顺序为首次出现的顺序，组内按名称排序）：

```text
Available Commands:
  serve       Start http server
Database:
  migrate     Run db migrate
  seed        Seed db data
```

> 提示：日志级别不再通过全局 `--verbose` 选项控制，改由环境变量 `GCLI_VERBOSE`（如 `GCLI_VERBOSE=debug`）或代码 `gcli.SetVerbose()` 设置，避免污染上层应用的选项列表。

### 绑定参数

关于参数定义：

- `必须的` 参数不能定义在 `可选参数` 之后
- 只允许有一个数组参数（多个值的）
- 数组参数只能定义在最后

绑定参数可用的方法:

```go
Add(arg Argument) *Argument
AddArg(name, desc string, requiredAndArrayed ...bool) *Argument
AddArgument(arg *Argument) *Argument
BindArg(arg Argument) *Argument
```

用法示例：

```go
cmd.AddArg("arg0", "the first argument, is required", true)
cmd.AddArg("arg1", "the second argument, is required", true)
cmd.AddArg("arg2", "the optional argument, is optional")
cmd.AddArg("arrArg", "the array argument, is array", false, true)
```

也可以使用 `Add()/BindArg()`:

```go
cmd.Add("arg0", gcli.Argument{
	Name: "ag0",
	Desc: "the first argument, is required",
	Require: true,
})
cmd.BindArg("arg2", gcli.Argument{
	Name: "ag0",
	Desc: "the third argument, is is optional",
})

cmd.BindArg("arrArg", gcli.Argument{
	Name: "arrArg",
	Desc: "the third argument, is is array",
	Arrayed: true,
})
```

### 获取参数

可以通过 `c.Arg(name string) *gcli.Argument` 获取参数，通过上面内置的方法可以将参数转换文常用的数据类型

```go
var MyCommand = &gcli.Command{
    Name: "example",
    Desc: "this is an example command",
    Config: func(c *gcli.Command) {
        cmd.BindArg("arg0", gcli.Argument{
            Name: "ag0",
            Desc: "the first argument, is required",
            Require: true,
        })
        cmd.Add("arg1", gcli.Argument{
            Name: "ag1",
            Desc: "the second argument, is is optional",
        })
    },
    Func: func(c *gcli.Command, args []string) error {
        arg0 := c.Arg("arg0").String()
        arg1 := c.Arg("arg1").Int()

        fmt.Println(arg0, arg1)
        return nil
    },
}
```

## 周边能力：颜色 / 交互 / 进度 / 展示

gcli 内置了颜色输出、交互输入（`Confirm` / `Select` / `ReadLine` 等）、
进度显示（`Bar` / `Spinner` / `Loading` 等）以及数据展示（表格 / 列表 / 树），
分别由 [gookit/color](https://github.com/gookit/color) 与 [gookit/cliui](https://github.com/gookit/cliui) 提供。

```go
color.Info.Tips("processing...")              // 颜色输出

ok := interact.Confirm("ensure continue?")    // 交互确认
if !ok {
    return nil
}

p := progress.Bar(100)                        // 进度条
p.Start();
/* 循环中 p.Advance() */
p.Finish()
```

> 更多用法详见 [gookit/color](https://github.com/gookit/color) 与 [gookit/cliui](https://github.com/gookit/cliui)。

## Gookit 工具包

- [gookit/ini](https://github.com/gookit/ini) INI配置读取管理，支持多文件加载，数据覆盖合并, 解析ENV变量, 解析变量引用
- [gookit/rux](https://github.com/gookit/rux) Simple and fast request router for golang HTTP
- [gookit/gcli](https://github.com/gookit/gcli) Go的命令行应用，工具库，运行CLI命令，支持命令行色彩，用户交互，进度显示，数据格式化显示
- [gookit/event](https://github.com/gookit/event) Go实现的轻量级的事件管理、调度程序库, 支持设置监听器的优先级, 支持对一组事件进行监听
- [gookit/config](https://github.com/gookit/config) Go应用配置管理，支持多种格式（JSON, YAML, TOML, INI, HCL, ENV, Flags），多文件加载，远程文件加载，数据合并
- [gookit/color](https://github.com/gookit/color) CLI 控制台颜色渲染工具库, 拥有简洁的使用API，支持16色，256色，RGB色彩渲染输出
- [gookit/filter](https://github.com/gookit/filter) 提供对Golang数据的过滤，净化，转换
- [gookit/validate](https://github.com/gookit/validate) Go通用的数据验证与过滤库，使用简单，内置大部分常用验证、过滤器
- [gookit/goutil](https://github.com/gookit/goutil) Go 的一些工具函数，格式化，特殊处理，常用信息获取等
- 更多请查看 https://github.com/gookit

## 参考项目

- `issue9/term` https://github.com/issue9/term
- `beego/bee` https://github.com/beego/bee
- `inhere/console` https://github/inhere/php-console
- [ANSI转义序列](https://zh.wikipedia.org/wiki/ANSI转义序列)
- [Standard ANSI color map](https://conemu.github.io/en/AnsiEscapeCodes.html#Standard_ANSI_color_map)
- go package: `golang.org/x/crypto/ssh/terminal`

## License

MIT
