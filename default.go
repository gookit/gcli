package gcli

var (
	// store default application instance
	DefaultApp *App
)

// NewDefaultApp create the default cli app.
func NewDefaultApp(fn ...func(a *App)) *App {
	DefaultApp = NewApp(fn...)
	return DefaultApp
}

// AllCommands returns all commands in the default app
func AllCommands() map[string]*Command {
	return DefaultApp.Commands()
}
