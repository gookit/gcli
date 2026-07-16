package main

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/_examples/cmd"
	"github.com/gookit/gcli/v3/builtin"
	"github.com/gookit/gcli/v3/gevent"
	// "github.com/gookit/gcli/v3/builtin/filewatcher"
	// "github.com/gookit/gcli/v3/builtin/reverseproxy"
)

var customGOpt string

// local run:
//
//	go run ./_examples/cliapp
//	go build ./_examples/cliapp && ./cliapp
//
// run on windows(cmd, powerShell):
//
//	go run ./_examples/cliapp
//	go build ./_examples/cliapp && ./cliapp
func main() {
	app := gcli.NewApp(func(app *gcli.App) {
		app.Version = "3.0.0"
		app.Desc = "this is my cli application"
		app.On(gcli.EvtAppInit, func(ctx *gcli.HookCtx) bool {
			// do something...
			gcli.Debugf("init app event", ctx.Name())
			return false
		})

		// app.SetVerbose(gcli.VerbDebug)
		// app.DefaultCommand("example")
		app.Logo.Text = `   ________    _______
  / ____/ /   /  _/   |  ____  ____
 / /   / /    / // /| | / __ \/ __ \
/ /___/ /____/ // ___ |/ /_/ / /_/ /
\____/_____/___/_/  |_/ .___/ .___/
                     /_/   /_/`
	})

	// disable global options
	// gcli.GOpts().SetDisable()

	// app.BeforeAddOpts = func(opts *gcli.Flags) {
	// 	opts.StrVar(&customGOpt, &gcli.CliOpt{Name: "custom", Desc: "desc message for the option"})
	// }

	app.On(gevent.OnAppBindOptsAfter, func(ctx *gcli.HookCtx) (stop bool) {
		gcli.Debugf("event OnAppBindOptsAfter", ctx.Name())
		return false
	})
	app.Flags().StrVar(&customGOpt, &gcli.CliOpt{
		Name: "custom",
		Desc: "desc message for the option",
	})

	// app.Strict = true
	app.Add(cmd.GitCmd)

	app.Add(cmd.Example)
	app.Add(cmd.DaemonRun)
	app.Add(cmd.EnvInfo)
	app.Add(cmd.CliColor, cmd.EmojiDemo)
	app.Add(
		cmd.ShowDemo,
		cmd.ProgressDemo,
		cmd.SpinnerDemo,
		cmd.InteractDemo,
	)

	// demos for recent features: B6(struct tag) / B4+B5(short merge) / B7(question)
	app.Add(
		cmd.StructFlagDemo,
		cmd.StructTypesDemo,
		cmd.ShortMergeDemo,
		cmd.QuestionDemo,
		cmd.ReorderDemo,
	)

	app.Add(builtin.GenEmojiMap)
	app.Add(builtin.GenDoc())
	app.Add(builtin.GenAutoComplete())

	// app.Add(filewatcher.FileWatcher(nil))
	// app.Add(reverseproxy.ReverseProxyCommand())

	var serverToken string
	app.Add(&gcli.Command{
		Name:    "test",
		Aliases: []string{"ts"},
		Desc:    "this is a description <info>message</> for command {$cmd}",
		Config: func(cmd *gcli.Command) {
			cmd.StrOpt(&serverToken, "token", "", "${SERVER_TOKEN}", "server token")
		},
		Func: func(cmd *gcli.Command, args []string) error {
			gcli.Print("hello, in the test command\n")
			return nil
		},
	})

	// create by func
	gcli.NewCommand("test1", "description1", func(c *gcli.Command) {
		// some config for the command
	}).WithFunc(func(c *gcli.Command, args []string) error {
		color.Green.Println("hello, command is: ", c.Name)
		return nil
	}).AttachTo(app)

	// fmt.Printf("%+v\n", gcli.CommandNames())

	// running
	app.Run(nil)
}
