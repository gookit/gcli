package progress

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProgress_Display(t *testing.T) {
	is := assert.New(t)
	ss := widgetMatch.FindAllString(TxtFormat, -1)
	is.Len(ss, 4)

	is.Contains(ss, "{@message}")
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

func ExampleBar() {
	maxStep := 105
	p := CustomBar(60, BarStyles[0], maxStep)
	p.MaxSteps = uint(maxStep)
	p.Format = FullBarFormat

	p.Start()
	for i := 0; i < maxStep; i++ {
		time.Sleep(80 * time.Millisecond)
		p.Advance()
	}
	p.Finish()
}

func ExampleDynamicText() {
	messages := map[int]string{
		// key is percent, range is 0 - 100.
		20:  " Prepare ...",
		40:  " Request ...",
		65:  " Transport ...",
		95:  " Saving ...",
		100: " Handle Complete.",
	}

	maxStep := 105
	p := DynamicText(messages, maxStep)

	p.Start()

	for i := 0; i < maxStep; i++ {
		time.Sleep(80 * time.Millisecond)
		p.Advance()
	}

	p.Finish()
}
