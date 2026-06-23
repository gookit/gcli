# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to semantic-ish versioning.

## [Unreleased]

### Changed

- **Duplicate-bind panics now include the command path.** When an option (or
  argument) is bound twice on the same command, the panic message now appends
  `in command '<path>'` (e.g. `redefined option flag 'dry-run' in command 'git
  branch'`), making the offending command easy to locate. The command's flag set
  is now named by its full path (`Command.Path()`); the flag-set name is only used
  for such diagnostics, so help output is unaffected.

## [v3.8.0] - 2026-06-22

### ⚠️ Breaking Changes

- **Per-app parse state moved out of `GlobalOpts` into a new `AppOptions` type.**
  The runtime fields `ShowHelp` / `ShowVersion` / `inCompletion` / `genCompletion`
  are no longer on `GlobalOpts`; each `App` (and standalone command) now owns its
  own `AppOptions`, so concurrent `App` instances no longer share these. `App.Opts()`
  still returns the process-level `*GlobalOpts` (so `app.Opts() == gcli.GOpts()` and
  `app.Opts().Verbose` are unchanged); use the new `App.AppOpts()` for the per-app
  state. Process-level config (`Verbose` / strict / `EnhanceShort` and the logger)
  stays in the package singleton, so log-level behavior is unaffected.

### Added

