# GCli v3.5 Updates: Changes Since v3.3.1

> [GCli](https://github.com/gookit/gcli) is a command-line application and tool library for Go.
> This post covers the main changes from `v3.3.1` to the recently released `v3.5` (including the v3.4 cycle). These updates focus on developer experience and underlying stability.

If you write CLI tools in Go, a few of these features might be useful. Here are the main changes.

## Key Updates

- **Shell completion**: Supports zero-registration generation and a dynamic completion mode (bash / zsh / PowerShell).
- **Command middleware**: Handle auth, logging, and other cross-cutting concerns via `Command.Use()` / `App.Use()`.
- **Grouped help**: Categorize commands and options into titled sections using `Category`.
- **Struct binding improvements**: Added a `field` tag rule and support for expanding anonymous nested structs.
- **Interactive input**: Automatically prompt for missing values via `Question`.
- **POSIX short-option merging**: Support for `-aux` splitting into `-a -u -x`.
- **Robustness fixes**: Panic handling, `help` command behavior, and more.
- A few breaking changes (migration guide at the end).

## 1. Shell Completion Improvements

Generating shell completion scripts previously required hardcoding command and option names, meaning adding a new command required regenerating the script. GCli v3.5 improves this workflow.

You no longer need to manually register the `genac` command. A built-in global option generates the static script directly:

```bash
# generate a completion script for your shell, then source it
myapp --gen-completion bash > myapp.bash
source myapp.bash

# zsh / PowerShell also supported
myapp --gen-completion zsh  > _myapp
myapp --gen-completion pwsh > myapp.ps1
```

Additionally, a **dynamic completion** mode is now available. The generated script no longer hardcodes names; instead, it calls back into your binary via the built-in `--in-completion` option to fetch candidates at completion time. New commands will immediately work with Tab-completion without regeneration.

For option values, you can define a list of candidates using `Choices`:

```go
c.StrOpt2(&format, "format", "output format",
    gflag.WithChoices("json", "yaml", "table"))
// typing `--format <Tab>` now suggests: json  yaml  table
```

## 2. Command Middleware

If you need to run auth checks or logging before a command's main logic without duplicating code across commands, you can use middleware.

Handlers registered with `Use()` run in order before the command's main `Func`. If a handler returns an error, the chain stops, and the error propagates upward.

```go
// command-level middleware
cmd.Use(func(c *gcli.Command, args []string) error {
    if os.Getenv("TOKEN") == "" {
        return c.NewErrf("missing TOKEN env")
    }
    return nil // return nil to continue the chain
})

// application-level middleware: applies before every command
app.Use(func(c *gcli.Command, args []string) error {
    gcli.Debugf("running command: %s", c.Name)
    return nil
})
```

Both `Command.Use()` and `App.Use()` return the receiver, supporting chaining. Apps without middleware behave exactly as before.

## 3. Grouped Help

When an app has many commands and options, the help output can become cluttered. You can now group them into titled sections using the `Category` field.

```go
// group commands
app.Add(&gcli.Command{Name: "migrate", Desc: "run db migrate", Category: "database"})
app.Add(&gcli.Command{Name: "serve",   Desc: "start http server"}) // default group

// group options
cmd.StrVar(&dsn, &gcli.CliOpt{Name: "db-dsn", Desc: "database dsn", Category: "database"})
cmd.StrOpt2(&port, "port", "bind port", gflag.WithCategory("network"))
```

Groups appear in the order of their first definition, and items within a group are sorted by name. If no category is set, the output format remains the same as in older versions.

## 4. More Flexible Struct Binding

`FromStruct` now supports a third tag rule (`TagRuleField`) and automatically expands anonymous nested structs.

The three available rules, selected via `c.FromStruct(ptr, ruleType)`:

- `gcli.TagRuleNamed` (default): `flag:"name=int0;shorts=i;required=true;desc=message"`
- `gcli.TagRuleSimple`: `flag:"desc;required;default;shorts"`
- `gcli.TagRuleField` (new): Uses the field name (SnakeCase) as the option name and reads metadata from independent tag keys.

```go
type commonOpts struct {
    Verbose bool `flag:"v" desc:"enable verbose output"`
}

type demoOpts struct {
    commonOpts        // anonymous: expands into a --verbose/-v option
    UserName string `flag:"u" desc:"the user name" required:"true"`
    Age      int    `desc:"the user age" default:"18"`
}

c.MustFromStruct(&demoOpts{}, gcli.TagRuleField)
// => options: --user-name/-u (required), --age (default 18), --verbose/-v
```

The `field` rule keeps the option name tied to the struct field, while `desc`, `default`, and `required` live in their own tags, making it easier to read and maintain.

## 5. Declarative Interactive Input

If a required option is missing at runtime, you can now attach a `Question`. GCli will detect the empty value and prompt the user for input interactively.

```go
c.StrOpt2(&token, "token", "the access token",
    gflag.WithQuestion("Please input your access token: "))
```

```text
$ myapp deploy
Please input your access token: ▮
```

If a custom `Collector` is also set, it takes priority over `Question`.

## 6. POSIX Short-Option Merging

GCli now supports merging short flags (e.g., `-a -u -x` becomes `-aux`) in a POSIX style. This is disabled by default and can be enabled via `Config.EnhanceShort`.

```go
c.ParserCfg().EnhanceShort = gcli.EnhanceShortMerge  // 1: -aux => -a -u -x
c.ParserCfg().EnhanceShort = gcli.EnhanceShortAttach // 2: also -Ostdout => -O stdout
```

It can also be enabled globally:
```go
gcli.SetEnhanceShort(gcli.EnhanceShortMerge)
```

| Level | Constant | Behavior |
|------|----------|----------|
| 0 | `EnhanceShortNone` | Off (default), fully compatible with old behavior |
| 1 | `EnhanceShortMerge` | Split a group only when all members are bool shorts |
| 2 | `EnhanceShortAttach` | Also support value-attached form `-Ostdout` = `-O stdout` |

A safety check is in place: a group is only split if every character is a boolean short option. Mixed forms like `-aO` (where `O` takes a value) are left untouched to prevent misparsing.

## 7. Robustness Fixes

Alongside new features, several long-standing issues were fixed:

- **Panics are no longer swallowed**: `gflag.Parser.Parse` previously ignored recovered panics. It now returns them as an error for easier upstream handling.
- **`help <command>` works on the first call**: Fixed an issue where it might print `unknown input command "help"`.
- **`findSimilarCmd` fix**: No longer writes a phantom `help` entry into the registry when an unknown command is run.
- **`Command.Copy()` fix**: No longer clears the source command's hooks due to a shared pointer.

## Breaking Changes & Migration

A few internal cleanups require adjustments if you depended on them:

| Before | After |
|---|---|
| `import ".../gcli/v3/helper"` | Now internal; inline your own helper |
| `import ".../gcli/v3/gclicom"` | Removed (unused after cliui migration) |
| Global `--verbose 4` flag | Env `GCLI_VERBOSE=debug`, or `gcli.SetVerbose(gcli.VerbDebug)` / `gcli.SetDebugMode()` |

The `--verbose` flag was removed because it bound to a per-app copy that the underlying logger never read, making it ineffective. Use the environment variable or code to control log levels.

Additionally, multiple `App` instances within the same process now share global options (verbose / help / version / strict / completion).

## Upgrade & Examples

```bash
go get -u github.com/gookit/gcli/v3@latest
```

The `_examples/cmd` directory in the repository includes runnable examples: `struct-flag` (field tags + anonymous structs), `short-merge` (short option merging), and `ask-demo` (interactive input).

> If you run into issues or have suggestions, feel free to open an issue or PR on [GitHub](https://github.com/gookit/gcli). For full API documentation, refer to [GoDoc](https://pkg.go.dev/github.com/gookit/gcli/v3).
