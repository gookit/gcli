package progress

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

var builtinWidgets = map[string]WidgetFunc{
	"elapsed": func(p *Progress) string { // 消耗时间
		// fmt.Sprintf("%.3f", time.Since(startTime).Seconds()*1000)
		sec := time.Since(p.StartedAt()).Seconds()
		return HowLongAgo(int64(sec))
	},
	"remaining": func(p *Progress) string { // 剩余时间
		step := p.Progress() // current progress

		// not set max steps OR current progress is 0
		if p.MaxSteps == 0 || step == 0 {
			return "unknown"
		}

		// calc remaining time
		sec64 := int64(time.Since(p.StartedAt()).Seconds())
		remaining := uint(sec64) / step * (p.MaxSteps - step)
		return HowLongAgo(int64(remaining))
	},
	"estimated": func(p *Progress) string { // 计算总的预计时间
		step := p.Progress() // current progress

		// not set max steps OR current progress is 0
		if p.MaxSteps == 0 || step == 0 {
			return "unknown"
		}

		// calc estimated time
		sec64 := int64(time.Since(p.StartedAt()).Seconds())
		estimated := uint(sec64) / step * p.MaxSteps
		return HowLongAgo(int64(estimated))
	},
	"memory": func(p *Progress) string {
		mem := new(runtime.MemStats)
		runtime.ReadMemStats(mem)
		return formatMemoryVal(mem.Sys)
	},
	"max": func(p *Progress) string {
		return fmt.Sprint(p.MaxSteps)
	},
	"current": func(p *Progress) string {
		step := fmt.Sprint(p.Progress())
		width := fmt.Sprint(p.StepWidth)
		diff := len(width) - len(step)
		if diff <= 0 {
			return step
		}

		return strings.Repeat(" ", diff) + step
	},
	"percent": func(p *Progress) string {
		return fmt.Sprintf("%.1f", p.Percent()*100)
	},
}
