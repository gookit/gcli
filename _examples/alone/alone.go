package main

import (
	"fmt"
	"github.com/gookit/cliapp"
)

var opts = struct {
	fontName   string
	visualMode bool
	list       bool
	sample     bool
}{}

// test run: go build ./demo/alone.go && ./alone -h
func main() {
	cmd := cliapp.Command{
		Name:    "test",
		Aliases: []string{"ts"},
		UseFor:  "this is a description <info>message</> for {$cmd}", // // {$cmd} will be replace to 'test'
		Func:    run,
	}

	cmd.BoolOpt(&opts.visualMode, "visual", "v", false, "Prints the font name.")
	cmd.StrOpt(&opts.fontName, "font", "", "", "Choose a font name. Default is a random font.")
	cmd.BoolOpt(&opts.list, "list", "", false, "Lists all available fonts.")
	cmd.BoolOpt(&opts.sample, "sample", "", false, "Prints a sample with that font.")

	// Alone Running
	cmd.AloneRun()
}

func run(cmd *cliapp.Command, args []string) int {
	cliapp.Print("hello, in the alone command\n")

	// fmt.Printf("%+v\n", cmd.Flags)
	fmt.Printf("opts %+v\n", opts)
	fmt.Printf("args is %v\n", args)

	return 0
}
