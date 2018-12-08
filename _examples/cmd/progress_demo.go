package cmd

import (
	"github.com/gookit/gcli"
	"github.com/gookit/gcli/progress"
	"time"
)

type progressDemo struct {
	maxSteps  int
	overwrite bool
}

func ProgressDemoCmd() *gcli.Command {
	pd := &progressDemo{}

	return &gcli.Command{
		Name:    "prog",
		UseFor:  "there are some progress bar run demos",
		Aliases: []string{"prg:demo", "progress"},
		Func:    pd.Run,
		Config: func(c *gcli.Command) {
			c.IntOpt(&pd.maxSteps, "max-step", "", 110, "setting the max step value")
			c.BoolOpt(&pd.overwrite, "overwrite", "o", true, "setting overwrite progress bar line")
			c.AddArg("name",
				"progress bar type name. allow: bar,txt,dtxt,loading,roundTrip",
				true,
			)
		},
		Examples: `Text progress bar:
  {$fullCmd} txt
Image progress bar:
  {$fullCmd} bar`,
	}
}

// Run command
func (d *progressDemo) Run(c *gcli.Command, _ []string) error {
	name := c.Arg("name").String()
	max := d.maxSteps

	switch name {
	case "bar":
		imgProgressBar(max)
	case "dt", "dtxt", "dynamicText":
		dynamicTextBar(max)
	case "txt", "text":
		txtProgressBar(max)
	case "load", "loading", "spinner":
		runLoadingBar(max)
	case "rt", "roundTrip":
		runRoundTripBar(max)
	default:
		return c.Errorf("the progress bar type name only allow: bar,txt,dtxt,loading,roundTrip. input is: %s", name)
	}
	return nil
}

func runRoundTripBar(max int) {
	p := progress.RoundTrip(0).WithMaxSteps(max)

	// running
	runProgressBar(p, max, 120)

	p.Finish()
}

func txtProgressBar(maxStep int) {
	txt := progress.Txt(maxStep)
	txt.AddMessage("message", "handling ... ")
	// txt.Overwrite = false
	// running
	runProgressBar(txt, maxStep, 80)

	txt.Finish()
}

func dynamicTextBar(maxStep int) {
	messages := map[int]string{
		// key is percent, range is 0 - 100.
		20:  " Prepare ...",
		40:  " Request ...",
		65:  " Transport ...",
		95:  " Saving ...",
		100: " Handle Complete.",
	}

	// maxStep = 10
	p := progress.DynamicText(messages, maxStep)
	// p.Overwrite = false

	// running
	runProgressBar(p, maxStep, 100)
	p.Finish()
}

func imgProgressBar(maxStep int) {
	cs := progress.RandomBarStyle()

	p := progress.CustomBar(60, cs)
	p.MaxSteps = uint(maxStep)
	p.Format = progress.FullBarFormat
	// p.Overwrite = false

	// p.AddMessage("message", " handling ...")

	// running
	runProgressBar(p, maxStep, 100)
	p.Finish()
}

func runLoadingBar(maxStep int) {
	p := progress.LoadingBar(progress.RandomCharsTheme())
	p.MaxSteps = uint(maxStep)
	p.AddMessage("message", " data loading ... ...")

	// running
	runProgressBar(p, maxStep, 70)

	// p.Finish()
	p.Finish("data load complete")
}

// running
func runProgressBar(p *progress.Progress, maxSteps int, speed int) {
	p.Start()
	for i := 0; i < maxSteps; i++ {
		time.Sleep(time.Duration(speed) * time.Millisecond)
		p.Advance()
	}
}
