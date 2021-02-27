package gcli

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
)

/*************************************************************
 * console log
 *************************************************************/

var level2name = map[VerbLevel]string{
	VerbError: "ERROR",
	VerbWarn:  "WARN",
	VerbInfo:  "INFO",
	VerbDebug: "DEBUG",
	VerbCrazy: "CRAZY",
}

var level2color = map[VerbLevel]color.Color{
	VerbError: color.FgRed,
	VerbWarn:  color.FgYellow,
	VerbInfo:  color.FgGreen,
	VerbDebug: color.FgCyan,
	VerbCrazy: color.FgMagenta,
}

// Logf print log message
func Debugf(format string, v ...interface{}) {
	logf(VerbDebug, format, v...)
}

// Logf print log message
func Logf(level VerbLevel, format string, v ...interface{}) {
	logf(level, format, v...)
}

// print log message
func logf(level VerbLevel, format string, v ...interface{}) {
	if gOpts.verbose < level {
		return
	}

	var fnName string
	pc, fName, line, ok := runtime.Caller(2)
	if !ok {
		fnName, fName, line = "UNKNOWN", "???.go", 0
	} else {
		fName = path.Base(fName)
		fnName = runtime.FuncForPC(pc).Name()
	}

	name := level.Upper()
	name = level2color[level].Render(name)
	color.Printf("GCli: [%s] [%s(), %s:%d] %s\n", name, fnName, fName, line, fmt.Sprintf(format, v...))
}

func defaultErrHandler(data ...interface{}) {
	if len(data) == 2 && data[1] != nil {
		if err, ok := data[1].(error); ok {
			color.Error.Tips(err.Error())
			// fmt.Println(color.Red.Render("ERROR:"), err.Error())
		}
	}
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
	return VerbError
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

func panicf(format string, v ...interface{}) {
	panic(fmt.Sprintf("GCli: "+format, v...))
}

// func exitWithMsg(format string, v ...interface{}) {
// 	fmt.Printf(format, v...)
// 	Exit(0)
// }

func isValidCmdName(name string) bool {
	if name[0] == '-' { // is option name.
		return false
	}

	return goodCmdName.MatchString(name)
}

func aliasNameCheck(name string) {
	if goodCmdName.MatchString(name) {
		return
	}

	panicf("alias name '%s' is invalid, must match: %s", name, regGoodCmdName)
}

// strictFormatArgs
// TODO mode:
//  POSIX '-ab' will split to '-a -b', '--o' -> '-o'
//  UNIX '-ab' will split to '-a b'
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
			// maybe ...

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

// split "ef" to ["e", "f"]
func splitShortStr(str string) (ss []string) {
	bs := []byte(str)

	for _, b := range bs {
		if strutil.IsAlphabet(b) {
			ss = append(ss, string(b))
		}
	}
	return
}

func shorts2str(ss []string) string {
	var newSs []string
	for _, s := range ss {
		newSs = append(newSs, "-"+s)
	}

	// eg: "-t, -o"
	return strings.Join(newSs, ", ")
}