- **Three-tier option model with shared (inherited) options: `Command.SharedOpts()`**
  (≈ cobra's `PersistentFlags`). Options bound on `c.SharedOpts()` are inherited by
  the command **and all of its descendant commands**, sharing the same bound variable
  (the same `flag.Value`/pointer). This adds a *shared* middle tier between the existing
  *global* (app) and *local* (per-command) options, so a parent option like
  `--git-dir` can be written and parsed in any sub-command segment —
  `app top sub --git-dir /x` and (with arg reorder) `app top sub arg --git-dir /x`
  both work. Use `c.SharedOpts()` with any binder (`BoolOpt/StrOpt/Opt[T]/FromStruct/...`).
  A local option of the same name on a sub-command takes priority; `Required` on a
  shared option is validated at the executing (leaf) command. New gflag primitive
  `Parser.InheritOptsFrom(src, category...)` re-registers another parser's options by
  their underlying `flag.Value`. In a sub-command's help, options inherited from
  ancestors are grouped under an **`Inherited Options`** section (a command's own
  shared options render with its local options).
- **Command documentation generation: new `docgen` package + builtin `GenDoc` command**
  (≈ cobra's `cobra/doc`). `docgen` renders a single command or a whole app to
  **markdown** (`CmdMarkdown` / `AppMarkdown` / `MarkdownTree`) and **man pages (roff)**
  (`CmdMan` / `ManTree`). Add `builtin.GenDoc()` to an app, then run
  `./cliapp gendoc -f md|man -o ./docs` to export docs. Adds a `gflag.CliOpt.TypeName()`
  accessor for the option type name.

### Fixed

- **`FromStruct` now expands anonymous embedded structs of an *unexported type*.**
  The unexported-field-name skip used to run *before* the anonymous-field check, so
  an embedded field whose **type name is lowercase** (e.g. `baseFlags`, `commonOpts`
  — as used by the `struct-flag` demo and the docs) was silently dropped and its
  options were never generated. The anonymous check now runs first; the exported
  inner fields of such an embed are reflectable/addressable and bind correctly. This
  is independent of the tag rule (`TagRuleNamed` / `TagRuleSimple` / `TagRuleField`).

[v3.8.0]: https://github.com/gookit/gcli/compare/v3.7.0...v3.8.0

## [v3.7.0] - 2026-06-22

### Added

- **Struct-tag binding: more field types + `enum`.** `FromStruct` now binds
  native `[]string` / `[]int` / `[]bool` (repeatable, e.g. `--name a --name b`),
  `time.Duration` (e.g. `--ttl 1h30m`), and `map[string]string` (repeatable
  `--meta k=v`) directly — no need to declare the special `gflag.Strings`/`KVString`
  types. A new `enum:"a,b,c"` tag key sets the option's value candidates (for
  completion) and adds membership validation. Internally the struct binder no
  longer uses `unsafe`.
- **Generic option binding: `gflag.Opt[T]` / `gflag.BindVar[T]`.** A type-safe
  generic API that dispatches on the pointer type to the matching binder, so one
  call replaces the per-type `BoolVar/IntVar/StrVar/...`. Supports the same set
  of types as struct binding (scalars, `time.Duration`, slices, `map[string]string`,
  and any `flag.Value`). Example: `gflag.Opt(fs, &name, "name", "n", "tom", "user name")`.

[v3.7.0]: https://github.com/gookit/gcli/compare/v3.6.0...v3.7.0

## [v3.6.0] - 2026-06-21

### ⚠️ Breaking Changes

- **Renamed package `github.com/gookit/gcli/v3/events` → `github.com/gookit/gcli/v3/gevent`.**
  The directory, file, and package name were renamed for naming consistency with
  the other sub-packages (e.g. `gflag`). The old `events` package no longer
  exists; update your imports. The event-name constants themselves are unchanged
  (`OnAppInitAfter`, `OnCmdRunBefore`, ...).

### Added

- **Full event-name aliases on the `gcli` package.** Every event name is now
  exposed as a `gcli.Evt*` constant (1:1 with the `gevent.On*` names), so you can
  reference event names directly from `gcli` **without importing the event
  package**. New aliases: `EvtAppInitBefore`, `EvtAppExit`,
  `EvtAppBindOptsBefore`, `EvtAppBindOptsAfter`, `EvtAppCmdAdd`, `EvtAppCmdAdded`,
  `EvtAppOptsParsed`, `EvtAppHelpBefore`, `EvtAppHelpAfter`, `EvtCmdInitBefore`.
- **Auto-reorder of input args (`Config.DisableReorderArgs`, enabled by default).**
  Before parsing options, the input args are rearranged into the canonical
  `--options... arguments` form, so options written *after* positional arguments
  are still parsed instead of being silently dropped — e.g. `cmd arg --name tom`
  now works the same as `cmd --name tom arg`. A known value-taking option keeps
  its value (`--name tom`); bool options, `--opt=val`, negative-number tokens
  (`-5`), a lone `-`, and everything after `--` are handled correctly. In a
  multi-level app **only the final executed command's args are reordered** —
  reordering stops at a sub-command name, so parent/sub option sets never mix.
  Disable per parser via `gflag.WithReorderArgs(false)` or
  `Config.DisableReorderArgs = true` to restore the strict std-flag order.

### Changed

- Mixed `arguments` + `--options` input on a command no longer loses the options.
  This is a **behavior change** but strictly more permissive: any input that
  parsed before still parses the same; only previously-failing orders now succeed.

### Migration

| Before | After |
|---|---|
| `import "github.com/gookit/gcli/v3/events"` | `import "github.com/gookit/gcli/v3/gevent"` |
| `events.OnAppInitAfter` | `gevent.OnAppInitAfter` — or `gcli.EvtAppInit` (no import needed) |
| `events.OnCmdRunBefore` | `gevent.OnCmdRunBefore` — or `gcli.EvtCmdRunBefore` (no import needed) |

> Tip: prefer the `gcli.Evt*` aliases to drop the event-package import entirely.

[v3.6.0]: https://github.com/gookit/gcli/compare/v3.5.0...v3.6.0

## [v3.5.0] - 2026-06-06

**Highlights:** more flexible struct binding (a new `field` tag rule plus automatic anonymous-field expansion), declarative interactive input via `Question`, and opt-in POSIX short-option merging through `EnhanceShort` — configurable per command or globally with `gcli.SetEnhanceShort()`.

### Added

- **Struct binding: `TagRuleField` tag rule.** A new rule for `FromStruct` that
  uses the **field name** (SnakeCase) as the option name and reads metadata from
  independent tag keys (`flag` for shorts, plus `desc` / `default` / `required`).
  Select it via `c.FromStruct(ptr, gcli.TagRuleField)`.
- **Struct binding: anonymous field expansion.** Anonymous nested structs are now
  expanded automatically, so a shared option set can be embedded and reused.
- **Declarative interactive input: `CliOpt.Question`.** When an option value is
  empty, GCli can collect it via an interactive prompt (a built-in default
  collector). Set it with `gflag.WithQuestion("...")`. A custom `Collector` still
  takes priority over `Question`.
- **POSIX short-option enhancement: `Config.EnhanceShort`.** Opt-in combining of
  short options with self-documenting levels `EnhanceShortNone` (0, default),
  `EnhanceShortMerge` (1, `-aux` => `-a -u -x` when all are bool), and
  `EnhanceShortAttach` (2, also `-Ostdout` => `-O stdout`). A group is split only
  when **all** members are bool short options, so value-taking shorts are never
  mis-parsed.
- **Global `EnhanceShort` setting.** `gcli.SetEnhanceShort(level)` /
  `gcli.EnhanceShort()` apply a level to every command at once; a command's own
  `Config.EnhanceShort` still takes priority.
- **Demo commands** under `_examples/cmd`: `struct-flag` (field tag + anonymous),
  `short-merge` (EnhanceShort), `ask-demo` (Question).

### Changed

- Strict mode now drives the safe `EnhanceShort` path internally instead of the
  old "blind split", which used to mis-split value-taking short options (e.g.
  `-Ostdout`). `strictFormatArgs` is reduced to long-option normalization only.

[v3.5.0]: https://github.com/gookit/gcli/compare/v3.4.1...v3.5.0

## [v3.4.1] - 2026-06-05

### Added

- **Built-in shell completion (no `genac` registration needed).**
  - `--gen-completion <bash|zsh|pwsh>` statically generates a completion script
    and exits.
  - `--in-completion` computes completion candidates at runtime; generated
    scripts are *thin* and delegate to it, so they need no regeneration when
    commands change.
  - PowerShell (pwsh) dynamic completion via `Register-ArgumentCompleter`.
  - A silent completion mode suppresses lifecycle hooks so stdout only contains
    candidates / the script.
- **Option value candidates: `CliOpt.Choices`** (`gflag.WithChoices(...)`) feed
  value completion for an option.
- **Command middleware: `Command.Use(...)`** runs handlers in registration order
  before the command's main `Func`; any handler returning an error aborts the
  chain.
- **Application middleware: `App.Use(...)`** applies before every command.

### Fixed

- `doExecute`'s `recover()` is now a proper `defer`, so panics during command
  execution are actually caught.

[v3.4.1]: https://github.com/gookit/gcli/compare/v3.4.0...v3.4.1

## [v3.4.0] - 2026-06-04

### ⚠️ Breaking Changes

- **Removed public package `github.com/gookit/gcli/v3/helper`.** It has been
  moved to the internal package `internal/helper`. These were internal-only
  utilities (`IsGoodName`, `IsGoodCmdId`, `IsGoodCmdName`, `Panicf`,
  `RenderText`). If you imported it directly, inline your own helpers instead.
- **Removed unused public package `github.com/gookit/gcli/v3/gclicom`** (it was
  a leftover after the show/progress migration to `gookit/cliui` and had no
  in-tree usages).
- **Removed the global `--verbose` / `--verb` option.** It bound to a per-app
  copy that was never read by the logger, so it had no effect. Control the log
  level via the env var `GCLI_VERBOSE` (e.g. `GCLI_VERBOSE=debug`) or
  `gcli.SetVerbose()` / `gcli.SetDebugMode()` in code.

### Added

- **Command grouping in help.** Set `Command.Category` to group commands under
  a titled section in the application help. Uncategorized commands fall back to
  the default `Available Commands` group (output unchanged when no category is
  used).
- **Option grouping in help.** Set `CliOpt.Category` (or use
  `gflag.WithCategory("name")`) to group options under a titled sub-section.
  Uncategorized options render first with no sub-title (backward compatible).

### Fixed

- **`help <command>` now works on first invocation.** Previously it printed
  `unknown input command "help"` because `help` was not treated as a command.
- **`findSimilarCmd` no longer pollutes the command registry.** It used to write
  a phantom `help` entry into the real `cmdNames` map on any unknown-command run.
- **`Command.Copy()` no longer clears the source command's hooks** (the shallow
  copy shared the `*Hooks` pointer and reset the original).
