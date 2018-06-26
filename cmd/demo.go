package cmd

import (
    "feedscenter/console/cli"
    "fmt"
)

type DemoCommand struct {
    cli.Command
    opts demoOptions
}

type demoOptions struct {
    id  int
    dir string
}

func NewDemoCommand() *DemoCommand {
    c := DemoCommand{Command: cli.Command{
        Name:        "demo",
        Description: "this is a description message",
        Aliases:     []string{"dm", "demo1"},
    }}

    // command definition
    c.Definition()

    return &c
}

func (c *DemoCommand) Name() string {
    return c.Command.Name
}

func (c *DemoCommand) Definition() {
    c.
        IntOpt(&c.opts.id, "id", 0, "the input id").
        StrOpt(&c.opts.dir, "d", "", "search `directory` for include files")
}

// execute
func (c *DemoCommand) Execute(app *cli.App, args []string) int {

    fmt.Printf("opts %+v\n", c.opts)

    return 0
}
