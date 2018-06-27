package cliapp

import (
	"fmt"
	"os"
)

func Stdout(msg ...interface{})  {
	fmt.Fprint(os.Stdout, msg...)
}

func Stdoutf(f string, v ...interface{})  {
	fmt.Fprintf(os.Stdout, f + "\n", v...)
}

func Stderr(msg ...interface{})  {
	fmt.Fprint(os.Stderr, msg...)
}

func Stderrf(f string, v ...interface{})  {
	fmt.Fprintf(os.Stderr, f + "\n", v...)
}
