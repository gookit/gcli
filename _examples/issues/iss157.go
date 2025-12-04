package main

import (
	"github.com/gookit/color/colorp"
	"github.com/gookit/gcli/v3"
)

type userOpts struct {
	Int  int    `flag:"name=int0;shorts=i;required=true;desc=int option message"`
	Bol  bool   `flag:"name=bol;shorts=b;desc=bool option message"`
	Str1 string `flag:"name=str1;shorts=o;required=true;desc=str1 message"`
	// use ptr
	Str2 *string `flag:"name=str2;required=true;desc=str2 message"`
	// custom type and implement flag.Value
	Verb0 gcli.VerbLevel `flag:"name=verb0;shorts=v0;desc=verb0 message"`
	// use ptr
	Verb1 *gcli.VerbLevel `flag:"name=verb1;desc=verb1 message"`
}

// run: go run ./_examples/issues/iss157.go
func main() {
	astr := "xyz"
	verb := gcli.VerbWarn

	cmd := gcli.NewCommand("test", "desc")
	cmd.Config = func(c *gcli.Command) {
		c.MustFromStruct(&userOpts{
			Str2:  &astr,
			Verb1: &verb,
		})
	}

	// disable auto bind global options: verbose,version, progress...
	gcli.GOpts().SetDisable()

	// direct run
	if err := cmd.Run(nil); err != nil {
		colorp.Errorln( err)
	}
}