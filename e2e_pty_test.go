package gcli_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/gookit/gopty/ptytest"
)

// The pty end-to-end cases below drive the real example binary
// (./_examples/cliapp) through a pseudo terminal: they cover what unit tests
// cannot — the prompt really appears on a terminal, what is typed becomes the
// option value, and a value given on the command line skips the question.
//
// The trickier semantics (which command asks, how often, that stdin is untouched
// when a value is given) are pinned without a terminal in
// shared_opts_test.go/gflag tests; these cases add the terminal layer itself.
//
// They run on unix (creack/pty) and on Windows (ConPTY, see gookit/gopty), so CI
// needs no Python and no Windows console harness. `go test -short` skips them.

// e2eBin is the example binary, built once per test process on first use.
var (
	e2eOnce sync.Once
	e2eBin  string
	e2eErr  error
)

// buildExampleApp builds ./_examples/cliapp into the repo's tmp/ (ignored) and
// returns its absolute path.
func buildExampleApp() (string, error) {
	root, err := moduleRoot()
	if err != nil {
		return "", err
	}
	outDir := filepath.Join(root, "tmp", "e2e")
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return "", err
	}
	name := "cliapp"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	out := filepath.Join(outDir, name)

	cmd := exec.Command("go", "build", "-o", out, "./_examples/cliapp")
	cmd.Dir = root
	if buildOut, berr := cmd.CombinedOutput(); berr != nil {
		return "", fmt.Errorf("go build ./_examples/cliapp: %v\n%s", berr, buildOut)
	}
	return out, nil
}

// moduleRoot is the directory holding this package (go test runs with the package
// directory as the working directory, but resolving it from the source file keeps
// the file paths absolute).
func moduleRoot() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot locate the test source file")
	}
	return filepath.Dir(file), nil
}

// exampleBin returns the built example binary, skipping the cases when they cannot
// run instead of reporting a build problem as a terminal-behaviour failure.
//
// The build happens once per test process, on the first case that needs it, so
// `go test -short` pays nothing for it.
func exampleBin(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("pty e2e cases are skipped with -short")
	}
	e2eOnce.Do(func() { e2eBin, e2eErr = buildExampleApp() })
	if e2eErr != nil {
		t.Fatalf("example app: %v", e2eErr)
	}
	return e2eBin
}

// TestE2E_askDemo covers a local option collected with CliOpt.Question.
func TestE2E_askDemo(t *testing.T) {
	bin := exampleBin(t)
	const ask = `Please input your access token:`

	t.Run("empty value: asked once, typed text used", func(t *testing.T) {
		s := ptytest.Start(t, bin, "ask-demo")
		m := s.Mark()
		if !s.Expect(ask, m) {
			t.Fatalf("the question never appeared, output:\n%s", s.Text())
		}
		s.SendLine("e2e-token")
		code, err := s.WaitExit()
		if err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if !s.Contains(`e2e-token`, m) {
			t.Fatalf("the typed text was not used, output:\n%s", s.Since(m))
		}
		if n := s.Count(ask, m); n != 1 {
			t.Fatalf("the question appeared %d times, want 1, output:\n%s", n, s.Since(m))
		}
	})

	t.Run("value on the command line: no question", func(t *testing.T) {
		s := ptytest.Start(t, bin, "ask-demo", "--token=cli-token")
		m := s.Mark()
		code, err := s.WaitExit()
		if err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if s.Contains(ask, m) {
			t.Fatalf("the question was asked although --token was given, output:\n%s", s.Since(m))
		}
		if !s.Contains(`cli-token`, m) {
			t.Fatalf("the given value was not used, output:\n%s", s.Since(m))
		}
	})

	t.Run("-h: option listed, no question", func(t *testing.T) {
		s := ptytest.Start(t, bin, "ask-demo", "-h")
		m := s.Mark()
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if !s.Contains(`--token`, m) {
			t.Fatalf("help does not list --token, output:\n%s", s.Since(m))
		}
		if s.Contains(ask, m) {
			t.Fatalf("help asked the question, output:\n%s", s.Since(m))
		}
	})
}

