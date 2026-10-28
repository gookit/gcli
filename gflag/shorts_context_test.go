package gflag_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/x/assert"
)

func newContextParser(level uint8) (*gflag.Parser, *bool, *bool, *string) {
	p := gflag.New("context")
	p.ParserCfg().EnhanceShort = level
	var alpha, beta bool
	var name string
	p.BoolOpt(&alpha, "alpha", "a", false, "alpha")
	p.BoolOpt(&beta, "beta", "b", false, "beta")
	p.StrOpt(&name, "name", "n", "", "name")
	return p, &alpha, &beta, &name
}

func TestParser_EnhanceShortPreservesArgumentContext(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantName string
		wantArgs []string
		wantBool bool
	}{
		{"terminator bool data", []string{"--", "-ab"}, "", []string{"-ab"}, false},
		{"terminator attached data", []string{"--", "-nfoo"}, "", []string{"-nfoo"}, false},
		{"long option bool data", []string{"--name", "-ab"}, "-ab", nil, false},
		{"short option bool data", []string{"-n", "-ab"}, "-ab", nil, false},
		{"long option attached data", []string{"--name", "-nfoo"}, "-nfoo", nil, false},
		{"short option attached data", []string{"-n", "-nfoo"}, "-nfoo", nil, false},
		{"terminator as value", []string{"--name", "--", "-ab"}, "--", nil, true},
		{"equals value", []string{"--name=-ab"}, "-ab", nil, false},
		{"group after value", []string{"--name", "Ada", "-ab"}, "Ada", nil, true},
		{"group after equals", []string{"--name=-ab", "-ab"}, "-ab", nil, true},
	}
	for _, level := range []uint8{1, 2} {
		for _, reorder := range []bool{false, true} {
			for _, tt := range tests {
				t.Run(fmt.Sprintf("level%d/reorder%t/%s", level, reorder, tt.name), func(t *testing.T) {
					p, alpha, beta, name := newContextParser(level)
					p.WithConfigFn(gflag.WithReorderArgs(reorder))
					assert.NoErr(t, p.Parse(tt.args))
					assert.Eq(t, tt.wantName, *name)
					assert.Eq(t, tt.wantBool, *alpha)
					assert.Eq(t, tt.wantBool, *beta)
					if !slices.Equal(tt.wantArgs, p.RawArgs()) {
						t.Errorf("positional arguments: want %q, got %q", tt.wantArgs, p.RawArgs())
					}
				})
			}
		}
	}
}

func TestParser_EnhanceShortStopsAtCommandContext(t *testing.T) {
	for _, level := range []uint8{1, 2} {
		for _, subcommand := range []bool{false, true} {
			t.Run(fmt.Sprintf("level%d/subcommand%t", level, subcommand), func(t *testing.T) {
				p, alpha, beta, name := newContextParser(level)
				if subcommand {
					p.SetReorderStop(func(token string) bool { return token == "serve" })
				} else {
					p.WithConfigFn(gflag.WithReorderArgs(false))
				}
				args := []string{"serve", "--name", "-ab", "-nfoo"}
				assert.NoErr(t, p.Parse(args))
				assert.Eq(t, "", *name)
				assert.False(t, *alpha)
				assert.False(t, *beta)
				assert.Eq(t, args, p.RawArgs())
			})
		}
	}
}

func TestParser_EnhanceShortKeepsUnknownGroupErrors(t *testing.T) {
	for _, level := range []uint8{1, 2} {
		for _, arg := range []string{"-az", "-an"} {
			t.Run(fmt.Sprintf("level%d/%s", level, arg), func(t *testing.T) {
				p, alpha, beta, name := newContextParser(level)
				assert.Err(t, p.Parse([]string{arg}))
				assert.False(t, *alpha)
				assert.False(t, *beta)
				assert.Eq(t, "", *name)
			})
		}
	}
}

func TestParser_EnhanceShortKeepsNumericAliases(t *testing.T) {
	for _, level := range []uint8{1, 2} {
		for _, reorder := range []bool{false, true} {
			for _, value := range []bool{false, true} {
				t.Run(fmt.Sprintf("level%d/reorder%t/value%t", level, reorder, value), func(t *testing.T) {
					p, alpha, beta, name := newContextParser(level)
					p.WithConfigFn(gflag.WithReorderArgs(reorder))
					var one, two bool
					p.BoolOpt(&one, "one", "1", false, "one")
					p.BoolOpt(&two, "two", "2", false, "two")
					args := []string{"-12", "-ab"}
					if value {
						args = []string{"--name", "-12", "-ab"}
					}
					assert.NoErr(t, p.Parse(args))
					assert.Eq(t, !value, one)
					assert.Eq(t, !value, two)
					assert.True(t, *alpha)
					assert.True(t, *beta)
					if value {
						assert.Eq(t, "-12", *name)
					} else {
						assert.Eq(t, "", *name)
					}
				})
			}
		}
	}

	for _, reorder := range []bool{false, true} {
		t.Run(fmt.Sprintf("attached/reorder%t", reorder), func(t *testing.T) {
			p := gflag.New("context")
			p.ParserCfg().EnhanceShort = 2
			p.WithConfigFn(gflag.WithReorderArgs(reorder))
			var name string
			p.StrOpt(&name, "name", "1", "", "name")
			assert.NoErr(t, p.Parse([]string{"-1Ada"}))
			assert.Eq(t, "Ada", name)
		})
	}
}

func TestParser_EnhanceShortPreservesNumericTails(t *testing.T) {
	for _, level := range []uint8{1, 2} {
		for _, boundary := range []string{"terminator", "positional", "subcommand"} {
			t.Run(fmt.Sprintf("level%d/%s", level, boundary), func(t *testing.T) {
				p, alpha, beta, name := newContextParser(level)
				var one, two bool
				p.BoolOpt(&one, "one", "1", false, "one")
				p.BoolOpt(&two, "two", "2", false, "two")
				args := []string{"--", "-12", "-1.5", "-.5", "-ab"}
				wantArgs := args[1:]
				if boundary != "terminator" {
					args[0] = "serve"
					wantArgs = args
					if boundary == "positional" {
						p.WithConfigFn(gflag.WithReorderArgs(false))
					} else {
						p.SetReorderStop(func(token string) bool { return token == "serve" })
					}
				}
				assert.NoErr(t, p.Parse(args))
				assert.False(t, one)
				assert.False(t, two)
				assert.False(t, *alpha)
				assert.False(t, *beta)
				assert.Eq(t, "", *name)
				assert.Eq(t, wantArgs, p.RawArgs())
			})
		}
	}
}
