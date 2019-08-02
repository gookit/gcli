package main

import (
	"flag"
	"fmt"
	"os"
)

var int0 int
var str0 string

// go run ./_examples/rawflag.go -int 10 -str abc
// go run ./_examples/rawflag.go --int 10 --str abc
// go run ./_examples/rawflag.go --int=10 --str=abc
func main() {
	useNewFlagSet()

	fmt.Println("int:", int0, "str:", str0)
}

func useDefaultFlag() {
	flag.IntVar(&int0, "int", 0, "int opt")
	flag.StringVar(&str0, "str", "", "str opt")

	flag.Parse()
}

func useNewFlagSet() {
	f := flag.NewFlagSet("user", flag.ExitOnError)
	f.IntVar(&int0, "int", 0, "int opt")
	f.StringVar(&str0, "str", "", "str opt")

	_ = f.Parse(os.Args[1:])
}
