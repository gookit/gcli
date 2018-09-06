package progress

import (
	"fmt"
	"github.com/gookit/color"
	"strings"
	"sync"
	"time"
)

// Spinner definition
type Spinner struct {
	// Delay is the running speed
	Delay  time.Duration
	Format string
	chars  []rune
	active bool
	lock   *sync.RWMutex
	// control the spinner running.
	stopCh chan struct{}
}

// SpinnerBar instance create
func SpinnerBar(cs []rune, speed time.Duration) *Spinner {
	return NewSpinner(cs, speed)
}

// NewSpinner instance
func NewSpinner(cs []rune, speed time.Duration) *Spinner {
	return &Spinner{
		Delay:  speed,
		Format: "%s",
		// color: color.Normal.Sprint,
		lock:  &sync.RWMutex{},
		chars: cs,
		// writer:   os.Stdout,
		stopCh: make(chan struct{}, 1),
	}
}

func (s *Spinner) prepare(format []string) {
	if len(format) > 0 {
		s.Format = format[0]
	}

	if s.Format != "" && !strings.Contains(s.Format, "%s") {
		s.Format = "%s " + s.Format
	}

	if len(s.chars) == 0 {
		s.chars = RandomCharsTheme()
	}

	if s.Delay == 0 {
		s.Delay = 100 * time.Millisecond
	}
}

// Start run spinner
func (s *Spinner) Start(format ...string) {
	if s.active {
		return
	}

	s.active = true
	s.prepare(format)

	go func() {
		index := 0
		length := len(s.chars)

		for {
			select {
			case <-s.stopCh:
				return
			default:
				s.lock.Lock()
				char := string(s.chars[index])
				if index+1 == length { // reset
					index = 0
				} else {
					index++
				}

				// \x0D - Move the cursor to the beginning of the line
				// \x1B[2K - Erase(Delete) the line
				fmt.Print("\x0D\x1B[2K")
				color.Printf(s.Format, char)
				s.lock.Unlock()

				time.Sleep(s.Delay)
			}
		}
	}()
}

// Stop run spinner
func (s *Spinner) Stop(finalMsg ...string) {
	if !s.active {
		return
	}

	s.lock.Lock()
	s.active = false
	fmt.Print("\x0D\x1B[2K")

	if len(finalMsg) > 0 {
		fmt.Println(finalMsg[0])
	}

	s.stopCh <- struct{}{}
	s.lock.Unlock()
}

// Restart will stop and start the spinner
func (s *Spinner) Restart() {
	s.Stop()
	s.Start()
}

// Active status
func (s *Spinner) Active() bool {
	return s.active
}
