# What's New in GCli v3.4 — A Friendly Tour Since v3.3.1

> [GCli](https://github.com/gookit/gcli) is a simple-to-use, full-featured
> command-line application & tool library for Go. The **v3.4** release line
> brings a batch of features focused on **developer experience** and
> **robustness**. This post walks through everything new since `v3.3.1`,
> from the things you'll use every day to the more advanced bits.

If you build CLI tools in Go, this update is worth a look. Let's start with the
highlights, then dig into each feature with real, runnable examples.

## Highlights at a glance

- 🧠 **Smarter shell completion** — zero-registration generation and a new
  **dynamic** mode that needs no script maintenance (bash / zsh / PowerShell).
- 🧅 **Command middleware** — `Command.Use()` / `App.Use()` for auth, logging,
  timing and other cross-cutting concerns.
- 🗂️ **Grouped help** — organize commands and options into titled sections via
  `Category`.
- 🏷️ **More flexible struct binding** — a new `field` tag rule plus automatic
  **anonymous struct** expansion.
- 💬 **Declarative interactive input** — collect a missing value with a single
  `Question`.
- ➖ **POSIX short-option merging** — `-aux` = `-a -u -x`, opt-in and safe.
- 🛡️ **Robustness fixes** — panics are no longer swallowed, `help <cmd>` works
  on first call, and more.
- ⚠️ **A few breaking changes** with a clear migration table (see the end).

---

## 1. Smarter shell completion

Shell completion used to mean generating a **static** script that hard-codes your
command and option names. The moment you add a command, the script is stale and
must be regenerated. v3.4 fixes this from both ends.

**Zero-registration static generation.** You no longer need to register the
`genac` command. Every app now has a built-in global option:

```bash
# generate a completion script for your shell, then source it
myapp --gen-completion bash > myapp.bash
source myapp.bash
# zsh / PowerShell also supported
myapp --gen-completion zsh  > _myapp
myapp --gen-completion pwsh > myapp.ps1
```

**Dynamic completion (zero maintenance).** By default the generated script is a
*thin* one: instead of hard-coding names, it calls back into your binary to ask
for candidates at completion time via the built-in `--in-completion` option.
Add a command tomorrow and Tab-completion just works — no regeneration needed.

**Value candidates for options.** Want completion to also suggest *values* for an
option? Give it a `Choices` list:

```go
c.StrOpt2(&format, "format", "output format",
    gflag.WithChoices("json", "yaml", "table"))
// typing `--format <Tab>` now suggests: json  yaml  table
```

The candidate computation is fully unit-tested; the shell glue (bash/zsh/pwsh)
is delegated to that single dynamic entry point, so behavior stays consistent.

## 2. Command middleware

Need to run auth checks, logging, or timing before a command's main logic —
without copy-pasting that code into every command? Middleware is here.

Register one or more handlers with `Use()`. They run **in registration order,
before** the command's main `Func`. If any handler returns an error, the chain
stops and the error propagates (the main `Func` is skipped).

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

Both `Command.Use()` and `App.Use()` return the receiver, so they chain nicely.
Apps with no middleware behave exactly as before.

## 3. Grouped help

As an app grows, a flat list of commands and options gets hard to scan. Now you
can group them under titled sections with a `Category`.

```go
// group commands
app.Add(&gcli.Command{Name: "migrate", Desc: "run db migrate", Category: "database"})
app.Add(&gcli.Command{Name: "serve",   Desc: "start http server"}) // default group

// group options
cmd.StrVar(&dsn, &gcli.CliOpt{Name: "db-dsn", Desc: "database dsn", Category: "database"})
cmd.StrOpt2(&port, "port", "bind port", gflag.WithCategory("network"))
```

Groups keep their first-appearance order; items inside a group are sorted by
name. When no category is set, output is identical to before — fully backward
compatible.

## 4. More flexible struct binding

GCli has long supported binding options straight from a struct. v3.4 adds a third
tag rule and anonymous-field support.

`FromStruct` now supports three rules, chosen via `c.FromStruct(ptr, ruleType)`:

- `gcli.TagRuleNamed` (default): `flag:"name=int0;shorts=i;required=true;desc=message"`
- `gcli.TagRuleSimple`: `flag:"desc;required;default;shorts"`
- `gcli.TagRuleField` **(new)**: use the **field name** (SnakeCase) as the option
  name, and read metadata from independent tag keys. **Anonymous nested structs
  are expanded automatically** — great for sharing a common option set.

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

The `field` rule is the most concise: the option name comes from the field, and
`desc` / `default` / `required` live in their own tag keys — easy to read and
maintain.

## 5. Declarative interactive input

Sometimes a value is required but the user forgot to pass it. Instead of writing
a manual collector, just attach a `Question`: when the option value is empty,
GCli prompts for it interactively (a built-in default collector).

```go
c.StrOpt2(&token, "token", "the access token",
    gflag.WithQuestion("Please input your access token: "))
```

```text
$ myapp deploy
Please input your access token: ▮
```

If you also set a custom `Collector`, it takes priority over `Question`.

## 6. POSIX short-option merging

Classic POSIX tools let you combine short flags: `-a -u -x` becomes `-aux`. GCli
now supports this — carefully, and **opt-in** so nothing changes by default.

Turn it on via `Config.EnhanceShort`, using the new self-documenting constants:

```go
c.ParserCfg().EnhanceShort = gcli.EnhanceShortMerge  // 1: -aux => -a -u -x
c.ParserCfg().EnhanceShort = gcli.EnhanceShortAttach // 2: also -Ostdout => -O stdout
```

| Level | Constant | Behavior |
|------|----------|----------|
| 0 | `EnhanceShortNone` | off (default) — full compatibility |
| 1 | `EnhanceShortMerge` | split a group **only when all members are bool** shorts |
| 2 | `EnhanceShortAttach` | also support value-attached form `-Ostdout` = `-O stdout` |

The key safety rule: a group is split **only if every character is a bool short
option**. Mixed forms like `-aO` (where `O` takes a value) are left untouched, so
value-taking shorts are never mis-parsed. (Strict mode now drives this same safe
path internally, replacing the old "blind split".)

## 7. Robustness fixes

Several long-standing rough edges were smoothed out:

- **Panics are no longer swallowed.** `gflag.Parser.Parse` used to print and
  ignore a recovered panic; it now returns it as an error so your code can react.
- **`help <command>` works on the first call.** Previously it could print
  `unknown input command "help"`.
- **`findSimilarCmd` no longer pollutes the command registry** with a phantom
  `help` entry on unknown-command runs.
- **`Command.Copy()` no longer clears the source command's hooks** (a shared
  pointer used to reset the original).

