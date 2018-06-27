package cliapp

import (
	"strings"
	"fmt"
)

// AddVar get command name
func (app *Application) AddVar(name string, value string) {
	app.vars[name] = value
}

// AddVars add multi tpl vars
func (app *Application) AddVars(vars map[string]string) {
	for n, v := range vars {
		app.AddVar(n, v)
	}
}

// GetVars get all tpl vars
func (app *Application) GetVars(name string, value string) map[string]string {
	return app.vars
}

// HelpVar like "{$script}"
const HelpVar = "{$%s}"

// ReplaceVars replace vars in the help info
func ReplaceVars(help string, vars map[string]string) string {
	// if not use var
	if !strings.Contains(help, "{$") {
		return help
	}

	var ss []string

	for n, v := range vars {
		ss = append(ss, fmt.Sprintf(HelpVar, n), v)
	}

	return strings.NewReplacer(ss...).Replace(help)
}
