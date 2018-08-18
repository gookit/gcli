package cliapp

import (
	"fmt"
	"github.com/gookit/color"
	"strings"
)

var level2name = map[uint]string{
	VerbError: "ERROR",
	VerbWarn:  "WARNING",
	VerbInfo:  "INFO",
	VerbDebug: "DEBUG",
	// VerbCrazy: "CRAZY",
}

var level2color = map[uint]color.Color{
	VerbError: color.FgRed,
	VerbWarn:  color.FgYellow,
	VerbInfo:  color.FgCyan,
	VerbDebug: color.FgGreen,
}

// Logf print log message
func Logf(level uint, format string, v ...interface{}) {
	if gOpts.verbose < level {
		return
	}

	name, has := level2name[level]
	if !has {
		return
	}

	c := level2color[level]
	fmt.Printf("cliapp: [%s] %s\n", c.Render(name), fmt.Sprintf(format, v...))
}

// replaceVars replace vars in the help info
func replaceVars(help string, vars map[string]string) string {
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

// strictFormatArgs '-ab' will split to '-a -b', '--o' -> '-o'
func strictFormatArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	var fmtdArgs []string
	for _, arg := range args {
		l := len(arg)

		if strings.Index(arg, "--") == 0 {
			if l == 3 {
				arg = "-" + string(arg[2])
			}

		} else if strings.Index(arg, "-") == 0 {
			if l > 2 {
				bools := strings.Split(strings.Trim(arg, "-"), "")
				for _, s := range bools {
					fmtdArgs = append(fmtdArgs, "-"+s)
				}

				continue
			}
		}

		fmtdArgs = append(fmtdArgs, arg)
	}

	return fmtdArgs
}
