package main

import (
    "runtime"
    "feedscenter/console/cli"
    "feedscenter/console/cmd"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    app := cli.NewApp()
    app.Verbose = cli.VerbDebug
    //app.Description = "a Des"

    app.Add(cmd.TestCommand())
    app.Add(cmd.NewPusher())

    app.Run()
}
