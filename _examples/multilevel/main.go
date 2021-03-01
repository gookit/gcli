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

var l1sub1opts = struct {
	aint int
}{}
var l2sub1opts = struct {
	astr string
}{}
var cmd = gcli.Command{
	Name:    "test",
	Aliases: []string{"ts"},
	Desc:    "this is a description <info>message</> for {$cmd}", // // {$cmd} will be replace to 'test'
	Subs: []*gcli.Command{
		{
			Name: "l1sub1",
			Desc: "desc message",
			Subs: []*gcli.Command{
				{
					Name: "l2sub1",
					Desc: "desc message",
					Config: func(c *gcli.Command) {
						c.StrVar(&l2sub1opts.astr, &gcli.FlagMeta{
							Name: "astr",
							Desc: "desc for astr",
						})
					},
				},
				{
					Name: "l2sub2",
					Desc: "desc message",
				},
			},
			Config: func(c *gcli.Command) {
				c.IntOpt(&l1sub1opts.aint, "aint", "", 2, "desc for aint")
			},
		},
		{
			Name: "l1sub2",
			Desc: "desc message",
		},
	},
}

// test run: go build ./_examples/multilevel && ./multilevel -h
// test run: go run ./_examples/multilevel
func main() {
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
