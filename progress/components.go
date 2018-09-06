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
	return New(maxSteps...).Config(func(p *Progress) {
		p.Format = FullFormat
	})
}

// Counter progress bar create
func Counter(maxSteps ...int) *Progress {
	return New(maxSteps...).Config(func(p *Progress) {
		p.Format = MinFormat
	})
}

// DynamicText progress bar create
func DynamicText(messages map[int]string) *Progress {
	return New().AddWidget("message", DynamicTextWidget(messages))
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

// some built in BarChars style
var (
	BarCharsStyle  = BarChars{'#', '>', ' '}
	BarCharsStyle1 = BarChars{'▉', '▉', '░'}
	BarCharsStyle2 = BarChars{'■', '■', ' '}
	BarCharsStyle3 = BarChars{'■', '▶', ' '}
	BarCharsStyle4 = BarChars{'=', '>', ' '}
)

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

// Config the progress instance
func (pb ProgressBar) Create(maxSteps ...int) *Progress {
	return New(maxSteps...).Config(func(p *Progress) {
		p.Format = DefBarFormat
	}).AddWidget("bar", ProgressBarWidget(pb.Width, pb.Chars))
}

// Bar create a default progress bar.
func Bar(maxSteps ...int) *Progress {
	return CustomBar(DefBarWidth, BarCharsStyle).WithMaxSteps(maxSteps...)
}

// Tape create new tape progress bar. is alias of Bar()
func Tape(maxSteps ...int) *Progress {
	return Bar(maxSteps...)
}

// CustomBar create a custom progress bar.
func CustomBar(width int, cs BarChars) *Progress {
	return New().Config(func(p *Progress) {
		p.Format = DefBarFormat
	}).AddWidget("bar", ProgressBarWidget(width, cs))
}

/*************************************************************
 * RoundTrip progress bar: `[ ====   ] Pending ...`
 *************************************************************/

// RoundTripBar config
type RoundTripBar struct {
	Char     rune
	CharNum  int
	BoxWidth int
}

// Create Progress bar from RoundTripBar config.
func (rt RoundTripBar) Create(maxSteps ...int) *Progress {
	return RoundTrip(rt.Char, rt.CharNum, rt.BoxWidth).WithMaxSteps(maxSteps...)
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

// CharsThemes collection
var CharsThemes = [][]rune{
	{'卍', '卐'},
	{'◐', '◒', '◓', '◑'},
	{'✣', '✤', '✥', '❉'},
	{'-', '\\', '|', '/'},
	[]rune("▖▘▝▗"),
	[]rune("◢◣◤◥"),
	[]rune("⌞⌟⌝⌜"),
	[]rune("◎●◯◌○⊙"),
	[]rune("◡◡⊙⊙◠◠"),
	[]rune("←↖↑↗→↘↓↙"),
	[]rune("㊎㊍㊌㊋㊏"),
	[]rune("⣾⣽⣻⢿⡿⣟⣯⣷"),
	[]rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏"),
	[]rune("▉▊▋▌▍▎▏▎▍▌▋▊▉"),
	[]rune("🌍🌎🌏"),
	[]rune("⠋⠙⠚⠒⠂⠂⠒⠲⠴⠦⠖⠒⠐⠐⠒⠓⠋"),
	[]rune("🕐🕑🕒🕓🕔🕕🕖🕗🕘🕙🕚🕛"),
}

// LoadingBar alias of load bar LoadBar()
func LoadingBar(chars []rune) *Progress {
	return LoadBar(chars)
}

// LoadBar create a loading progress bar
func LoadBar(chars []rune) *Progress {
	return New().Config(func(p *Progress) {
		p.Format = "{@loading} {@message}"
		p.AddWidget("loading", LoadingWidget(chars))
	})
}

// RandomCharsTheme get
func RandomCharsTheme() []rune {
	rand.Seed(time.Now().UnixNano())
	return CharsThemes[rand.Intn(len(CharsThemes)-1)]
}
