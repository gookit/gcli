package main

import (
	"fmt"

	"github.com/gookit/gcli/v3"
)

var opts = struct {
	fontName   string
	visualMode bool
	list       bool
	sample     bool
}{}
var cmd = gcli.Command{
	Name:    "test",
	Aliases: []string{"ts"},
	Desc:    "this is a description <info>message</> for {$cmd}", // // {$cmd} will be replace to 'test'
	Subs: []*gcli.Command{
		{
			Name: "sub1",
			Desc: "desc message",
			Subs: []*gcli.Command{
				{
					Name: "sub11",
					Desc: "desc message",
				},
				{
					Name: "sub12",
					Desc: "desc message",
				},
			},
		},
		{
			Name: "sub2",
			Desc: "desc message",
		},
	},
}

// test run: go build ./_examples/mutlilevel && ./mutlilevel -h
// test run: go run ./_examples/mutlilevel
func main() {
	cmd.BoolOpt(&opts.visualMode, "visual", "v", false, "Prints the font name.")
	cmd.StrOpt(&opts.fontName, "font", "", "", "Choose a font name. Default is a random font.")
	cmd.BoolOpt(&opts.list, "list", "", false, "Lists all available fonts.")
	cmd.BoolOpt(&opts.sample, "sample", "", false, "Prints a sample with that font.")

	cmd.Func = func(c *gcli.Command, args []string) error {
		c.Infoln("hello, in the alone command\n")

		// fmt.Printf("%+v\n", cmd.Flags)
		fmt.Printf("opts %+v\n", opts)
		fmt.Printf("args is %v\n", args)

		return nil
	}

	// Alone Running
	cmd.MustRun(nil)
	// cmd.Run(os.Args[1:])
}
