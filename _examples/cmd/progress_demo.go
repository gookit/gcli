package cmd

import (
	"fmt"
	"time"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/progress"
)

type progressDemo struct {
	maxSteps  int
	overwrite, random bool

	handlers map[string]func(int)
}

var pd = &progressDemo{}

func ProgressDemoCmd() *gcli.Command {
	c := &gcli.Command{
		Name:    "prog",
		Desc:    "there are some progress bar run demos",
		Aliases: []string{"prg-demo", "progress"},
		Func:    pd.Run,
		Config: func(c *gcli.Command) {
			c.IntOpt(&pd.maxSteps, "max-step", "", 100, "setting the max step value")
			c.BoolOpt(&pd.overwrite, "overwrite", "o", true, "setting overwrite progress bar line")
			c.BoolVar(&pd.random, &gcli.FlagMeta{Name: "random", Desc: "use random style for progress bar"})
			// c.AddArg("name",
			// 	"progress bar type name. allow: bar,txt,dtxt,loading,roundTrip",
			// 	true,
			// )
			c.BindArg(gcli.Argument{
				Name: "name",
				Desc: "progress bar type name. allow: bar,txt,dtxt,loading,roundTrip",

				Required: true,
			})
		},
		Examples: `Text progress bar:
  {$fullCmd} txt
Image progress bar:
  {$fullCmd} bar`,
	}

	return c
}

// Run command
func (pd *progressDemo) Run(c *gcli.Command, _ []string) error {
	name := c.Arg("name").String()
	max := pd.maxSteps

	color.Infoln("Progress Demo:")
	switch name {
	case "bar":
		showProgressBar(max)
	case "bars", "all-bar":
		showAllProgressBar(max)
	case "dt", "dtxt", "dynamicText":
		dynamicTextBar(max)
	case "txt", "text":
		txtProgressBar(max)
	case "spr", "load", "loading", "spinner":
		runLoadingBar(max)
	case "rt", "roundTrip":
		runRoundTripBar(max)
	default:
		return c.Errorf("the progress bar type name only allow: bar,txt,dtxt,loading,roundTrip. input is: %s", name)
	}
	return nil
}

func showProgressBar(maxStep int) {
	cs := progress.BarStyles[3]
	if pd.random {
		cs = progress.RandomBarStyle()
	}

	p := progress.CustomBar(40, cs)
	p.MaxSteps = uint(maxStep)
	p.Format = progress.FullBarFormat
	// p.Overwrite = true

	// p.AddMessage("message", " handling ...")

	// running
	runProgressBar(p, maxStep, 60)
	p.Finish()
}

func showAllProgressBar(maxStep int) {
	ln := len(progress.BarStyles)
	ch := make(chan bool, ln)

	for i, style := range progress.BarStyles {
		go func(i int, style progress.BarChars) {
			p := progress.CustomBar(40, style)

			// p.Newline = true
			p.MaxSteps = uint(maxStep)
			// p.Format = progress.FullBarFormat
			p.Format = progress.BarFormat
			p.AddMessage("message", fmt.Sprintf("Bar %d", i+1))

			// run
			runProgressBar(p, maxStep, 100)

			// end
			p.Finish()
			ch <- true
		}(i, style) // NOTICE: must use arguments
	}

	// waiting
	for range progress.BarStyles {
		<-ch
	}

	fmt.Println("- Done with progress, number ", ln)
}

func runRoundTripBar(max int) {
	p := progress.RoundTrip(0).WithMaxSteps(max)

	// running
	runProgressBar(p, max, 120)

	p.Finish()
}

func txtProgressBar(maxStep int) {
	txt := progress.Txt(maxStep)
	txt.AddMessage("message", "Handling ... ")
	// txt.Overwrite = false
	// running
	runProgressBar(txt, maxStep, 80)

	txt.Finish("Completed")
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
