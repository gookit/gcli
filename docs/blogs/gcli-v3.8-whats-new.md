# What's New in GCli v3.8 — Since v3.5

> [GCli](https://github.com/gookit/gcli) is a simple-to-use, full-featured command-line application & tool library
> for Go. The **v3.6 → v3.8** cycle is a focused modernization pass: friendlier parsing, richer type binding, a
> type-safe generic API, a proper three-tier option model, and built-in documentation generation. This post walks
> through everything new since `v3.5`, organized by capability rather than by tag.

If you build CLI tools in Go, these three releases close several long-standing gaps versus cobra/kong while keeping
GCli's "batteries-included" character. Let's start with the highlights, then dig into each feature with real, runnable
examples.

## Highlights at a glance

- 🔀 **Auto-reorder of input args** (v3.6) — options written *after* positional arguments are still parsed; on by default, safe for multi-level commands.
- 🧰 **Richer struct-tag binding** (v3.7) — native `[]string`/`[]int`/`[]bool`, `time.Duration`, `map[string]string`, and an `enum:"a,b,c"` tag — and the binder no longer uses `unsafe`.
- 🧬 **Type-safe generic API** (v3.7) — `gflag.Opt[T]` / `gflag.BindVar[T]` replace the per-type `BoolVar/IntVar/StrVar/...`.
- 🪜 **Three-tier option model** (v3.8) — `Command.SharedOpts()` (≈ cobra's `PersistentFlags`) are inherited by a command and all its descendants.
- 📄 **Command documentation generation** (v3.8) — export **markdown** and **man pages** via the new `docgen` package or the builtin `gendoc` command.
- 🏷️ **Package rename `events` → `gevent`** (v3.6), plus `gcli.Evt*` aliases so you can reference event names without importing the event package.
- ⚠️ **A couple of breaking changes** with a clear migration table (see the end).

---

## 1. Auto-reorder of input args

GCli used to require the canonical `--options ... arguments` order, just like the standard library `flag`: the first
non-flag token stops option parsing, so `cmd arg --name tom` would silently drop `--name tom`. In practice users mix
the order all the time. As of v3.6 GCli rearranges the input into the canonical form **before** parsing — so options
written after positional arguments are still picked up.

```bash
# these now behave identically — the option is no longer lost
myapp build --name tom src/
myapp build src/ --name tom
```

It is **on by default** and strictly more permissive: any input that parsed before parses the same; only
previously-failing orders now succeed. The reorder is careful — a known value-taking option keeps its value
(`--name tom`), while bool options, `--opt=val`, negative-number tokens (`-5`), a lone `-`, and everything after `--`
are all handled correctly.

Crucially, in a **multi-level app only the final executed command's args are reordered** — reordering stops at a
sub-command name, so a parent's and a sub-command's option sets never get mixed up.

Prefer the strict std-flag order? Turn it off per parser:

```go
// disable for one command
c.ParserCfg().DisableReorderArgs = true
// or via the config func
gflag.WithReorderArgs(false)
```

See the runnable [`reorder-args`](https://github.com/gookit/gcli/tree/master/_examples/cmd/reorder_demo.go) demo.

## 2. Richer struct-tag binding (and no more `unsafe`)

Binding options from a struct is one of GCli's nicest features. v3.7 makes the field types much richer — the common
collection and time types now bind **natively**, with no need to declare the special `gflag.Strings` / `KVString`
helper types:

```go
type deployOpts struct {
    Names []string          `flag:"name=names;shorts=n;desc=name list (repeatable)"`
    Ports []int             `flag:"name=ports;shorts=p;desc=port list (repeatable)"`
    TTL   time.Duration     `flag:"name=ttl;desc=time to live, eg: 1h30m"`
    Meta  map[string]string `flag:"name=meta;shorts=m;desc=key=value metadata (repeatable)"`
    Lang  string            `flag:"name=lang;shorts=l;desc=language;enum=go,php,java"`
}

c.MustFromStruct(&deployOpts{})
```

```bash
myapp deploy -n a -n b -p 80 -p 443 --ttl 1h30m -m k1=v1 -m k2=v2 -l go
# slices repeat (-n a -n b), maps repeat (-m k=v), duration parses 1h30m, lang must be one of go/php/java
```

What's new in the binder:

- **Slices** — `[]string` / `[]int` / `[]bool` bind as repeatable options (`--name a --name b`).
- **Duration** — `time.Duration` parses Go duration strings like `1h30m`.
- **Maps** — `map[string]string` binds as a repeatable `--meta k=v` option.
- **Enums** — the new `enum:"a,b,c"` tag key sets the option's value candidates (used for completion) **and** adds membership validation, so an out-of-set value is rejected.

Under the hood, the struct binder **no longer uses `unsafe`** — it now resolves the field pointer via the safe
`reflect.Value.Addr().Interface()`, closing the last unsafe path in struct binding. Try the
[`struct-types`](https://github.com/gookit/gcli/tree/master/_examples/cmd/structtypes_demo.go) demo.

## 3. A type-safe generic API

The classic per-type binders (`BoolVar`, `IntVar`, `StrVar`, `Float64Var`, ...) still work, but v3.7 adds a single
generic entry point that infers the binder from your pointer's type:

```go
var (
    name string
    age  int
    tags []string
    ttl  time.Duration
)

gflag.Opt(&c.Flags, &name, "name", "n", "tom", "the user name")
gflag.Opt(&c.Flags, &age,  "age",  "a", 18,    "the user age")
gflag.Opt(&c.Flags, &tags, "tag",  "t", nil,   "the tags, repeatable")
gflag.Opt(&c.Flags, &ttl,  "ttl",  "",  time.Duration(0), "time to live")
```

`gflag.Opt[T]` dispatches on the pointer type to the matching binder, covering the same set of types as struct binding
— scalars, `time.Duration`, slices, `map[string]string`, and any `flag.Value`. For full control over the option
metadata, use `gflag.BindVar[T]` with a `*gflag.CliOpt`:

```go
var langs []string
gflag.BindVar(&c.Flags, &langs, gflag.NewOpt("langs", "language list", nil))
```

One call, no per-type method names, full type safety from the pointer.

## 4. Three-tier option model: shared (inherited) options

This is the headline of v3.8. GCli has always had **global** (app-level) options and **local** (per-command) options.
v3.8 adds the missing middle tier — **shared options** via `Command.SharedOpts()` — the direct equivalent of cobra's
`PersistentFlags`.

Options bound on `c.SharedOpts()` are inherited by the command **and all of its descendant commands**, sharing the
same bound variable (the same underlying `flag.Value` / pointer). So a parent option can be written and parsed in any
sub-command segment:

```go
var gitDir string

top := &gcli.Command{Name: "git", Desc: "git-like demo"}
// bind a shared option on the parent — visible to every sub-command
top.SharedOpts().StrOpt(&gitDir, "git-dir", "", ".git", "the git data dir")

top.Add(&gcli.Command{
    Name: "status",
    Func: func(c *gcli.Command, _ []string) error {
        // gitDir is populated whether --git-dir is written here or on the parent
        gcli.Printf("git dir: %s\n", gitDir)
        return nil
    },
})
```

```bash
myapp git --git-dir /x status      # written on the parent
myapp git status --git-dir /x      # written on the sub-command — both work
myapp git status arg --git-dir /x  # and (thanks to arg reorder) even after an argument
```

The semantics mirror cobra closely:

- A **local** option of the same name on a sub-command takes priority over the inherited one.
- A `Required` shared option is validated at the **executing (leaf) command**, not at every intermediate ancestor.
- In a sub-command's help, options inherited from ancestors are grouped under an **`Inherited Options`** section; a
  command's *own* shared options render alongside its local options.

Under the hood this is powered by a new gflag primitive, `Parser.InheritOptsFrom(src, category...)`, which
re-registers another parser's options by their underlying `flag.Value` — so the parent and child genuinely write the
same variable.

## 5. Command documentation generation (markdown + man)

v3.8 ships a documentation generator, the rough equivalent of cobra's `cobra/doc`. The new `docgen` package renders a
single command or a whole app to **markdown** and **man pages (roff)**:

```go
import "github.com/gookit/gcli/v3/docgen"

// whole-app trees
docgen.MarkdownTree(app, "./docs")  // index.md + one .md per command
docgen.ManTree(app, "./man")        // one .1 per command

// single command
md  := docgen.CmdMarkdown(cmd)
man := docgen.CmdMan(cmd)
```

Prefer to drive it from the CLI? Add the builtin command and run it:

```go
import "github.com/gookit/gcli/v3/builtin"

app.Add(builtin.GenDoc())
```

```bash
./cliapp gendoc -f md  -o ./docs   # markdown
./cliapp gendoc -f man -o ./man    # man pages
```

A few niceties baked in:

- **Examples are cleaned and rendered** — color tags like `<cyan>...</>` are stripped and built-in variables such as
  `{$fullCmd}` are expanded, so docs read as plain runnable commands.
- **Multi-line Examples are preserved** — man output wraps examples in a `.nf/.fi` (no-fill) block, so each example
  line survives instead of being folded into one.
- The app overview (`index.md`) includes the app **version**, and option tables include each option's type via the new
  `gflag.CliOpt.TypeName()` accessor.

## 6. Package rename: `events` → `gevent`, plus `gcli.Evt*` aliases

For naming consistency with the other sub-packages (`gflag`, `gevent`), the event package was renamed from
`github.com/gookit/gcli/v3/events` to `github.com/gookit/gcli/v3/gevent`. The event-name constants themselves are
unchanged (`OnAppInitAfter`, `OnCmdRunBefore`, ...).

Even better: every event name is now also exposed as a `gcli.Evt*` constant, so you can reference event names directly
from the `gcli` package **without importing the event package at all**:

```go
// no event-package import needed
app.On(gcli.EvtAppInit,      func(ctx *gcli.HookCtx) bool { /* ... */ return false })
app.On(gcli.EvtCmdRunBefore, func(ctx *gcli.HookCtx) bool { /* ... */ return false })
```

## ⚠️ Breaking changes & migration

Two changes across this cycle may require action:

| Before | After |
|---|---|
| `import "github.com/gookit/gcli/v3/events"` | `import "github.com/gookit/gcli/v3/gevent"` (constants unchanged) |
| `events.OnCmdRunBefore` | `gevent.OnCmdRunBefore` — or `gcli.EvtCmdRunBefore` (no import needed) |
| per-app parse state on `GlobalOpts` (`ShowHelp`/`ShowVersion`/`inCompletion`/`genCompletion`) | moved to a new per-app `AppOptions`; use `app.AppOpts()` |

About the `AppOptions` split (v3.8): the runtime fields that describe *one app's* parse state moved out of the
process-level `GlobalOpts` into a per-`App` `AppOptions`, so concurrent `App` instances no longer share them.
**`App.Opts()` still returns the process-level `*GlobalOpts`** (so `app.Opts().Verbose` and `app.Opts() == gcli.GOpts()`
are unchanged) — only the four per-app fields above moved. Process-level config (verbose / strict / `EnhanceShort` and
the logger) deliberately stays in the package singleton, so log-level behavior is unaffected.

> Note: the arg auto-reorder (v3.6) is a behavior change, but a strictly more permissive one — if you relied on the old
> "stop at first non-flag" behavior, disable it with `Config.DisableReorderArgs = true`.

## Upgrade

```bash
go get -u github.com/gookit/gcli/v3@latest
```

Then explore the runnable demos under [`_examples/cmd`](https://github.com/gookit/gcli/tree/master/_examples):
`reorder-args` (arg reorder), `struct-types` (slice/duration/map/enum), and `struct-flag` (field tag + anonymous).

## Wrapping up

v3.6 → v3.8 is about catching up on the fundamentals that heavy CLI users expect: forgiving argument order, rich and
safe struct binding, a clean generic API, persistent/shared options across a command tree, and first-class doc
generation — all while keeping GCli's batteries-included color/interactive/progress stack intact.

Give it a try, and if you hit anything or have ideas, issues and PRs are very welcome on
[GitHub](https://github.com/gookit/gcli). Happy CLI building! 🎉

---

*Links: [GitHub](https://github.com/gookit/gcli) ·
[GoDoc](https://pkg.go.dev/github.com/gookit/gcli/v3) ·
[Chinese version](zh-CN/gcli-v3.8-whats-new.md)*
