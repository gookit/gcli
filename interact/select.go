package interact

import (
	"bytes"
	"fmt"
	"github.com/gookit/color"
	"sort"
	"strings"
)

// Select definition
type Select struct {
	// Title message for select. e.g "Your city?"
	Title string
	// Options the options data for select. allow: []int,[]string,map[string]string
	Options interface{}
	// DefOpt default option when not input answer
	DefOpt string
	// NoQuit option. if is false, will display "quit" option
	NoQuit bool
	// build from field DefOpt
	defMsg string
}

// NewSelect instance.
// usage:
//	s := NewSelect("Your city?", []string{"chengdu", "beijing"})
//	val := s.Run()
func NewSelect(title string, options interface{}) *Select {
	return &Select{Title: title, Options: options}
}

func (s *Select) prepare() (valArr []string, valMap map[string]string) {
	s.Title = strings.TrimSpace(s.Title)
	if s.Title == "" || s.Options == nil {
		exitWithErr("show.Select: must provide title title and options data")
	}

	switch optsData := s.Options.(type) {
	case map[string]string:
		valMap = optsData
		valArr = make([]string, len(optsData))
		i := 0
		for v := range optsData {
			valArr[i] = v
			i++
		}
	case []int:
	case []int64:
	case []string:
		valArr = make([]string, len(optsData))
		valMap = make(map[string]string, len(optsData))
		for i, n := range optsData {
			v := fmt.Sprint(i)
			valMap[v] = fmt.Sprint(n)
			valArr[i] = v
		}
	default:
		exitWithErr("(show.Select) invalid options data")
	}

	// has default opt
	if s.DefOpt != "" {
		_, has := valMap[s.DefOpt]
		if !has {
			exitWithErr("(show.Select) default option '%s' don't exists", s.DefOpt)
		}

		s.defMsg = fmt.Sprintf("[default:%s]", color.Green.Render(s.DefOpt))
	}

	// sort opt values
	sort.Strings(valArr)
	return
}

// Render select and options to terminal
func (s *Select) render(valArr []string, valMap map[string]string) {
	buf := new(bytes.Buffer)
	green := color.Green.Render

	buf.WriteString(color.Comment.Render(s.Title))
	for _, opt := range valArr {
		buf.WriteString(fmt.Sprintf("\n  %s) %s", green(opt), valMap[opt]))
	}

	if !s.NoQuit {
		valMap["q"] = "quit"
		buf.WriteString(fmt.Sprintf("\n  %s) quit", green("q")))
	}

	// render select and options message to terminal
	color.Println(buf.String())
}

// Run select and receive use input answer
func (s *Select) Run() *Value {
	valArr, valMap := s.prepare()

	// render to console
	s.render(valArr, valMap)

DoSelect:
	ans, err := ReadLine("Your choice" + s.defMsg + ": ")
	if err != nil {
		exitWithErr("(show.Select) %s", err.Error())
	}

	// don't input
	if ans == "" {
		if s.DefOpt != "" { // has default option
			return &Value{s.DefOpt}
		}
		goto DoSelect
	}

	// select error, retry ...
	if _, has := valMap[ans]; !has {
		goto DoSelect
	}

	if !s.NoQuit && ans == "q" { // quit select.
		exitWithMsg(OK, "\n  Quit,ByeBye")
	}

	return &Value{ans}
}
