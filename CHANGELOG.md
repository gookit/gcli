# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to semantic-ish versioning.

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