## ⚠️ Breaking changes & migration

A small number of cleanups require action if you depended on them:

| Before | After |
|---|---|
| `import ".../gcli/v3/helper"` | now internal — inline your own helper |
| `import ".../gcli/v3/gclicom"` | removed (was unused after the cliui migration) |
| global `--verbose 4` flag | env `GCLI_VERBOSE=debug`, or `gcli.SetVerbose(gcli.VerbDebug)` / `gcli.SetDebugMode()` |

Why remove `--verbose`? It bound to a per-app copy that the logger never read, so
it had no real effect — it only cluttered your app's option list. Control the log
level via the environment variable or code instead.

> Note: multiple `App` instances in one process now share the global options
> (verbose / help / version / strict / completion), as part of unifying them into
> a single source of truth.

## Upgrade

```bash
go get -u github.com/gookit/gcli/v3@latest
```

Then explore the runnable demos under [`_examples/cmd`](https://github.com/gookit/gcli/tree/master/_examples):
`struct-flag` (field tag + anonymous), `short-merge` (EnhanceShort), and
`ask-demo` (Question).

## Wrapping up

v3.4 is all about making GCli nicer to build with — completion that maintains
itself, middleware for the boring-but-necessary cross-cutting logic, cleaner help
output, and more flexible flag binding — while quietly fixing several robustness
issues underneath.

Give it a try, and if you hit anything or have ideas, issues and PRs are very
welcome on [GitHub](https://github.com/gookit/gcli). Happy CLI building! 🎉

---

*Links: [GitHub](https://github.com/gookit/gcli) ·
[GoDoc](https://pkg.go.dev/github.com/gookit/gcli/v3) ·
[Chinese version](zh-CN/gcli-v3.4-whats-new.md)*
