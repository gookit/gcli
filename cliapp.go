package main

import (
    "runtime"
    "feedscenter/consumer/cli"
    "feedscenter/consumer/cmd"
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
