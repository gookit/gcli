package interact

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/gookit/color"
)

// Select definition
type Select struct {
	// Title message for select. e.g "Your city?"
	Title string
	// Options the options data for select. allow: []int,[]string,map[string]string
	Options interface{}
	// DefOpt default option when not input answer
	DefOpt string
	// DefOpts use for `MultiSelect` is true
	DefOpts []string
	// DisableQuit option. if is false, will display "quit" option. default False
	DisableQuit bool
	// MultiSelect allow multi select. default False
	MultiSelect bool
	// QuitHandler func()
	// parsed options data
	// {
	// 	"option value": "option name",
	// }
	valMap map[string]string
}

// NewSelect instance.
// usage:
// 	s := NewSelect("Your city?", []string{"chengdu", "beijing"})
// 	val := s.Run().String() // "1"
func NewSelect(title string, options interface{}) *Select {
	return &Select{
		Title:   title,
		Options: options,
	}
}

func (s *Select) prepare() (valArr []string) {
	s.Title = strings.TrimSpace(s.Title)
	if s.Title == "" || s.Options == nil {
		exitWithErr("(interact.Select) must provide title and options data")
	}

	s.valMap = make(map[string]string)
	handleArrItem := func(i int, v interface{}) {
		nv := fmt.Sprint(i)
		s.valMap[nv] = fmt.Sprint(v)
		valArr = append(valArr, nv)
	}

	switch optsData := s.Options.(type) {
	case map[string]int:
		valArr = make([]string, len(optsData))
		i := 0
		for v, n := range optsData {
			valArr[i] = v
			s.valMap[v] = fmt.Sprint(n)
			i++
		}

		sort.Strings(valArr) // sort
	case map[string]string:
		s.valMap = optsData
		valArr = make([]string, len(optsData))
		i := 0
		for v := range optsData {
			valArr[i] = v
			i++
		}

		sort.Strings(valArr) // sort
	case string:
		ss := stringToArr(optsData, ",")
		for i, v := range ss {
			handleArrItem(i, v)
		}
	case []int:
		for i, v := range optsData {
			handleArrItem(i, v)
		}
	case []string:
		for i, v := range optsData {
			handleArrItem(i, v)
		}
	default:
		exitWithErr("(interact.Select) invalid options data")
	}

	// format some field data
	s.DefOpt = strings.TrimSpace(s.DefOpt)
	if len(s.DefOpts) > 0 {
		var ss []string
		for _, v := range s.DefOpts {
			if v = strings.TrimSpace(v); v != "" {
				ss = append(ss, v)
			}
		}

		s.DefOpts = ss
	}

	return
}

// Render select and options to terminal
func (s *Select) render(valArr []string) {
	buf := new(bytes.Buffer)
	green := color.Green.Render

	buf.WriteString(color.Comment.Render(s.Title))
	for _, opt := range valArr {
		buf.WriteString(fmt.Sprintf("\n  %s) %s", green(opt), s.valMap[opt]))
	}

	if !s.DisableQuit {
		s.valMap["q"] = "quit"
		buf.WriteString(fmt.Sprintf("\n  %s) quit", green("q")))
	}

	// render select and options message to terminal
	color.Println(buf.String())
	buf = nil
}

func (s *Select) selectOne() *Value {
	tipsText := "Your choice: "

	// has default opt
	if s.DefOpt != "" {
		if _, has := s.valMap[s.DefOpt]; !has {
			exitWithErr("(interact.Select) default option '%s' don't exists", s.DefOpt)
		}

		defMsg := fmt.Sprintf("[default:%s]", color.Green.Render(s.DefOpt))
		tipsText = "Your choice" + defMsg + ": "
	}

DoSelect:
	ans, err := ReadLine(tipsText)
	if err != nil {
		exitWithErr("(interact.Select) %s", err.Error())
	}

	if ans == "" { // empty input
		if s.DefOpt != "" { // has default option
			return &Value{s.DefOpt}
		}

		goto DoSelect // retry ...
	}

	// check input
	if _, has := s.valMap[ans]; !has {
		color.Error.Println("Unknown option value: ", ans)
		goto DoSelect // retry ...
	}

	// quit select.
	if !s.DisableQuit && ans == "q" {
		exitWithMsg(OK, "\n  Quit,ByeBye")
	}

	return &Value{ans}
}

// for enable MultiSelect
func (s *Select) selectMulti() *Value {
	hasDefault := len(s.DefOpts) > 0
	tipsText := "Your choice(multi use <magenta>,</> separate): "
	if hasDefault {
		tipsText = fmt.Sprintf(
			"Your choice(multi use <magenta>,</> separate)[default:%s]: ",
			color.Green.Render(strings.Join(s.DefOpts, ",")),
		)
	}

DoSelect:
	ans, err := ReadLine(tipsText)
	if err != nil {
		exitWithErr("(interact.Select) %s", err.Error())
	}

	values := stringToArr(ans, ",")
	if len(values) == 0 { // empty input
		// has default options
		if hasDefault {
			return &Value{s.DefOpts}
		}

		goto DoSelect // retry ...
	}

	// check input
	for _, v := range values {
		if _, has := s.valMap[v]; !has {
			color.Error.Println("Unknown option value: ", v)
			goto DoSelect // retry ...
		}

		// quit select.
		if !s.DisableQuit && v == "q" {
			exitWithMsg(OK, "\n  Quit,ByeBye")
		}
	}

	return &Value{values}
}

// Run select and receive use input answer
func (s *Select) Run() *Value {
	valArr := s.prepare()
	// render to console
	s.render(valArr)

	// if enable MultiSelect
	if s.MultiSelect {
		return s.selectMulti()
	}

	return s.selectOne()

}
