package main

import (
	"github.com/gookit/gcli/v3/_examples/cmd"
)

// test run:
//  go build ./_examples/ggit && ./ggit -h
// test run:
//  go run ./_examples/ggit
func main() {
	// Alone Running
	cmd.GitCmd.MustRun(nil)
}
