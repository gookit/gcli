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
	number     int
}{}

// test run:
// 	go build ./_examples/simpleone && ./simpleone -h
// test run:
// 	go run ./_examples/simpleone
func main() {
	cmd := gcli.Command{
		Name:    "test",
		Aliases: []string{"ts"},
		Desc:    "this is a description <info>message</> for {$cmd}", // // {$cmd} will be replace to 'test'
	}

	cmd.BoolOpt(&opts.visualMode, "visual", "v", false, "Prints the font name.")
	cmd.StrOpt(&opts.fontName, "font", "fn", "", "Choose a font name. Default is a random name;true")
	cmd.BoolOpt(&opts.list, "list", "", false, "Lists all available fonts.")
	cmd.BoolOpt(&opts.sample, "sample", "", false, "Prints a sample with that font.\nmessage at new line")
	cmd.IntOpt(&opts.number, "number", "n,num", 0, "a integer option")

	cmd.AddArg("arg1", "this is a argument")
	cmd.AddArg("arg2", "this is a argument2")

	cmd.WithConfigFn(func(opt *gcli.FlagsConfig) {
		opt.DescNewline = true
	})

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
