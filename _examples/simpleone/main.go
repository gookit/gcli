package main

import (
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

var opts = struct {
	fontName   string
	visualMode bool
	list       bool
	sample     bool
}{}

// test run: go build ./_examples/simpleone && ./simpleone -h
// test run: go run ./_examples/simpleone
func main() {
	cmd := gcli.Command{
		Name:    "test",
		Aliases: []string{"ts"},
		Desc:    "this is a description <info>message</> for {$cmd}", // // {$cmd} will be replace to 'test'
	}

	cmd.BoolOpt(&opts.visualMode, "visual", "v", false, "Prints the font name.")
	cmd.StrOpt(&opts.fontName, "font", "", "", "Choose a font name. Default is a random font.")
	cmd.BoolOpt(&opts.list, "list", "", false, "Lists all available fonts.")
	cmd.BoolOpt(&opts.sample, "sample", "", false, "Prints a sample with that font.")

	cmd.Func = func(c *gcli.Command, args []string) error {
		c.Infoln("hello, in the alone command:", c.Name)

		dump.Print(args)
		dump.Print(opts)

		return nil
	}

	// Alone Running
	cmd.MustRun(nil)
	// cmd.Run(os.Args[1:])
}
