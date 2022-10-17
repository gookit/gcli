package cmd

import (
	"fmt"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

// Names The string flag list, implemented flag.Value interface
type Names []string

func (ns *Names) String() string {
	return fmt.Sprint(*ns)
}

func (ns *Names) Set(value string) error {
	*ns = append(*ns, value)
	return nil
}

// options for the command
var exampleOpts = struct {
	id      int
	c       string
	dir     string
	opt     string
	showErr bool
	names   Names
}{}

// Example command definition
var Example = &gcli.Command{
	Func:    exampleExecute,
	Name:    "example",
	Aliases: []string{"module-exp", "exp", "ex"},
	Desc:    "this is command description message",
	// {$binName} {$cmd} is help vars. '{$cmd}' will replace to 'example'
	Examples: `
  {$binName} {$cmd} --id 12 -c val ag0 ag1
  <cyan>{$fullCmd} --names tom --names john -n c</> 	test use special option
`,
	Config: func(c *gcli.Command) {
		// bind options
		c.IntOpt(&exampleOpts.id, "id", "", 2, "the id option")
		c.BoolOpt(&exampleOpts.showErr, "err", "e", false, "display error example")
		c.StrOpt(&exampleOpts.c, "config", "c", "value", "the config option")
		// notice `DIRECTORY` will replace to option value type
		c.StrOpt(&exampleOpts.dir, "dir", "d", "", "the `DIRECTORY` option")
		// setting option name and short-option name
		c.StrOpt(&exampleOpts.opt, "opt", "o", "", "the option message")
		// setting a special option var, it must implement the flag.Value interface
		c.VarOpt(&exampleOpts.names, "names", "n", "the option message")

		// bind args with names
		c.AddArg("arg0", "the first argument, is required", true)
		c.AddArg("arg1", "the second argument, is required", true)
		c.AddArg("arg2", "the optional argument, is optional")
		c.AddArg("arrArg", "the array argument, is array", false, true)

	},
}

// command running
// example run:
//
//	go run ./_examples/cliapp.go ex -c some.txt -d ./dir --id 34 -n tom -n john val0 val1 val2 arrVal0 arrVal1 arrVal2
func exampleExecute(c *gcli.Command, args []string) error {
	color.Infoln("hello, in example command")

	if exampleOpts.showErr {
		return c.NewErrf("OO, An error has occurred!!")
	}

	magentaln := color.Magenta.Println

	color.Cyanln("All Aptions:")
	// fmt.Printf("%+v\n", exampleOpts)
	dump.V(exampleOpts)

	color.Cyanln("Remain Args:")
	// fmt.Printf("%v\n", args)
	dump.P(args)

	magentaln("Get arg by name:")
	arr := c.Arg("arg0")
	fmt.Printf("named arg '%s', value: %#v\n", arr.Name, arr.Value)

	magentaln("All named args:")
	for _, arg := range c.Args() {
		fmt.Printf("- named arg '%s': %+v\n", arg.Name, arg.Value)
	}

	return nil
}
