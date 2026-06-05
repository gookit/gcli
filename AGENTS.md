# GCli 项目说明

## 项目概述

GCli 是一个用 Golang 编写的简单易用的命令行应用程序和工具库。它提供了丰富的功能，包括运行命令、颜色样式、数据展示、进度显示、交互方法等。

## 核心功能

- **多命令支持**: 支持添加多个命令及命令别名
- **多级命令**: 支持多级命令结构，每级命令可绑定独立的选项
- **选项/参数绑定**: 支持从结构体绑定命令选项，支持 `flag` 风格的选项和 `argument` 参数
- **颜色输出**: 支持丰富的颜色输出，基于 [gookit/color](https://github.com/gookit/color)
- **交互功能**: 内置用户交互方法：`ReadLine`, `Confirm`, `Select`, `MultiSelect` 等
- **进度显示**: 内置进度显示方法：`Txt`, `Bar`, `Loading`, `RoundTrip`, `DynamicText` 等
- **自动补全**: 支持生成 `zsh` 和 `bash` 命令补全脚本
- **帮助信息**: 自动生成命令帮助信息，支持颜色显示
- **错误提示**: 命令输入错误时会提示相似命令（包含别名提示）

## 技术架构

项目采用模块化设计，主要包括以下核心模块：

- `app.go`: 定义 CLI 应用程序结构和主逻辑
- `cmd.go`: 定义命令结构和命令执行逻辑
- `gflag/flags.go`: 实现命令行参数解析系统
- `interact/interact.go`: 提供交互式用户输入功能
- `progress/progress.go`: 实现终端进度条显示
- `show/show.go`: 提供数据格式化显示工具

## 构建与运行

### 构建

```bash
# 安装依赖
go get github.com/gookit/gcli/v3

# 编译示例应用
go build ./_examples/cliapp.go
```

### 运行

```bash
# 显示版本信息
./cliapp --version
# 或
./cliapp -V

# 显示应用帮助
./cliapp
# 或
./cliapp -h

# 运行具体命令
./cliapp example -c some.txt -d ./dir --id 34 -n tom -n john val0 val1 val2
```

## 开发约定

### 代码结构

- 主要功能模块按目录组织（app, cmd, gflag, interact, progress, show 等）
- 每个模块内有明确的接口定义和实现
- 遵循 Go 语言编码规范

### 依赖管理

项目使用 Go Modules 进行依赖管理，主要依赖包括：
- `github.com/gookit/color`: 提供颜色输出功能
- `github.com/gookit/goutil`: 提供通用工具函数
- `golang.org/x/crypto`: 提供加密相关功能
- `github.com/gookit/goutil/x/assert` - 单元测试断言功能

### 测试

项目包含完整的测试用例，可通过以下命令运行：
```bash
go test -v ./...
```

> `github.com/gookit/goutil/x/assert` 提供单元测试断言功能

## 使用示例

### 快速开始

```go
package main

import (
    "github.com/gookit/gcli/v3"
    "github.com/gookit/gcli/v3/_examples/cmd"
)

func main() {
    app := gcli.NewApp()
    app.Version = "1.0.3"
    app.Desc = "this is my cli application"

    app.Add(cmd.Example)
    app.Add(&gcli.Command{
        Name: "demo",
        Desc: "this is a description <info>message</> for {$cmd}", 
        Subs: []*gcli.Command {
            // ... allow add subcommands
        },
        Aliases: []string{"dm"},
        Func: func (cmd *gcli.Command, args []string) error {
            gcli.Print("hello, in the demo command\n")
            return nil
        },
    })

    app.Run(nil)
}
```

### 生成自动补全脚本

```go
import  "github.com/gookit/gcli/v3/builtin"

// 在应用中添加生成命令
app.Add(builtin.GenAutoComplete())
```

运行命令生成补全脚本：
```bash
./cliapp genac
```

## 项目特点

- 丰富的功能且易于使用
- POSIX 风格的短选项合并（`-a -b` = `-ab`）
- 支持设置 `Required` 选项参数
- 支持自定义参数验证器
- 自动检测并收集命令运行时的参数
- 支持独立运行单个命令作为独立应用
