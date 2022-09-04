package events

// constants for hooks event, there are default allowed event names
const (
	OnAppInit = "app.init"

	OnAppPrepareAfter = "app.prepare.after"

	OnAppRunBefore = "app.run.before"
	OnAppRunAfter  = "app.run.after"
	OnAppRunError  = "app.run.error"

	OnCmdInit = "cmd.init"

	// OnCmdNotFound app or sub command not found.
	//
	// Data:
	// 	{name: command-name}
	OnCmdNotFound = "cmd.not.found"
	// OnAppCmdNotFound app command not found
	OnAppCmdNotFound = "app.cmd.not.found"
	// OnCmdSubNotFound sub command not found
	OnCmdSubNotFound = "cmd.sub.not.found"

	// OnCmdOptParsed event
	//
	// Data:
	// 	{args: command-args}
	OnCmdOptParsed = "cmd.opts.parsed"

	// OnCmdRunBefore cmd run
	OnCmdRunBefore = "cmd.run.before"
	OnCmdRunAfter  = "cmd.run.after"
	OnCmdRunError  = "cmd.run.error"

	// OnCmdExecBefore cmd exec
	OnCmdExecBefore = "cmd.exec.before"
	OnCmdExecAfter  = "cmd.exec.after"
	OnCmdExecError  = "cmd.exec.error"

	// OnGOptionsParsed event
	//
	// Data:
	// 	{args: remain-args}
	OnGOptionsParsed = "gcli.gopts.parsed"
	// OnStop   = "stop"
)
