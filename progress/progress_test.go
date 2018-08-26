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
