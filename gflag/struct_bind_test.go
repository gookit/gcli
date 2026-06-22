package gflag_test

import (
	"testing"
	"time"

	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/x/assert"
)

// native slice + time.Duration fields can be bound from struct tags (D1.2)
func TestFlags_FromStruct_nativeSliceAndDuration(t *testing.T) {
	type opts struct {
		Names []string      `flag:"name=names;shorts=n;desc=name list"`
		Ports []int         `flag:"name=ports;shorts=p;desc=port list"`
		Flags []bool        `flag:"name=flags;desc=bool list"`
		TTL   time.Duration `flag:"name=ttl;desc=time to live"`
	}

	o := &opts{}
	fs := gflag.New("test")
	assert.NoErr(t, fs.FromStruct(o))
	assert.True(t, fs.HasOption("names"))
	assert.True(t, fs.HasOption("ttl"))

	err := fs.Parse([]string{
		"--names", "a", "-n", "b",
		"--ports", "80", "-p", "443",
		"--flags", "true", "--flags", "false",
		"--ttl", "1h30m",
	})
	assert.NoErr(t, err)
	assert.Eq(t, []string{"a", "b"}, o.Names)
	assert.Eq(t, []int{80, 443}, o.Ports)
	assert.Eq(t, []bool{true, false}, o.Flags)
	assert.Eq(t, 90*time.Minute, o.TTL)
}

// unsupported slice elem type reports a clear error, not a panic
func TestFlags_FromStruct_unsupportedSlice(t *testing.T) {
	type opts struct {
		Rates []float64 `flag:"name=rates;desc=rate list"`
	}

	fs := gflag.New("test")
	err := fs.FromStruct(&opts{})
	assert.Err(t, err)
	assert.StrContains(t, err.Error(), "unsupport slice type")
}
