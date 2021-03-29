package gcli

import (
	"fmt"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
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

func defaultErrHandler(data ...interface{}) (stop bool) {
	if len(data) == 2 && data[1] != nil {
		if err, ok := data[1].(error); ok {
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
	return DefaultVerb
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

// parse tag value. eg: "name=int0;shorts=i;required=true;desc=int option message"
func parseTagValue(name, str string) (mp map[string]string) {
	ss := strutil.Split(str, ";")
	if len(ss) == 0 {
		return
	}

	mp = make(map[string]string, len(flagTagKeys))
	for _, s := range ss {
		if strings.ContainsRune(s, '=') == false {
			panicf("parse tag error on field '%s': item must match `KEY=VAL`", name)
		}

		kvNodes := strings.SplitN(s, "=", 2)
		key, val := kvNodes[0], kvNodes[1]
		if !flagTagKeys.Has(key) {
			panicf("parse tag error on field '%s': invalid key name '%s'", name, key)
		}

		mp[key] = val
	}
	return
}

// func exitWithMsg(format string, v ...interface{}) {
// 	fmt.Printf(format, v...)
// 	Exit(0)
// }

const (
	// match an good option, argument name
	regGoodName = `^[a-zA-Z][\w-]*$`
	// match an good command name
	regGoodCmdName = `^[a-zA-Z][\w-]*$`
	// match command id. eg: "self:init"
	regGoodCmdId = `^[a-zA-Z][\w:-]*$`
	// match command path. eg: "self init"
	// regGoodCmdPath = `^[a-zA-Z][\w -]*$`
)

var (
	// good name for option and argument
	goodName = regexp.MustCompile(regGoodName)
	// match an good command name
	goodCmdId = regexp.MustCompile(regGoodCmdId)
	// match an good command name
	goodCmdName = regexp.MustCompile(regGoodCmdName)
)

func isValidCmdName(name string) bool {
	if name[0] == '-' { // is option name.
		return false
	}

	return goodCmdName.MatchString(name)
}

func isValidCmdId(name string) bool {
	if name[0] == '-' { // is option name.
		return false
	}

	return goodCmdId.MatchString(name)
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

// split "ef" to ["e", "f"]
// split "e,f" to ["e", "f"]
func splitShortcut(str string) (ss []string) {
	bs := []byte(str)
	if len(bs) == 0 {
		return
	}

	for _, b := range bs {
		// if b == ',' {
		// 	continue
		// }
		if strutil.IsAlphabet(b) {
			ss = append(ss, string(b))
		}
	}
	return
}

func shorts2string(ss []string) string {
	var newSs []string
	for _, s := range ss {
		newSs = append(newSs, "-"+s)
	}

	// eg: "-t, -o"
	return strings.Join(newSs, ", ")
}

// regex: "`[\w ]+`"
var codeReg = regexp.MustCompile("`" + `[\w ]+` + "`")

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
