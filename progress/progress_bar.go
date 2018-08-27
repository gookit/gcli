package progress

import "strings"

// internal format for ProgressBar
const (
	DefBarWidth   = 60
	DefBarFormat  = "[{@bar}] {@percent:4s}%({@current}/{@max}){@message}"
	FullBarFormat = "[{@bar}] {@percent:4s}%({@current}/{@max}) {@elapsed:6s}/{@estimated:-6s} {@memory:6s}"
)

var barWidgets = map[string]WidgetFunc{
	"bar": func(p *Progress) string {
		var completeLen float32
		b := p.Binding().(*ProgressBar)

		if p.MaxSteps > 0 { // MaxSteps is valid
			completeLen = p.percent * float32(b.Width)
		} else { // not set MaxSteps
			completeLen = float32(p.step % uint(b.Width))
		}

		bar := strings.Repeat(string(b.Chars.Completed), int(completeLen))

		if diff := int(b.Width) - int(completeLen); diff > 0 {
			ingChar := string(b.Chars.Processing)
			bar += ingChar + strings.Repeat(string(b.Chars.Remaining), diff-len(ingChar))
		}

		return bar
	},
}

// BarChars setting for a progress bar. default {'#', '>', ' '}
type BarChars struct {
	Completed, Processing, Remaining rune
}

// ProgressBar definition.
// Preview:
// 		1 [->--------------------------]
// 		3 [■■■>------------------------]
// 	25/50 [==============>-------------]  50%
//
type ProgressBar struct {
	Progress
	// Width for the bar. default is 100
	Width uint8
	// Chars config for the bar. default {'#', '>', ' '}
	Chars *BarChars
}

// Config the progress instance
func (p *ProgressBar) Config(fn func(p *ProgressBar)) *ProgressBar {
	fn(p)
	return p
}

// Start progress bar
func (p *ProgressBar) Start(maxSteps ...int) {
	if p.Width == 0 {
		p.Width = DefBarWidth
	}

	if p.Chars == nil {
		p.Chars = defaultBarChars()
	}

	p.AddWidgets(barWidgets)
	p.Progress.Start(maxSteps...)
}

// default chars config
func defaultBarChars() *BarChars {
	return &BarChars{'#', '>', ' '}
}

// Bar create new image progress bar.
func Bar(maxSteps ...int) *ProgressBar {
	p := &ProgressBar{
		Progress: *New(maxSteps...),
		// settings for bar
		Width: DefBarWidth,
		Chars: defaultBarChars(),
	}

	p.Format = DefBarFormat
	p.SetBinding(p)
	return p
}

// Tape create new tape progress bar. is alias of Bar()
func Tape(maxSteps ...int) *ProgressBar {
	return Bar(maxSteps...)
}

// FullBar create new progress bar, contains all widgets
func FullBar(maxSteps ...int) *ProgressBar {
	return Bar(maxSteps...).Config(func(p *ProgressBar) {
		p.Format = FullBarFormat
	})
}
