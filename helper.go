package gcli

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
)

/*************************************************************
 * console log
 *************************************************************/

var level2name = map[uint]string{
	VerbError: "ERROR",
	VerbWarn:  "WARN",
	VerbInfo:  "INFO",
	VerbDebug: "DEBUG",
	VerbCrazy: "CRAZY",
}

var level2color = map[uint]color.Color{
	VerbError: color.FgRed,
	VerbWarn:  color.FgYellow,
	VerbInfo:  color.FgGreen,
	VerbDebug: color.FgCyan,
	VerbCrazy: color.FgMagenta,
}

// Logf print log message
func Logf(level uint, format string, v ...interface{}) {
	if gOpts.verbose < level {
		return
	}

	name, has := level2name[level]
	if !has {
		name = "CRAZY"
		level = VerbCrazy
	}

	name = level2color[level].Render(name)
	fmt.Printf("GCLI: [%s] %s\n", name, fmt.Sprintf(format, v...))
}

/*************************************************************
 * simple events manage
 *************************************************************/

// SimpleHooks struct
type SimpleHooks struct {
	// Hooks can setting some hooks func on running.
	hooks map[string]HookFunc
}

// On register event hook by name
func (sh *SimpleHooks) On(name string, handler HookFunc) {
	if handler != nil {
		// init map
		if sh.hooks == nil {
			sh.hooks = make(map[string]HookFunc)
		}

		sh.hooks[name] = handler
	}
}

// Add register on not exists hook.
func (sh *SimpleHooks) Add(name string, handler HookFunc) {
	if _, ok := sh.hooks[name]; !ok {
		sh.On(name, handler)
	}
}

// Fire event by name, allow with event data
func (sh *SimpleHooks) Fire(event string, data ...interface{}) {
	if handler, ok := sh.hooks[event]; ok {
		handler(data...)
	}
}

// ClearHooks clear hooks data
func (sh *SimpleHooks) ClearHooks() {
	sh.hooks = nil
}

func defaultErrHandler(data ...interface{}) {
	if len(data) == 2 && data[1] != nil {
		if err, ok := data[1].(error); ok {
			color.Error.Tips(err.Error())
			// fmt.Println(color.Red.Render("ERROR:"), err.Error())
		}
	}
}

/*************************************************************
 * app/cmd help vars
 *************************************************************/

// HelpVarFormat allow var replace on render help info.
// Default support:
// 	"{$binName}" "{$cmd}" "{$fullCmd}" "{$workDir}"
const HelpVarFormat = "{$%s}"

// HelpVars struct. provide string var function for render help template.
type HelpVars struct {
	// Vars you can add some vars map for render help info
	Vars map[string]string
}

// AddVar get command name
func (hv *HelpVars) AddVar(name, value string) {
	if hv.Vars == nil {
		hv.Vars = make(map[string]string)
	}

	hv.Vars[name] = value
}

// AddVars add multi tpl vars
func (hv *HelpVars) AddVars(vars map[string]string) {
	for n, v := range vars {
		hv.AddVar(n, v)
	}
}

// GetVar get a help var by name
func (hv *HelpVars) GetVar(name string) string {
	return hv.Vars[name]
}

// GetVars get all tpl vars
func (hv *HelpVars) GetVars() map[string]string {
	return hv.Vars
}

// ReplaceVars replace vars in the input string.
func (hv *HelpVars) ReplaceVars(input string) string {
	// if not use var
	if !strings.Contains(input, "{$") {
		return input
	}

	var ss []string
	for n, v := range hv.Vars {
		ss = append(ss, fmt.Sprintf(HelpVarFormat, n), v)
	}

	return strings.NewReplacer(ss...).Replace(input)
}

/*************************************************************
 * some helper methods
 *************************************************************/

// Print messages
func Print(args ...interface{}) {
	color.Print(args...)
}

// Println messages
func Println(args ...interface{}) {
	color.Println(args...)
}

// Printf messages
func Printf(format string, args ...interface{}) {
	color.Printf(format, args...)
}

func exitWithErr(format string, v ...interface{}) {
	color.Error.Tips(format, v...)
	Exit(ERR)
}

// func exitWithMsg(format string, v ...interface{}) {
// 	fmt.Printf(format, v...)
// 	Exit(0)
// }

// strictFormatArgs '-ab' will split to '-a -b', '--o' -> '-o'
func strictFormatArgs(args []string) (fmtArgs []string) {
	if len(args) == 0 {
		return args
	}

	for _, arg := range args {

		// eg: --a ---name
		if strings.Index(arg, "--") == 0 {
			farg := strings.TrimLeft(arg, "-")
			if rl := len(farg); rl == 1 { // fix: "--a" -> "-a"
				arg = "-" + farg
			} else if rl > 1 { // fix: "---name" -> "--name"
				arg = "--" + farg
			}
			// TODO No change remain OR remove like "--" "---"

		} else if strings.IndexByte(arg, '-') == 0 {
			ln := len(arg)
			// fix: "-abc" -> "-a -b -c"
			if ln > 2 {
				chars := strings.Split(strings.Trim(arg, "-"), "")

				for _, s := range chars {
					fmtArgs = append(fmtArgs, "-"+s)
				}
				continue
			}
		}

		fmtArgs = append(fmtArgs, arg)
	}

	return fmtArgs
}

// flags parser is flag#FlagSet.Parse(), so:
// - if args like: "arg0 arg1 --opt", will parse fail
// - if args convert to: "--opt arg0 arg1", can correctly parse
func moveArgumentsToEnd(args []string) []string {
	if len(args) < 2 {
		return args
	}

	var argEnd int
	for i, arg := range args {
		// strop on the first option
		if strings.IndexByte(arg, '-') == 0 {
			argEnd = i
			break
		}
	}

	// the first is an option
	if argEnd == -1 {
		return args
	}

	return append(args[argEnd:], args[0:argEnd]...)
}
