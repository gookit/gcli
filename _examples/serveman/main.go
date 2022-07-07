package main

import "github.com/gookit/gcli/v3"

func main() {
	app := gcli.NewApp()
	app.Version = "1.0.0"
	app.Desc = "manage the http server start,stop,restart"

	app.Add(ServerStart, ServerStop, ServerRestart)
	app.Run(nil)
}