- **`gflag.Parser.Parse` no longer silently swallows panics** — a recovered
  panic is now returned as an error instead of being printed and ignored.
- Fixed a compile error in the `_examples` progress demo after the
  `cliui/progress` `int64` signature change.

### Changed (internal)

- Global options are now a single source of truth: `App` reuses the package-level
  `gOpts`. **NOTE:** multiple `App` instances in the same process now share the
  global options (verbose/help/version/strict/completion).
- `App.findCommandName` is now side-effect free (returns a `foundCmd` instead of
  mutating `app.args` / `app.inputName` mid-parse).
- Merged the duplicated option-validation logic shared by `Parser.Parse` and
  `CliOpts.ParseOpts` into a single `validateAll`.
- Moved internal-only `helper` utilities under `internal/`.

### Migration

| Before | After |
|---|---|
| `import "github.com/gookit/gcli/v3/helper"` | internal now — inline your own helper |
| `import "github.com/gookit/gcli/v3/gclicom"` | removed |
| CLI flag `--verbose 4` | env `GCLI_VERBOSE=debug` or `gcli.SetVerbose(gcli.VerbDebug)` |

[v3.4.0]: https://github.com/gookit/gcli/compare/v3.3.1...v3.4.0

---

## v3.0.1 - 2021-04-23

**new**

- [x] add some special flag type vars
- [x] support hidden command on render help by `c.Hidden=true`

**fixed**

- [x] alias not works on command ID
- [x] render color on command/option/argument description

## v3.0.0 - 2021-04-23

**new**

- [x] support multi level sub commands
- [x] support parse flags from struct tags
- [x] support flag/argument validate
- [ ] support controller on application `app.controllers []Controller`
  - 独立于commands之外的。Independent of commands.
  - 支持组选项，全部子命令都拥有这些选项 `Config/GroupOptions()` 里绑定组选项。
- [x] 支持单个command、controller独立运行