// TestE2E_sharedAsk covers the shared option whose collection is deferred to the
// executing command (Command.prepareSharedOpts): asked once, on the leaf, and never
// when any command segment already carries the value.
func TestE2E_sharedAsk(t *testing.T) {
	bin := exampleBin(t)
	const ask = `Please input your api token:`

	t.Run("empty value: asked once at the executed command", func(t *testing.T) {
		s := ptytest.Start(t, bin, "shared-ask", "run")
		m := s.Mark()
		if !s.Expect(ask, m) {
			t.Fatalf("the question never appeared, output:\n%s", s.Text())
		}
		s.SendLine("e2e-shared")
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if !s.Contains(`leaf: token=e2e-shared`, m) {
			t.Fatalf("the leaf did not see the collected value, output:\n%s", s.Since(m))
		}
		if n := s.Count(ask, m); n != 1 {
			t.Fatalf("the question appeared %d times, want 1, output:\n%s", n, s.Since(m))
		}
	})

	t.Run("value in the leaf segment: no question", func(t *testing.T) {
		s := ptytest.Start(t, bin, "shared-ask", "run", "--token=leaf-token")
		m := s.Mark()
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if s.Contains(ask, m) {
			t.Fatalf("asked although the leaf segment had the value, output:\n%s", s.Since(m))
		}
		if !s.Contains(`leaf: token=leaf-token`, m) {
			t.Fatalf("the leaf did not see the given value, output:\n%s", s.Since(m))
		}
	})

	t.Run("value in the parent segment: no question", func(t *testing.T) {
		s := ptytest.Start(t, bin, "shared-ask", "--token=parent-token", "run")
		m := s.Mark()
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if s.Contains(ask, m) {
			t.Fatalf("asked although the parent segment had the value, output:\n%s", s.Since(m))
		}
		if !s.Contains(`leaf: token=parent-token`, m) {
			t.Fatalf("the leaf did not see the given value, output:\n%s", s.Since(m))
		}
	})

	t.Run("-h on the leaf: inherited option listed, no question", func(t *testing.T) {
		s := ptytest.Start(t, bin, "shared-ask", "run", "-h")
		m := s.Mark()
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if !s.Contains(`--token`, m) {
			t.Fatalf("help does not list the inherited --token, output:\n%s", s.Since(m))
		}
		if s.Contains(ask, m) {
			t.Fatalf("help asked the question, output:\n%s", s.Since(m))
		}
	})
}

// TestE2E_cliuiPrompts covers cliui's line-based prompts inside a gcli app: the
// builtin completion command asks with interact.Unconfirmed/AnswerIsYes, and the
// example's read demo echoes one line back.
func TestE2E_cliuiPrompts(t *testing.T) {
	bin := exampleBin(t)
	const confirm = `Ensure continue \[yes\|no\]`

	t.Run("confirm: n is unconfirmed", func(t *testing.T) {
		s := ptytest.Start(t, bin, "interact", "confirm")
		m := s.Mark()
		if !s.Expect(confirm, m) {
			t.Fatalf("the confirmation never appeared, output:\n%s", s.Text())
		}
		s.SendLine("n")
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if !s.Contains(`Unconfirmed`, m) {
			t.Fatalf("n did not take the unconfirmed path, output:\n%s", s.Since(m))
		}
	})

	t.Run("confirm: y is confirmed", func(t *testing.T) {
		s := ptytest.Start(t, bin, "interact", "confirm")
		m := s.Mark()
		if !s.Expect(confirm, m) {
			t.Fatalf("the confirmation never appeared, output:\n%s", s.Text())
		}
		s.SendLine("y")
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if !s.Contains(`Confirmed`, m) {
			t.Fatalf("y did not take the confirmed path, output:\n%s", s.Since(m))
		}
		if s.Contains(`Please try again`, m) {
			t.Fatalf("the answer was rejected, output:\n%s", s.Since(m))
		}
	})

	t.Run("read: input echoed back", func(t *testing.T) {
		s := ptytest.Start(t, bin, "interact", "read")
		m := s.Mark()
		if !s.Expect(`Your name\?`, m, 20*time.Second) {
			t.Fatalf("the question never appeared, output:\n%s", s.Text())
		}
		s.SendLine("gopher")
		if code, err := s.WaitExit(); err != nil || code != 0 {
			t.Fatalf("exit = (%d, %v), output:\n%s", code, err, s.Text())
		}
		if !s.Contains(`Your input:\s+gopher`, m) {
			t.Fatalf("the input was not echoed back, output:\n%s", s.Since(m))
		}
	})
}
