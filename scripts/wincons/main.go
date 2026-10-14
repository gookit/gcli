//go:build windows

// Command wincons drives the example app (./_examples/cliapp) on a real Windows
// console: the Windows twin of scripts/ptysmoke (PTY, Linux/macOS only).
//
// Python's pty module does not exist on Windows, and a ConPTY child cannot be
// started from a gofer job (CreateProcess with PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE
// fails with 0xC0000142 there), so this harness drives the console Windows has:
// it detaches from the console it was started with, creates a private one
// (AllocConsole), starts the app with that console's CONIN$ as stdin and its
// stdout/stderr on a pipe, then types the answer as KEY_EVENT records with
// WriteConsoleInput - exactly what a keyboard writes into the console input
// buffer. The app therefore sees a genuine console on stdin (cooked line input,
// the prompt is real), while the harness can still read what it printed.
//
//	go build -o tmp/wincons/cliapp.exe ./_examples/cliapp
//	go run ./scripts/wincons --bin tmp/wincons/cliapp.exe
//
// Every case is independent: it gets a fresh console, so no stale key event can
// leak into the next one. Exit status 0 when every check passes.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32                        = syscall.NewLazyDLL("kernel32.dll")
	user32                          = syscall.NewLazyDLL("user32.dll")
	procFreeConsole                 = kernel32.NewProc("FreeConsole")
	procAllocConsole                = kernel32.NewProc("AllocConsole")
	procWriteConsoleInputW          = kernel32.NewProc("WriteConsoleInputW")
	procGetConsoleMode              = kernel32.NewProc("GetConsoleMode")
	procGetNumberOfConsoleInputEvts = kernel32.NewProc("GetNumberOfConsoleInputEvents")
	procVkKeyScanW                  = user32.NewProc("VkKeyScanW")
)

const (
	keyEventType = 0x0001
	vkReturn     = 0x0D

	// prompts printed by _examples/cmd (question_demo.go / shared_opt_demo.go)
	askLocal  = "Please input your access token:"
	askShared = "Please input your api token:"

	promptWait = 15 * time.Second
	exitWait   = 15 * time.Second
)

// keyEventRecord is the KEY_EVENT_RECORD of a console input record.
type keyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

// inputRecord is one INPUT_RECORD written with WriteConsoleInputW.
type inputRecord struct {
	EventType uint16
	_         uint16
	Event     [16]byte
}

// step is one run of the app plus the checks on its output.
type step struct {
	name string
	args []string
	// ask is the prompt the case is about; it is waited for before answer is typed.
	ask string
	// answer is typed on the console once the prompt showed up.
	answer string
	// wantAsk: the prompt must show up exactly once (false: it must not show up).
	wantAsk bool
	want    []string
	notWant []string
}

type checker struct {
	results []string
	failed  int
}

func (c *checker) check(name string, ok bool) {
	if ok {
		fmt.Println("PASS " + name)
		c.results = append(c.results, name)
		return
	}
	fmt.Println("FAIL " + name)
	c.failed++
	c.results = append(c.results, name+" (FAIL)")
}

