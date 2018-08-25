package progress

// Progress definition
type Progress struct {
	// current step value
	step int
	// default is 1
	stepSize uint8
	// completed percent. eg: "83.8"
	percent float32
	// RedrawFreq redraw freq. default is 1
	RedrawFreq uint8
	// Overwrite default is True
	Overwrite bool
}

// New Progress instance
func New() *Progress {
	return &Progress{}
}

// Start the progress bar
func (p *Progress) Start() {

}

// Advance one step
func (p *Progress) Advance() {

}

// AdvanceTo a special step number
func (p *Progress) AdvanceTo(steps uint) {

}

func Text() {

}

const (
	DefFormat  = "[{@bar}] {@percent:3s}%({@current}/{@max})"
	FullFormat = "[{@bar}] {@percent:3s}%({@current}/{@max}) {@elapsed:6s}/{@estimated:-6s} {@memory:6s}"
)

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
	// MaxSteps maximal steps.
	MaxSteps uint
	// Chars config for the bar. default {'#', '>', ' '}
	Chars struct {
		Completed, Processing, Remaining rune
	}
}

func Bar() {

}
