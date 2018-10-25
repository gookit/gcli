package main

import (
	"github.com/gookit/cliapp"
	"github.com/gookit/cliapp/_examples/cmd"
	"github.com/gookit/cliapp/builtin"
	"github.com/gookit/color"

	// "github.com/gookit/cliapp/builtin/filewatcher"
	// "github.com/gookit/cliapp/builtin/reverseproxy"
	"runtime"
)

// run:
// go run ./_examples/cliapp.go
// go build ./_examples/cliapp.go && ./cliapp
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app := cliapp.NewApp(func(app *cliapp.App) {
		app.Version = "1.0.6"
		app.Description = "this is my cli application"
		app.Hooks[cliapp.EvtInit] = func(a *cliapp.App, data interface{}) {
			// do something...
			// fmt.Println("init app")
		}
		// app.SetVerbose(cliapp.VerbDebug)
		// app.DefaultCommand("example")
		app.Logo.Text = `   ________    _______
  / ____/ /   /  _/   |  ____  ____
 / /   / /    / // /| | / __ \/ __ \
/ /___/ /____/ // ___ |/ /_/ / /_/ /
\____/_____/___/_/  |_/ .___/ .___/
                     /_/   /_/`
	})

	app.Add(cmd.ExampleCommand())
	app.Add(cmd.EnvInfoCommand())
	app.Add(cmd.GitCommand())
	app.Add(cmd.ColorCommand(), cmd.EmojiDemoCmd())
	app.Add(cmd.ShowDemoCommand(), cmd.ProgressDemoCmd(), cmd.SpinnerDemoCmd(), cmd.InteractDemoCommand())
	app.Add(builtin.GenEmojiMapCommand())

	// app.Add(filewatcher.FileWatcher(nil))
	// app.Add(reverseproxy.ReverseProxyCommand())

	app.Add(&cliapp.Command{
		Name:    "test",
		Aliases: []string{"ts"},
		UseFor:  "this is a description <info>message</> for command {$cmd}",
		Func: func(cmd *cliapp.Command, args []string) int {
			cliapp.Print("hello, in the test command\n")
			return 0
		},
	})

	app.Add(builtin.GenAutoCompleteScript())
	// create by func
	app.NewCommand("test1", "description1", func(c *cliapp.Command) {
		// some config for the command
	}).SetFunc(func(c *cliapp.Command, args []string) int {
		color.Green.Println("hello, command is: ", c.Name)
		return 0
	}).AttachTo(app)

	// fmt.Printf("%+v\n", cliapp.CommandNames())

	// running
	app.Run()
}
