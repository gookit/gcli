package gflag

import (
	"testing"

	"github.com/gookit/goutil/x/assert"
)

func TestExpandShortArgs(t *testing.T) {
	shorts := map[string]string{"a": "aaa", "u": "uuu", "x": "xxx", "O": "out", "b": "bbb"}
	isBool := func(s string) bool {
		return map[string]bool{"a": true, "u": true, "x": true}[s]
	}

	tests := []struct {
		name  string
		level uint8
		in    []string
		want  []string
	}{
		{"level0 raw", 0, []string{"-aux"}, []string{"-aux"}},
		{"level1 all bool split", 1, []string{"-aux"}, []string{"-a", "-u", "-x"}},
		{"level1 mixed not split", 1, []string{"-aO"}, []string{"-aO"}},
		{"level1 no attached", 1, []string{"-Ostdout"}, []string{"-Ostdout"}},
		{"level2 attached value", 2, []string{"-Ostdout"}, []string{"-O", "stdout"}},
		{"level2 still split all bool", 2, []string{"-aux"}, []string{"-a", "-u", "-x"}},
		{"level2 with =", 2, []string{"-a=1"}, []string{"-a=1"}},
		{"level2 long opt", 2, []string{"--name"}, []string{"--name"}},
		{"level2 terminator", 2, []string{"--"}, []string{"--"}},
		{"level2 single short", 2, []string{"-O"}, []string{"-O"}},
		{"level2 unknown short", 2, []string{"-zzz"}, []string{"-zzz"}},
		{
			"level1 mixed args", 1,
			[]string{"-aux", "val", "--name", "x"},
			[]string{"-a", "-u", "-x", "val", "--name", "x"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandShortArgs(tt.in, shorts, isBool, tt.level)
			assert.Eq(t, tt.want, got)
		})
	}
}

func TestParser_EnhanceShort(t *testing.T) {
	newParser := func(level uint8) (*Parser, *bool, *bool, *bool, *string) {
		p := New("test")
		p.ParserCfg().EnhanceShort = level

		var a, u, x bool
		var out string
		p.BoolOpt(&a, "aaa", "a", false, "bool a")
		p.BoolOpt(&u, "uuu", "u", false, "bool u")
		p.BoolOpt(&x, "xxx", "x", false, "bool x")
		p.StrOpt(&out, "out", "O", "", "string out")
		return p, &a, &u, &x, &out
	}

	t.Run("level0 not defined error", func(t *testing.T) {
		p, _, _, _, _ := newParser(0)
		err := p.Parse([]string{"-aux"})
		assert.Err(t, err)
	})

	t.Run("level1 multi bool split", func(t *testing.T) {
		p, a, u, x, _ := newParser(1)
		err := p.Parse([]string{"-aux"})
		assert.NoErr(t, err)
		assert.True(t, *a)
		assert.True(t, *u)
		assert.True(t, *x)
	})

	t.Run("level2 attached value", func(t *testing.T) {
		p, _, _, _, out := newParser(2)
		err := p.Parse([]string{"-Ostdout"})
		assert.NoErr(t, err)
		assert.Eq(t, "stdout", *out)
	})
}
