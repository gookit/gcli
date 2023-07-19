package events

// constants for hooks event, there are default allowed event names
const (
	// OnAppInitBefore On app init before
	OnAppInitBefore = "app.init.before"
	// OnAppInitAfter On app init after
	OnAppInitAfter = "app.init.after"
	// OnAppExit On app exit before
	OnAppExit = "app.exit"

	// OnAppBindOptsBefore before bind app options
	OnAppBindOptsBefore = "app.bind.opts.before"
	// OnAppBindOptsAfter after bind app options.
	//
	// support binding custom global options
	OnAppBindOptsAfter = "app.bind.opts.after"

	// OnAppCmdAdd on app cmd add
	OnAppCmdAdd = "app.cmd.add.before"

	// OnAppCmdAdded on app cmd added
	OnAppCmdAdded = "app.cmd.added"

	// OnAppOptsParsed event
	//
	// Data:
	// 	{args: app-args}
	OnAppOptsParsed = "app.opts.parsed"

	// OnAppPrepared prepare for run, after the OnAppOptsParsed
	OnAppPrepared = "app.run.prepared"

	// OnAppRunBefore app run before, after the OnAppPrepared
	OnAppRunBefore = "app.run.before"
	OnAppRunAfter  = "app.run.after"
	OnAppRunError  = "app.run.error"

	OnCmdInitBefore = "cmd.init.before"
	OnCmdInitAfter  = "cmd.init.after"

	// OnCmdNotFound on top-command or subcommand not found.
	//
	// Ctx:
	// 	{"name": name, "args": []string}
	OnCmdNotFound = "cmd.not.found"

	// OnAppCmdNotFound on top command not found.
	// ctx: {"name": name, "args": []string}
	OnAppCmdNotFound = "app.cmd.not.found"
	// OnCmdSubNotFound on subcommand not found.
	// ctx: {"name": name, "args": []string}
	OnCmdSubNotFound = "cmd.sub.not.found"

	// OnCmdOptParsed event
	//
	// Data:
	// 	{args: command-args}
	OnCmdOptParsed = "cmd.opts.parsed"

	// OnCmdRunBefore cmd run, flags has been parsed.
	OnCmdRunBefore = "cmd.run.before"
	// OnCmdRunAfter after cmd success run
	OnCmdRunAfter = "cmd.run.after"
	OnCmdRunError = "cmd.run.error"

	// OnCmdExecBefore cmd exec
	OnCmdExecBefore = "cmd.exec.before"
	OnCmdExecAfter  = "cmd.exec.after"
	OnCmdExecError  = "cmd.exec.error"

	// OnGlobalOptsParsed app or cmd parsed the global options
	//
	// Data:
	// 	{args: remain-args}
	OnGlobalOptsParsed = "gcli.gopts.parsed"
)
