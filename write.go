package cliapp

import (
	"fmt"
	"os"
)

func Stdout(msg ...interface{}) (int, error) {
	return fmt.Fprint(os.Stdout, msg...)
}

func Stdoutf(f string, v ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stdout, f + "\n", v...)
}

func Stderr(msg ...interface{}) (int, error) {
	return fmt.Fprint(os.Stderr, msg...)
}

func Stderrf(f string, v ...interface{})  {
	fmt.Fprintf(os.Stderr, f + "\n", v...)
}
