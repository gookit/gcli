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

// native map[string]string field can be bound from struct tags (D1.3)
func TestFlags_FromStruct_nativeMap(t *testing.T) {
	type opts struct {
		Meta map[string]string `flag:"name=meta;shorts=m;desc=key-value metadata"`
	}

	o := &opts{}
	fs := gflag.New("test")
	assert.NoErr(t, fs.FromStruct(o))
	assert.True(t, fs.HasOption("meta"))

	err := fs.Parse([]string{"--meta", "k1=v1", "-m", "k2=v2"})
	assert.NoErr(t, err)
	assert.Len(t, o.Meta, 2)
	assert.Eq(t, "v1", o.Meta["k1"])
	assert.Eq(t, "v2", o.Meta["k2"])
}

// enum tag sets value candidates(choices) + membership validation (D1.4)
func TestFlags_FromStruct_enum(t *testing.T) {
	type opts struct {
		Lang string `flag:"name=lang;shorts=l;desc=language;enum=go,php,java"`
	}

	t.Run("choices populated for completion", func(t *testing.T) {
		fs := gflag.New("test")
		assert.NoErr(t, fs.FromStruct(&opts{}))
		assert.Eq(t, []string{"go", "php", "java"}, fs.Opt("lang").Choices)
	})

	t.Run("valid value passes", func(t *testing.T) {
		o := &opts{}
		fs := gflag.New("test")
		assert.NoErr(t, fs.FromStruct(o))
		assert.NoErr(t, fs.Parse([]string{"--lang", "go"}))
		assert.Eq(t, "go", o.Lang)
	})

	t.Run("invalid value rejected", func(t *testing.T) {
		fs := gflag.New("test")
		assert.NoErr(t, fs.FromStruct(&opts{}))
		err := fs.Parse([]string{"--lang", "ruby"})
		assert.Err(t, err)
		assert.StrContains(t, err.Error(), "allowed list")
	})
}

// 匿名内嵌*未导出类型*也必须被展开 (回归: 旧实现里未导出字段名跳过会在展开前把它丢掉,
// 导致 _examples 的 baseFlags / 文档里的 commonOpts 等小写内嵌的选项实际未生成)。
func TestFlags_FromStruct_anonymousUnexportedEmbed(t *testing.T) {
	type baseNamed struct {
		Verbose bool `flag:"name=verbose;shorts=v;desc=verbose"`
	}
	type baseField struct {
		Verbose bool `flag:"v" desc:"verbose"`
	}

	t.Run("TagRuleNamed expands unexported embed", func(t *testing.T) {
		type opts struct {
			baseNamed
			Name string `flag:"name=name;shorts=n;desc=name"`
		}
		o := &opts{}
		fs := gflag.New("test")
		assert.NoErr(t, fs.FromStruct(o))
		assert.True(t, fs.HasOption("verbose"))
		assert.True(t, fs.HasOption("name"))

		assert.NoErr(t, fs.Parse([]string{"-v", "--name", "tom"}))
		assert.True(t, o.Verbose)
		assert.Eq(t, "tom", o.Name)
	})

	t.Run("TagRuleField expands unexported embed", func(t *testing.T) {
		type opts struct {
			baseField
			UserName string `flag:"u" desc:"user name"`
		}
		o := &opts{}
		fs := gflag.New("test")
		assert.NoErr(t, fs.FromStruct(o, gflag.TagRuleField))
		assert.True(t, fs.HasOption("verbose"))
		assert.True(t, fs.HasOption("user-name"))

		assert.NoErr(t, fs.Parse([]string{"-v", "-u", "tom"}))
		assert.True(t, o.Verbose)
		assert.Eq(t, "tom", o.UserName)
	})
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
