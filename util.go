package gcli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/helper"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/stdutil"
)

/*************************************************************
 * console log
 *************************************************************/

var level2color = map[VerbLevel]color.Color{
	VerbError: color.FgRed,
	VerbWarn:  color.FgYellow,
	VerbInfo:  color.FgGreen,
	VerbDebug: color.FgCyan,
	VerbCrazy: color.FgMagenta,
}

// Debugf print log message
func Debugf(format string, v ...any) {
	logf(VerbDebug, format, v...)
}

// Logf print log message
func Logf(level VerbLevel, format string, v ...any) {
	logf(level, format, v...)
}

// print log message
func logf(level VerbLevel, format string, v ...any) {
	if gOpts.Verbose < level {
		return
	}

	name := level2color[level].Render(level.Upper())
	logAt := stdutil.GetCallerInfo(3)

	color.Printf("GCli: [%s] [<gray>%s</>] %s \n", name, logAt, fmt.Sprintf(format, v...))
}

func defaultErrHandler(ctx *HookCtx) (stop bool) {
	if errV := ctx.Get("err"); errV != nil {
		if err, ok := errV.(error); ok {
			color.Error.Tips(err.Error())
			// fmt.Println(color.Red.Render("ERROR:"), err.Error())
		}
	}

	return
}

func name2verbLevel(name string) VerbLevel {
	switch strings.ToLower(name) {
	case "quiet":
		return VerbQuiet
	case "error":
		return VerbError
	case "warn":
		return VerbWarn
	case "info":
		return VerbInfo
	case "debug":
		return VerbDebug
	case "crazy":
		return VerbCrazy
	}

	// default level
	return defaultVerb
}

/*************************************************************
 * some helper methods
 *************************************************************/

// Print messages
func Print(args ...any) {
	color.Print(args...)
}

// Println messages
func Println(args ...any) {
	color.Println(args...)
}

// Printf messages
func Printf(format string, args ...any) {
	color.Printf(format, args...)
}

func panicf(format string, v ...any) {
	panic(fmt.Sprintf("GCli: "+format, v...))
}

func sepStr(seps []string) string {
	if len(seps) > 0 {
		return seps[0]
	}
	return comdef.DefaultSep
}

const (
	// match a good option, argument name
	regGoodName = `^[a-zA-Z][\w-]*$`
	// match a good command name
	regGoodCmdName = `^[a-zA-Z][\w-]*$`
	// match command id. eg: "self:init"
	regGoodCmdId = `^[a-zA-Z][\w:-]*$`
	// match command path. eg: "self init"
	// regGoodCmdPath = `^[a-zA-Z][\w -]*$`
)

var (
	// good name for option and argument
	goodName = regexp.MustCompile(regGoodName)
	// match a good command name
	goodCmdId = regexp.MustCompile(regGoodCmdId)
	// match a good command name
	goodCmdName = regexp.MustCompile(regGoodCmdName)
)

func aliasNameCheck(name string) {
	if helper.IsGoodCmdName(name) {
		return
	}
	panicf("alias name '%s' is invalid, must match: %s", name, regGoodCmdName)
}

// strictFormatArgs
// TODO mode:
//
//	POSIX '-ab' will split to '-a -b', '--o' -> '-o'
//	UNIX '-ab' will split to '-a b'
func strictFormatArgs(args []string) (fmtArgs []string) {
	if len(args) == 0 {
		return args
	}

	for _, arg := range args {
		// if contains '=' append self
		// TODO mode:
		//  '--test=x', '-t=x' , '-test=x', '-test'
		if strings.ContainsRune(arg, '=') {
			fmtArgs = append(fmtArgs, arg)
			continue
		}

		// eg: --a ---name
		if strings.HasPrefix(arg, "--") {
			farg := strings.TrimLeft(arg, "-")
			if rl := len(farg); rl == 1 { // fix: "--a" -> "-a"
				arg = "-" + farg
			} else if rl > 1 { // fix: "---name" -> "--name"
				arg = "--" + farg
			}

			// TODO No change remain OR remove like "--" "---"
			// maybe ...

		} else if strings.HasPrefix(arg, "-") {
			ln := len(arg)
			// fix: "-abc" -> "-a -b -c"
			if ln > 2 {
				for _, s := range arg[1:] {
					fmtArgs = append(fmtArgs, "-"+string(s))
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

func splitPath2names(path string) []string {
	var names []string
	path = strings.TrimSpace(path)
	if path != "" {
		if strings.ContainsRune(path, ':') { // command ID
			names = strings.Split(path, CommandSep)
		} else if strings.ContainsRune(path, ' ') { // command path
			names = strings.Split(path, " ")
		} else {
			names = []string{path}
		}
	}

	return names
}

// regex: "`[\w ]+`"
// regex: "`.+`"
var codeReg = regexp.MustCompile("`" + `.+` + "`")

// convert "`keywords`" to "<mga>keywords</>"
func wrapColor2string(s string) string {
	if strings.ContainsRune(s, '`') {
		s = codeReg.ReplaceAllStringFunc(s, func(code string) string {
			code = strings.Trim(code, "`")
			return color.WrapTag(code, "mga")
		})
	}
	return s
}
