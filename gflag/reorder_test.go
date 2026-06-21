package gflag

import (
	"testing"

	"github.com/gookit/goutil/x/assert"
)

// build a parser with: -n/--name, -o/--out (value opts); -v/--verbose, -f/--force (bool opts)
func newReorderParser() *Flags {
	p := New("test")
	var name, out string
	var verbose, force bool
	p.StrOpt(&name, "name", "n", "", "name opt")
	p.StrOpt(&out, "out", "o", "", "out opt")
	p.BoolOpt(&verbose, "verbose", "v", false, "verbose opt")
	p.BoolOpt(&force, "force", "f", false, "force opt")
	return p
}

func TestRearrangeArgs(t *testing.T) {
	fs := newReorderParser().fSet

	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"arg then long opt", []string{"arg1", "--name", "tom"}, []string{"--name", "tom", "arg1"}},
		{"value opt groups its value", []string{"file1", "-n", "tom", "file2"}, []string{"-n", "tom", "file1", "file2"}},
		{"bool opt not consume next", []string{"arg", "-v", "x"}, []string{"-v", "arg", "x"}},
		{"eq form not consume next", []string{"a", "--name=tom", "b"}, []string{"--name=tom", "a", "b"}},
		{"already canonical", []string{"--name", "tom", "a", "b"}, []string{"--name", "tom", "a", "b"}},
		{"mixed multi", []string{"a", "-v", "b", "--out", "log", "c"}, []string{"-v", "--out", "log", "a", "b", "c"}},
		{"negative number is arg", []string{"-n", "tom", "-5"}, []string{"-n", "tom", "-5"}},
		{"lone dash is arg", []string{"-", "-v"}, []string{"-v", "-"}},
		{"double dash stops reorder", []string{"a", "--", "-n", "b"}, []string{"a", "--", "-n", "b"}},
		{"unknown opt not consume next", []string{"x", "--unknown", "y"}, []string{"--unknown", "x", "y"}},
		{"single token unchanged", []string{"arg"}, []string{"arg"}},
		{"empty unchanged", []string{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Eq(t, tt.want, rearrangeArgs(tt.in, fs, nil))
		})
	}
}

func TestRearrangeArgs_stopAt(t *testing.T) {
	fs := newReorderParser().fSet
	isSub := func(name string) bool { return name == "push" }

	t.Run("stop at subcommand, keep sub opts verbatim", func(t *testing.T) {
		got := rearrangeArgs([]string{"-v", "push", "--force", "origin"}, fs, isSub)
		assert.Eq(t, []string{"-v", "push", "--force", "origin"}, got)
	})

	t.Run("stop even when arg precedes subcommand", func(t *testing.T) {
		got := rearrangeArgs([]string{"push", "x", "-v"}, fs, isSub)
		assert.Eq(t, []string{"push", "x", "-v"}, got)
	})

	t.Run("reorder before reaching subcommand", func(t *testing.T) {
		got := rearrangeArgs([]string{"arg", "-v", "push", "-f"}, fs, isSub)
		assert.Eq(t, []string{"-v", "arg", "push", "-f"}, got)
	})
}