func main() {
	bin := flag.String("bin", filepath.Join("tmp", "wincons", "cliapp.exe"), "cliapp binary (built from ./_examples/cliapp)")
	verbose := flag.Bool("verbose", false, "print the app output and console state of every case")
	flag.Parse()

	abs, err := filepath.Abs(*bin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "binary path:", err)
		os.Exit(2)
	}
	if _, err := os.Stat(abs); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nbuild it first: go build -o %s ./_examples/cliapp\n", err, *bin)
		os.Exit(2)
	}

	// 1. empty value: the console hands over a typed line
	// 2. value on the command line: no question on the console
	// 3. shared option, empty value: asked once, at the executed command
	// 4. shared option, value in the leaf segment: no question either
	// 5. shared option, -h: help never asks
	steps := []step{
		{
			name: "ask-demo: empty value", args: []string{"ask-demo"},
			ask: askLocal, answer: "win-token", wantAsk: true,
			want: []string{"win-token"},
		},
		{
			name: "ask-demo: value on the command line", args: []string{"ask-demo", "--token=cli-token"},
			ask: askLocal, want: []string{"cli-token"},
		},
		{
			name: "shared-ask run: empty value", args: []string{"shared-ask", "run"},
			ask: askShared, answer: "win-shared-token", wantAsk: true,
			want: []string{"leaf: token=win-shared-token"},
		},
		{
			name: "shared-ask run: value in the leaf segment", args: []string{"shared-ask", "run", "--token=leaf-token"},
			ask: askShared, want: []string{"leaf: token=leaf-token"},
		},
		{
			name: "shared-ask run -h: help never asks", args: []string{"shared-ask", "run", "-h"},
			ask: askShared, want: []string{"--token"},
		},
	}

	ck := &checker{}
	fmt.Printf("binary: %s\n", abs)
	fmt.Println("note: a console window opens for each case; the app's output is read through a pipe")
	for _, st := range steps {
		runStep(ck, abs, st, *verbose)
	}

	fmt.Printf("\n%d/%d checks passed\n", len(ck.results)-ck.failed, len(ck.results))
	if ck.failed > 0 {
		os.Exit(1)
	}
}

// runStep drives one case on a fresh private console.
func runStep(ck *checker, bin string, st step, verbose bool) {
	if st.ask == "" {
		panic("step " + st.name + ": ask pattern is required")
	}

	in, childIn, err := newConsole()
	if err != nil {
		ck.check(st.name+": console", false)
		fmt.Printf("  %v\n", err)
		return
	}
	defer func() {
		_ = in.Close()
		_ = childIn.Close()
		_, _, _ = procFreeConsole.Call() // release it before the next case
	}()

	out := &lockedBuffer{}
	cmd := exec.Command(bin, st.args...)
	cmd.Stdin = childIn
	cmd.Stdout, cmd.Stderr = out, out
	cmd.Env = append(cleanEnv(os.Environ()), "TERM=xterm", "NO_COLOR=1")
	if err := cmd.Start(); err != nil {
		ck.check(st.name+": start", false)
		fmt.Printf("  %v\n", err)
		return
	}

	code, exitErr := waitExit(cmd, out, st, in)
	text := out.String()

	if verbose {
		fmt.Printf("--- %s: console mode %s, %d input event(s) left ---\n%s\n",
			st.name, modeText(in), pendingEvents(in), text)
	}
	ck.check(fmt.Sprintf("%s: exit code 0", st.name), code == 0 && exitErr == nil)
	if st.wantAsk {
		ck.check(st.name+": question asked once", strings.Count(text, st.ask) == 1)
	} else {
		ck.check(st.name+": no question", !strings.Contains(text, st.ask))
	}
	for _, w := range st.want {
		ck.check(fmt.Sprintf("%s: output has %q", st.name, w), strings.Contains(text, w))
	}
	for _, w := range st.notWant {
		ck.check(fmt.Sprintf("%s: output has no %q", st.name, w), !strings.Contains(text, w))
	}
	if code != 0 || exitErr != nil {
		fmt.Printf("  exit=%d err=%v output:\n%s\n", code, exitErr, tail(text, 12))
	}
}

// modeText names the console input mode flags that decide how a read behaves.
func modeText(in *os.File) string {
	var mode uint32
	r, _, _ := procGetConsoleMode.Call(in.Fd(), uintptr(unsafe.Pointer(&mode)))
	if r == 0 {
		return "<not a console>"
	}
	name := func(bit uint32, s string) string {
		if mode&bit != 0 {
			return s
		}
		return "-" + s
	}
	return fmt.Sprintf("[%s %s %s %s]", name(0x0001, "PROCESSED_INPUT"), name(0x0002, "LINE_INPUT"),
		name(0x0004, "ECHO_INPUT"), name(0x0100, "VT_INPUT"))
}

// pendingEvents reports how many console input records nobody has read yet.
func pendingEvents(in *os.File) uint32 {
	var n uint32
	_, _, _ = procGetNumberOfConsoleInputEvts.Call(in.Fd(), uintptr(unsafe.Pointer(&n)))
	return n
}

