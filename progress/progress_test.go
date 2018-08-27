package progress

import (
	"fmt"
	"strings"
	"testing"
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
