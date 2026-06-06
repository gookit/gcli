package gcli_test

import (
	"runtime"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/x/assert"
)

func TestGcliBasic(t *testing.T) {
	is := assert.New(t)
	is.NotEmpty(gcli.Version())
	is.NotEmpty(gcli.CommitID())
}

func TestVerbose(t *testing.T) {
	is := assert.New(t)

	old := gcli.Verbose()
	is.Eq(gcli.VerbError, old)
	is.False(gcli.GOpts().NoColor)

	gcli.SetDebugMode()
	is.Eq(gcli.VerbDebug, gcli.Verbose())

	gcli.SetQuietMode()
	is.Eq(gcli.VerbQuiet, gcli.Verbose())

	info := gcli.VerbInfo
	gcli.SetVerbose(info)
	is.Eq(info, gcli.Verbose())
	is.Eq(3, info.Int())
	is.Eq("info", info.Name())
	is.Eq("INFO", info.Upper())

	gcli.SetVerbose(old)
	is.Eq(gcli.VerbError, gcli.Verbose())
}

func TestVerbLevel(t *testing.T) {
	is := assert.New(t)

	verb := gcli.VerbLevel(23)
	is.Eq("unknown", verb.Name())
	is.Eq(23, verb.Int())

	err := verb.Set("2")
	is.NoErr(err)
	is.Eq(gcli.VerbWarn, verb)
	is.Eq("warn", verb.Name())

	err = verb.Set("debug")
	is.NoErr(err)
	is.Eq(gcli.VerbDebug, verb)
	is.Eq("debug", verb.Name())

	err = verb.Set("30")
	is.NoErr(err)
	is.Eq(gcli.VerbCrazy, verb)
	is.Eq("crazy", verb.Name())
}

func TestStrictMode(t *testing.T) {
	is := assert.New(t)

	old := gcli.StrictMode()
	defer func() {
		gcli.SetStrictMode(old)
	}()

	gcli.SetStrictMode(false)
	is.False(gcli.StrictMode())

	gcli.SetStrictMode(true)
	is.True(gcli.StrictMode())
}

func TestNewCtx(t *testing.T) {
	is := assert.New(t)

	ctx := gcli.NewCtx()
	ctx.InitCtx()

	is.True(ctx.PID() > 0)
	is.Eq(runtime.GOOS, ctx.OsName())
	is.NotEmpty(ctx.BinName())
	is.NotEmpty(ctx.WorkDir())

	args := ctx.OsArgs()
	is.NotEmpty(args)

	if len(args) > 1 {
		is.NotEmpty(ctx.ArgLine())
	} else {
		is.Empty(ctx.ArgLine())
	}
}

func TestSetStrictMode(t *testing.T) {
	stm := gcli.StrictMode()
	defer gcli.SetStrictMode(stm)

	opts := struct {
		name   string
		ok, bl bool
	}{}

	// gcli.SetVerbose(gcli.VerbDebug)
	app := gcli.NewApp(gcli.NotExitOnEnd())
	app.Add(&gcli.Command{
		Name: "test",
		Config: func(c *gcli.Command) {
			c.StrOpt(&opts.name, "name", "n", "", "1")
			c.BoolOpt(&opts.ok, "ok", "o", true, "2")
			c.BoolOpt(&opts.bl, "bl", "b", false, "3")
		},
		Func: func(c *gcli.Command, _ []string) error {
			return nil
		},
	})

	app.Run([]string{"test", "-o", "-n", "inhere"})
	assert.Eq(t, "inhere", opts.name)
	assert.True(t, opts.ok)

	app.Run([]string{"test", "-o=false", "-n=inhere"})
	assert.Eq(t, "inhere", opts.name)
	assert.False(t, opts.ok)

	app.Run([]string{"test", "-ob"})
	// assert.StrContains(t, errMsg, "ddd")

	gcli.SetStrictMode(true)
	app.Run([]string{"test", "-ob"})
	assert.True(t, opts.ok)
	assert.True(t, opts.bl)
}

// TestStrictMode_safeShortSplit 验证 strictMode 收敛为驱动 gflag EnhanceShort 后,
// 取值型短选项不再被旧的"盲拆"逻辑错误拆分(回归 B4+B5)。
func TestStrictMode_safeShortSplit(t *testing.T) {
	stm := gcli.StrictMode()
	defer gcli.SetStrictMode(stm)

	// 每个用例独立构建 app/opts，避免上一次 Run 的状态干扰
	build := func() (*gcli.App, *struct {
		out     string
		verbose bool
	}) {
		opts := &struct {
			out     string
			verbose bool
		}{}
		app := gcli.NewApp(gcli.NotExitOnEnd())
		app.Add(&gcli.Command{
			Name: "test",
			Config: func(c *gcli.Command) {
				c.StrOpt(&opts.out, "out", "O", "", "out file")      // 取值型短选项
				c.BoolOpt(&opts.verbose, "verbose", "V", false, "v") // bool 短选项
			},
			Func: func(c *gcli.Command, _ []string) error { return nil },
		})
		return app, opts
	}

	gcli.SetStrictMode(true)

	// 混合写法: -VO 含取值型短选项 O，非全 bool。
	// 旧盲拆会错拆为 -V -O 并把后续 "val" 吞作 O 的值(verbose=true,out="val");
	// 新方案下 -VO 不拆，整体作为未知长选项解析失败，opt 值不被污染(保持零值)。
	t.Run("mixed value-short not blindly split", func(t *testing.T) {
		app, opts := build()
		app.Run([]string{"test", "-VO", "val"})
		assert.False(t, opts.verbose) // 未被盲拆污染
		assert.Eq(t, "", opts.out)    // "val" 未被错误吞为 O 的值
	})

	// 纯 bool 组合: 拆分仍正常工作。
	t.Run("all-bool short still split", func(t *testing.T) {
		app, opts := build()
		app.Run([]string{"test", "-VV"})
		assert.True(t, opts.verbose)
	})

	// 单个 bool 短选项正常。
	t.Run("single bool short", func(t *testing.T) {
		app, opts := build()
		app.Run([]string{"test", "-V"})
		assert.True(t, opts.verbose)
	})

	// 取值型短选项分开写正常取值，未受影响。
	t.Run("value-short with separate value", func(t *testing.T) {
		app, opts := build()
		app.Run([]string{"test", "-V", "-O", "val"})
		assert.True(t, opts.verbose)
		assert.Eq(t, "val", opts.out)
	})
}

func TestString(t *testing.T) {
	s := gcli.String("ab,cd")
	assert.Eq(t, "ab,cd", s.String())
	assert.Eq(t, []string{"ab", "cd"}, s.Split(","))
}