// waitExit waits for the prompt (and types the answer), then for the exit.
func waitExit(cmd *exec.Cmd, out *lockedBuffer, st step, in *os.File) (int, error) {
	typed := false
	deadline := time.Now().Add(timeoutFor(st))
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	for {
		select {
		case err := <-done:
			return exitCode(err), err
		default:
		}

		if !typed && st.answer != "" && strings.Contains(out.String(), st.ask) {
			typed = true
			time.Sleep(200 * time.Millisecond) // the prompt is on screen now
			if err := typeLine(in, st.answer); err != nil {
				fmt.Printf("  type: %v\n", err)
			}
		}
		if time.Now().After(deadline) {
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
			<-done
			return -1, fmt.Errorf("did not exit within %s", timeoutFor(st))
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func timeoutFor(st step) time.Duration {
	if st.answer != "" {
		return promptWait + exitWait
	}
	return promptWait
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return ee.ExitCode()
	}
	return -1
}

// newConsole detaches from the console the harness was started with and creates a
// private one, returning the console input handle for the harness (typing) and a
// second one to hand to the child as stdin.
func newConsole() (*os.File, *os.File, error) {
	_, _, _ = procFreeConsole.Call()
	ok, _, err := procAllocConsole.Call()
	if ok == 0 {
		return nil, nil, fmt.Errorf("AllocConsole: %v", err)
	}

	in, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("open CONIN$: %w", err)
	}
	childIn, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
	if err != nil {
		_ = in.Close()
		return nil, nil, fmt.Errorf("open CONIN$ for the child: %w", err)
	}
	return in, childIn, nil
}

// typeLine writes text plus Enter as console key events: the console hands them to
// whoever reads the input buffer, which is what a keyboard does.
func typeLine(in *os.File, text string) error {
	records := make([]inputRecord, 0, 2*len(text)+2)
	for _, r := range text {
		vk, state := keyFor(uint16(r))
		records = append(records, keyRecord(vk, uint16(r), state, true), keyRecord(vk, uint16(r), state, false))
	}
	// Enter: the console would otherwise look up a char for it, which fails for CR.
	records = append(records, keyRecord(vkReturn, vkReturn, 0, true), keyRecord(vkReturn, vkReturn, 0, false))

	var written uint32
	r, _, e := procWriteConsoleInputW.Call(in.Fd(), uintptr(unsafe.Pointer(&records[0])),
		uintptr(len(records)), uintptr(unsafe.Pointer(&written)))
	if r == 0 {
		return fmt.Errorf("WriteConsoleInput: %v", e)
	}
	return nil
}

// keyFor returns the virtual key and shift state a keyboard would report for the
// character. Inventing the virtual key (e.g. using '-' itself, which is
// VK_INSERT) makes the console act on a function key and drop the character.
func keyFor(char uint16) (uint16, uint32) {
	r, _, _ := procVkKeyScanW.Call(uintptr(char))
	if int32(r) == -1 { // no key on this layout
		return char, 0
	}
	return uint16(r & 0xff), uint32(r>>8) & 0xff // low byte VK, high byte shift state
}

// keyRecord builds one KEY_EVENT record.
func keyRecord(vk, char uint16, state uint32, down bool) inputRecord {
	rec := inputRecord{EventType: keyEventType}
	k := keyEventRecord{RepeatCount: 1, VirtualKeyCode: vk, UnicodeChar: char, ControlKeyState: state}
	if down {
		k.KeyDown = 1
	}
	copy(rec.Event[:], (*[16]byte)(unsafe.Pointer(&k))[:])
	return rec
}

// cleanEnv keeps the harness' environment but forces plain text output.
func cleanEnv(env []string) []string {
	out := make([]string, 0, len(env))
	for _, kv := range env {
		name, _, _ := strings.Cut(kv, "=")
		switch strings.ToUpper(name) {
		case "NO_COLOR", "TERM":
			continue
		}
		out = append(out, kv)
	}
	return out
}

type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// tail keeps the last n lines of the output.
func tail(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = append([]string{fmt.Sprintf("... %d earlier lines", len(lines)-n)}, lines[len(lines)-n:]...)
	}
	return strings.Join(lines, "\n")
}
