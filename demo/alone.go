package main

import "github.com/golangkit/cliapp"

var (
	fontName      string
	visualMode    bool
	list          bool
	sample        bool
)

// test run: go build ./demo/alone.go && ./alone -h
func main() {
	cmd := cliapp.Command{
		Name:        "test",
		Aliases:     []string{"ts"},
		Description: "this is a description <info>message</> for {$cmd}", // // {$cmd} will be replace to 'test'
		Fn: run,
	}

	cmd.Flags.BoolVar(&visualMode, "visual", false,"Prints the font name.")
	cmd.Flags.StringVar(&fontName, "font", "", "Choose a font name. Default is a random font.")
	cmd.Flags.BoolVar(&list, "list", false, "Lists all available fonts.")
	cmd.Flags.BoolVar(&sample, "sample", false, "Prints a sample with that font.")

	cmd.AloneRun()
}

func run(cmd *cliapp.Command, args []string) int  {
	cliapp.Stdout("hello, in the alone command\n")
	return 0
}
