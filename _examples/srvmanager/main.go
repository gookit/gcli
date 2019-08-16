package main

import "github.com/gookit/gcli/v2"

func main() {
	app := gcli.NewApp()
	app.Version = "1.0.0"
	app.Description = "manage the http server start,stop,restart"

	app.Add(ServerStart(), ServerStop(), ServerRestart())
	app.Run()
}
