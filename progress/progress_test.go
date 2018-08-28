package progress

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestProgress_Display(t *testing.T) {
	ss := widgetMatch.FindAllString(TxtFormat, -1)
	fmt.Println(ss)

	widgetMatch.ReplaceAllStringFunc(TxtFormat, func(s string) string {
		fmt.Println(s, strings.Trim(s, "{@}"))
		return s
	})
}

func TestSpinner(t *testing.T) {
	chars := []rune(`你\|/`)
	str := `你\|/`

	fmt.Println(chars, string(chars[0]), string(str[0]))
}

func TestLoading(t *testing.T) {
	chars := []rune("◐◑◒◓")
	str := "◐◑◒◓"

	fmt.Println(chars, string(chars[0]), str, string(str[0]))
}

func ExampleDynamicTextWidget() {
	widget := DynamicTextWidget(map[int]string{
		// int is percent, range is 0 - 100.
		20:  " Prepare ...",
		40:  " Request ...",
		65:  " Transport ...",
		95:  " Saving ...",
		100: " Handle Complete.",
	})

	maxStep := 105
	p := New(maxStep).Config(func(p *Progress) {
		p.Format = "{@percent:4s}%({@current}/{@max}) {@message}"
	}).AddWidget("message", widget)

	runProgressBar(p, maxStep, 80)

	p.Finish()
}

// running
func runProgressBar(p *Progress, maxStep int, speed int) {
	p.Start()
	for i := 0; i < maxStep; i++ {
		time.Sleep(time.Duration(speed) * time.Millisecond)
		p.Advance()
	}
}
