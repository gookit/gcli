package main

import (
    "runtime"
    "github.com/golangkit/cliapp"
)

// for test run: go build console/cliapp.go && ./cliapp
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    app := cliapp.NewApp()
    app.Verbose = cliapp.VerbDebug
    app.Description = "this is my cli application"

    app.Add(ExampleCommand())
    app.Add(GitCommand())

    app.Run()
}
