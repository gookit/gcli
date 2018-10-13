package progress

import (
	"math/rand"
	"time"
)

// some built in chars
const (
	CharStar    rune = '*'
	CharPlus    rune = '+'
	CharWell    rune = '#'
	CharEqual   rune = '='
	CharEqual1  rune = '═'
	CharSpace   rune = ' '
	CharCenter  rune = '●'
	CharSquare  rune = '■'
	CharSquare1 rune = '▇'
	CharSquare2 rune = '▉'
	CharSquare3 rune = '░'
	CharSquare4 rune = '▒'
	CharSquare5 rune = '▢'
	// Hyphen Minus
	CharHyphen      rune = '-'
	CharCNHyphen    rune = '—'
	CharUnderline   rune = '_'
	CharLeftArrow   rune = '<'
	CharRightArrow  rune = '>'
	CharRightArrow1 rune = '▶'
)

// Txt progress bar create.
func Txt(maxSteps ...int) *Progress {
	return New(maxSteps...).Config(func(p *Progress) {
		p.Format = TxtFormat
	})
}

// Full text progress bar create.
func Full(maxSteps ...int) *Progress {
	return NewWithConfig(func(p *Progress) {
		p.Format = FullFormat
	}, maxSteps...)
}

// Counter progress bar create
func Counter(maxSteps ...int) *Progress {
	return NewWithConfig(func(p *Progress) {
		p.Format = MinFormat
	}, maxSteps...)
}

// DynamicText progress bar create
func DynamicText(messages map[int]string, maxSteps ...int) *Progress {
	return NewWithConfig(func(p *Progress) {
		p.Format = "{@percent:4s}%({@current}/{@max}){@message}"
		p.AddWidget("message", DynamicTextWidget(messages))
	}, maxSteps...)
}

/*************************************************************
 * Generic progress bar
 *************************************************************/

// internal format for ProgressBar
const (
	DefBarWidth   = 60
	DefBarFormat  = "{@bar} {@percent:4s}%({@current}/{@max}){@message}"
	FullBarFormat = "{@bar} {@percent:4s}%({@current}/{@max}) {@elapsed:7s}/{@estimated:-7s} {@memory:6s}"
)

// BarChars setting for a progress bar. default {'#', '>', ' '}
type BarChars struct {
	Completed, Processing, Remaining rune
}

// BarStyles some built in BarChars style collection
var BarStyles = []BarChars{
	{'=', '>', ' '},
	{'#', '>', ' '},
	{'-', '>', '-'},
	{'▉', '▉', '░'},
	{'■', '■', ' '},
	{'■', '■', '▢'},
	{'■', '▶', ' '},
}

// ProgressBar definition.
// Preview:
// 		1 [->--------------------------]
// 		3 [■■■>------------------------]
// 	25/50 [==============>-------------]  50%
//
type ProgressBar struct {
	// Width for the bar. default is 100
	Width int
	// Chars config for the bar. default {'#', '>', ' '}
	Chars BarChars
}

// Create the progress instance
func (pb ProgressBar) Create(maxSteps ...int) *Progress {
	return New(maxSteps...).Config(func(p *Progress) {
		p.Format = DefBarFormat
	}).AddWidget("bar", ProgressBarWidget(pb.Width, pb.Chars))
}

// RandomBarStyle get
func RandomBarStyle() BarChars {
	rand.Seed(time.Now().UnixNano())
	return BarStyles[rand.Intn(len(BarStyles)-1)]
}

// Bar create a default progress bar.
func Bar(maxSteps ...int) *Progress {
	return CustomBar(DefBarWidth, BarStyles[0], maxSteps...)
}

// Tape create new tape progress bar. is alias of Bar()
func Tape(maxSteps ...int) *Progress {
	return Bar(maxSteps...)
}

// CustomBar create a custom progress bar.
func CustomBar(width int, cs BarChars, maxSteps ...int) *Progress {
	return New(maxSteps...).Config(func(p *Progress) {
		p.Format = DefBarFormat
	}).AddWidget("bar", ProgressBarWidget(width, cs))
}

/*************************************************************
 * RoundTrip progress bar: `[ ====   ] Pending ...`
 *************************************************************/

// RoundTripBar alias of RoundTrip()
func RoundTripBar(char rune, charNumAndBoxWidth ...int) *Progress {
	return RoundTrip(char, charNumAndBoxWidth...)
}

// RoundTrip create a RoundTrip progress bar.
// Usage:
// 	p := RoundTrip(CharEqual)
// 	// p := RoundTrip('*') // custom char
// 	p.Start()
// 	....
// 	p.Finish()
func RoundTrip(char rune, charNumAndBoxWidth ...int) *Progress {
	charNum := 4
	boxWidth := 12
	if ln := len(charNumAndBoxWidth); ln > 0 {
		charNum = charNumAndBoxWidth[0]
		if ln > 1 {
			boxWidth = charNumAndBoxWidth[1]
		}
	}

	return New().
		AddWidget("rtBar", RoundTripWidget(char, charNum, boxWidth)).
		Config(func(p *Progress) {
			p.Format = "[{@rtBar}] {@percent:4s}% ({@current}/{@max}){@message}"
		})
}

/*************************************************************
 * Loading bar
 *************************************************************/

// LoadingBar alias of load bar LoadBar()
func LoadingBar(chars []rune, maxSteps ...int) *Progress {
	return LoadBar(chars, maxSteps...)
}

// SpinnerBar alias of load bar LoadBar()
func SpinnerBar(chars []rune, maxSteps ...int) *Progress {
	return LoadBar(chars, maxSteps...)
}

// LoadBar create a loading progress bar
func LoadBar(chars []rune, maxSteps ...int) *Progress {
	return New(maxSteps...).Config(func(p *Progress) {
		p.Format = "{@loading} ({@current}/{@max}){@message}"
		p.AddWidget("loading", LoadingWidget(chars))
	})
}
