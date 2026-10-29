package gflag

import (
	"testing"

	"github.com/gookit/goutil/x/assert"
)

func TestStrictFormatArgs(t *testing.T) {
	p := New("strict")
	var a, b, v bool
	var name, out string
	p.BoolOpt(&a, "alpha", "a", false, "bool alpha")
	p.BoolOpt(&b, "beta", "b", false, "bool beta")
	p.BoolOpt(&v, "verbose", "v", false, "bool verbose")
	p.StrOpt(&name, "name", "n", "", "string name")
	p.StrOpt(&out, "out", "O", "", "string out")
	fs := p.FSet()

	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"empty", nil, nil},
		// long option form normalization
		{"two dashes one char", []string{"--a"}, []string{"-a"}},
		{"three dashes", []string{"---name", "tom"}, []string{"--name", "tom"}},
		{"already normal", []string{"--name", "tom"}, []string{"--name", "tom"}},
		// groups and attached values are left to ExpandShortArgs
		{"short group", []string{"-ab"}, []string{"-ab"}},
		{"attached value", []string{"-Oab"}, []string{"-Oab"}},
		// the value of a value-taking option is not an option
		{"long value one char", []string{"--name", "--a"}, []string{"--name", "--a"}},
		{"long value three dashes", []string{"--name", "---deep"}, []string{"--name", "---deep"}},
		{"short value three dashes", []string{"-n", "---deep"}, []string{"-n", "---deep"}},
		{"short value one char", []string{"-O", "--x"}, []string{"-O", "--x"}},
		{"value after normalized option", []string{"---name", "--deep"}, []string{"--name", "--deep"}},
		{"attached value with =", []string{"--name=--a", "---d"}, []string{"--name=--a", "-d"}},
		// a bool option does not take a value
		{"bool then option", []string{"-v", "--deep"}, []string{"-v", "--deep"}},
		{"bool then malformed option", []string{"-v", "---deep"}, []string{"-v", "--deep"}},
		// unknown options are reported by the parser, but their form is normalized
		{"unknown option", []string{"--unknown", "---x"}, []string{"--unknown", "-x"}},
		// the terminator keeps the rest verbatim
		{"terminator", []string{"--", "--a", "---b"}, []string{"--", "--a", "---b"}},
		{"terminator after value", []string{"--name", "--", "---b"}, []string{"--name", "--", "---b"}},
		// a lone dash is positional
		{"lone dash", []string{"-"}, []string{"-"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Eq(t, tt.want, StrictFormatArgs(tt.in, fs))
		})
	}

	// without a flag set only the form is normalized
	assert.Eq(t, []string{"-a"}, StrictFormatArgs([]string{"--a"}, nil))
	assert.Eq(t, []string{"--name", "-a"}, StrictFormatArgs([]string{"--name", "--a"}, nil))
}
