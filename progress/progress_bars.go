package progress

import "strings"

// some built in chars
const (
	CharStar   rune = '*'
	CharPlus   rune = '+'
	CharWell   rune = '#'
	CharEqual  rune = '='
	CharSpace  rune = ' '
	CharSquare rune = '■'
	// Hyphen Minus
	CharHyphen     rune = '-'
	CharCNHyphen   rune = '—'
	CharUnderline  rune = '_'
	CharLeftArrow  rune = '<'
	CharRightArrow rune = '>'
)

// Counter progress bar create
func Counter() *Progress {
	return New().Config(func(p *Progress) {
		p.Format = MinFormat
	})
}

// RoundTripBar create a RoundTrip progress bar.
// Usage:
// 	p := RoundTrip(CharEqual, 100)
// 	// p := RoundTrip('*', 100) // custom char
// 	p.Start()
// 	....
// 	p.Finish()
func RoundTrip(char rune, charNumAndBoxWidth ...int) *Progress {
	if char == 0 {
		char = CharEqual
	}

	charNum := 4
	boxWidth := 12

	if ln := len(charNumAndBoxWidth); ln > 0 {
		charNum = charNumAndBoxWidth[0]
		if ln > 1 {
			boxWidth = charNumAndBoxWidth[1]
		}
	}

	cursor := string(repeatRune(char, charNum))

	// direction: <- OR ->
	left := false
	// record cursor position
	position := 0

	return New().Config(func(p *Progress) {
		p.Format = "[{@bar}] {@percent:4s}% ({@current}/{@max}){@message}"
	}).AddWidget("bar", func(p *Progress) string {
		var bar string
		if position > 0 {
			bar += strings.Repeat(" ", position)
		}

		bar += cursor + strings.Repeat(" ", boxWidth-position-charNum)

		if left { // left <-
			if position <= 0 { // begin ->
				left = false
			} else {
				position--
			}
		} else { // -> right
			if position+charNum >= boxWidth { // begin <-
				left = true
			} else {
				position++
			}
		}

		return bar
	})
}

/*************************************************************
 * loading bar
 *************************************************************/

// default spinner chars: -\|/
var (
	LoadingTheme1 = []rune{'-', '\\', '|', '/'}
	LoadingTheme2 = []rune{'◐', '◒', '◓', '◑'}
	LoadingTheme3 = []rune{'✣', '✤', '✥', '❉'}
	LoadingTheme4 = []rune{'卍', '卐'}
)

// LoadingBar create a loading progress bar
func LoadingBar(chars []rune) *Progress {
	// chars := []rune(`-\|/`)
	if len(chars) == 0 {
		chars = LoadingTheme1
	}

	length := len(chars)
	counter := 0

	return New().Config(func(p *Progress) {
		p.Format = "{@loading} {@message}"
	}).
		AddWidget("loading", func(p *Progress) string {
			char := string(chars[counter])
			if counter+1 == length {
				counter = 0 // reset
			} else {
				counter++
			}

			return char
		})
}
