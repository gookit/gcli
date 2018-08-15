package cliapp

import (
	"fmt"
	"strings"
)

var level2name = map[uint]string{
	VerbError: "ERROR",
	VerbWarn:  "WARNING",
	VerbInfo:  "INFO",
	VerbDebug: "DEBUG",
}

// print debug logging
func debugf(level uint, format string, v ...interface{}) {
	if Verbose < level {
		return
	}

	fmt.Printf("[DEBUG] "+format, v...)
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
