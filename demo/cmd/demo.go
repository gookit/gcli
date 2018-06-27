package cmd

import (
	"github.com/golangkit/cliapp"
	"github.com/golangkit/cliapp/color"
	"fmt"
)

type DemoOpts struct {
	id  int
	c   string
	dir string
}

type DemoCommand struct {
	Cmd cliapp.Command
	Opts DemoOpts
}

func (c *DemoCommand) Configure() *cliapp.Command  {
	cmd := &cliapp.Command{
		Name:        "demo",
		Description: "this is a description for demo",
		Aliases:     []string{"dm"},
		//Fn:     demoExec,
		ArgList: map[string]string{
			"arg0": "the first argument",
			"arg1": "the second argument",
		},
		Examples: "{{script}} {{cmd}} --id 12 -c val ag0 ag1",
	}

	f := &cmd.Flags
	f.IntVar(&c.Opts.id, "id", 2, "the id option")
	f.StringVar(&c.Opts.c, "c", "value", "the short option")
	f.StringVar(&c.Opts.dir, "dir", "", "the dir option")

	return cmd
}

func (c *DemoCommand) Execute(app *cliapp.App, args []string) int {
	fmt.Println(color.Render("<suc>hello, in the demo commander.</>"))

	return 0
}
