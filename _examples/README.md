# example

## run example

### cli application

show help:

```bash
go run ./cliapp -h
```

run application:

```bash
go run ./cliapp demo
```

#### feature demo commands

`cliapp` 内置了几个演示最近新增特性的命令（源码在 `cmd/`）：

```bash
# B6: TagRuleField 标签规则 + 匿名字段展开
go run ./cliapp struct-flag --user-name tom --age 22 -v
go run ./cliapp struct-flag -u tom            # age 用默认值 18

# B4+B5: EnhanceShort POSIX 短选项合并
go run ./cliapp short-merge -aux              # = -a -u -x (全 bool 才拆)
go run ./cliapp short-merge -Ostdout          # = -O stdout (level2 紧贴取值)

# B7: Question 声明式交互收集
go run ./cliapp ask-demo                       # 不带 --token 会自动提问收集
go run ./cliapp ask-demo --token abc123        # 已提供值则不提问
```

### Only one command

show help:

```bash
go run ./simpleone -h
```

run command:

```bash
go run ./simpleone
```

### Multi level commands

show help:

```bash
go run ./multilevel -h
```

run command:

```bash
go run ./multilevel
```

### Simple git commands

show help:

```bash
go run ./ggit -h
```

run command:

```bash
go run ./ggit
```
