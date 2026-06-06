# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to semantic-ish versioning.

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
