package alert

import (
	"fmt"

	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/goutil/x/ccolor"
)

/*
提示消息: Error, Success, Warning, Info, Debug, Notice, Question, Alert, Fatal, Panic

   ╭──────────────────────────────────────────────────────────────────╮
   │   Error : Update available! 3.21.0 → 3.27.0  	 	 	 	 	  │
   ╰──────────────────────────────────────────────────────────────────╯

*/

type MsgBox struct {
	TypeText  string
	TypeColor string
	Content   string
}

var (
	ErrorMsg = MsgBox{
		TypeText:  "ERROR",
		TypeColor: "red1",
	}
)

// Error tips message print
func Error(format string, v ...any) int {
	prefix := ccolor.Red.Sprint("ERROR: ")
	_, _ = fmt.Fprintf(gclicom.Output, prefix+format+"\n", v...)
	return 1
}

// Success tips message print
func Success(format string, v ...any) int {
	prefix := ccolor.Green.Sprint("SUCCESS: ")
	_, _ = fmt.Fprintf(gclicom.Output, prefix+format+"\n", v...)
	return 0
}
