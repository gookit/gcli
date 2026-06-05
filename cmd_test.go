package gcli_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/gookit/goutil/x/assert"
)

var simpleArgs = []string{"hi"}

func TestNewCommand(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.Aliases = []string{"alias1", "alias2"}
	})

	is.NotEmpty(c)
	is.False(c.IsRunnable())
	is.Nil(c.App())

	err := c.Run([]string{})
	is.NoErr(err)
	is.True(c.IsStandalone())
	is.False(c.NotStandalone())

	is.False(c.IsDisabled())
	c.Disable()
	is.True(c.IsDisabled())

	// is.Eq("", c.ArgLine())
	is.Eq("alias1,alias2", c.Aliases.String())
	is.Eq("alias1alias2", c.Aliases.Join(""))

	c = gcli.NewCommand("test1", "desc test")
	app := gcli.NewApp()
	c.AttachTo(app)
	is.True(app.HasCommand("test1"))
}

func TestCommand_NewErrf(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		c.AddArg("arg0", "desc message")
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		is.Eq("hi", c.Arg("arg0").String())
		return c.NewErrf("error message")
	})

	is.NotEmpty(c)
	is.True(c.IsRunnable())

	err := c.Run(simpleArgs)
	is.Err(err)
	is.Eq("error message", err.Error())
	is.Eq([]string{"hi"}, c.RawArgs())

	is.NotPanics(func() {
		c.MustRun(simpleArgs)
	})
}

func TestCommand_AddArg(t *testing.T) {
	is := assert.New(t)
	c := gcli.NewCommand("test", "test desc", nil)

	arg := c.AddArg("arg0", "arg desc", true)
	is.Eq(0, arg.Index())

	ret := c.ArgByIndex(0)
	is.Eq(ret, arg)

	assert.PanicsMsg(t, func() {
		c.ArgByIndex(1)
	}, "gflag: get not exists argument #1")

	arg = c.AddArg("arg1", "arg1 desc")
	is.Eq(1, arg.Index())

	ret = c.Arg("arg1")
	is.Eq(ret, arg)

	is.PanicsMsg(func() {
		c.Arg("not-exist")
	}, "gflag: get not exists argument 'not-exist'")

	is.Len(c.Args(), 2)

	is.PanicsMsg(func() {
		c.AddArg("", "desc")
	}, "gflag: the command argument name cannot be empty")

	is.PanicsMsg(func() {
		c.AddArg(":)&dfd", "desc")
	}, "gflag: the argument name ':)&dfd' is invalid, must match: ^[0-9a-zA-Z][\\w-]*$")

	is.PanicsMsg(func() {
		c.AddArg("arg1", "desc")
	}, "gflag: the argument name 'arg1' already exists in command 'test'")
	is.PanicsMsg(func() {
		c.AddArg("arg2", "arg2 desc", true)
	}, "gflag: required argument 'arg2' cannot be defined after optional argument")

	c.AddArg("arg3", "arg3 desc", false, true)
	is.PanicsMsg(func() {
		c.AddArg("argN", "desc", true)
	}, "gflag: have defined an array argument, you cannot add argument 'argN'")
}

func TestCommand_Run(t *testing.T) {
	is := assert.New(t)

	// use struct
	c := &gcli.Command{
		Name: "test",
		Desc: "test desc",
		Config: func(c *gcli.Command) {
			is.Eq("test", c.Name)
			c.Aliases = []string{"alias1"}
		},
		Func: func(c *gcli.Command, args []string) error {
			return nil
		},
	}

	err := c.Run([]string{})
	is.NoErr(err)
	is.True(c.IsStandalone())
	is.False(c.NotStandalone())
	is.Eq("alias1", c.Aliases.String())

	err = c.Run([]string{"-h"})
	is.NoErr(err)
}

func TestCommand_Run_recoverPanic(t *testing.T) {
	// c.Run 会改包级 gOpts, 跑完后重置
	defer gcli.ResetGOpts()

	// 1. 主函数 panic 被恢复成 error
	t.Run("func panic", func(t *testing.T) {
		c := gcli.NewCommand("test-panic-func", "desc")
		c.SetFunc(func(c *gcli.Command, args []string) error {
			panic("boom func")
		})

		var err error
		assert.NotPanics(t, func() {
			err = c.Run([]string{})
		})
		assert.Err(t, err)
		assert.StrContains(t, err.Error(), "boom func")
	})

	// 2. 中间件 panic 同样被恢复, 且主函数不会被执行
	t.Run("middleware panic", func(t *testing.T) {
		funcCalled := false
		c := gcli.NewCommand("test-panic-mw", "desc")
		c.Use(func(c *gcli.Command, args []string) error {
			panic("boom mw")
		})
		c.SetFunc(func(c *gcli.Command, args []string) error {
			funcCalled = true
			return nil
		})

		var err error
		assert.NotPanics(t, func() {
			err = c.Run([]string{})
		})
		assert.Err(t, err)
		assert.StrContains(t, err.Error(), "boom mw")
		// 中间件 panic 后主函数不应被执行
		assert.False(t, funcCalled)
	})

	// 3. panic 为 error 类型时原样返回
	t.Run("panic with error", func(t *testing.T) {
		wantErr := errors.New("xx panic err")
		c := gcli.NewCommand("test-panic-err", "desc")
		c.SetFunc(func(c *gcli.Command, args []string) error {
			panic(wantErr)
		})

		var err error
		assert.NotPanics(t, func() {
			err = c.Run([]string{})
		})
		assert.Eq(t, wantErr, err)
	})
}

func TestNewCommand_Run(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("test", "desc test", func(c *gcli.Command) {
		is.Eq("test", c.Name)
		c.Aliases = []string{"alias1"}
	})
	c.SetFunc(func(c *gcli.Command, args []string) error {
		return nil
	})

	is.NotEmpty(c)
	err := c.Run([]string{})
	is.NoErr(err)
	is.True(c.IsStandalone())
	is.False(c.NotStandalone())

	is.Eq("alias1", c.Aliases.String())

	err = c.Run([]string{"-h"})
	is.NoErr(err)
	is.Eq("alias1", c.Aliases.String())

	// error run on app
	g := gcli.NewApp()
	g.AddCommand(c)
	err = c.Run(simpleArgs)
	is.NoErr(err)
}

// newGitCmd 构造一个与原共享 r 等价的 git 命令树, 每个用例独立调用获得全新实例,
// 各 Func 闭包写入返回的局部 bf, 避免跨用例共享可变状态(支持 -shuffle=on)。
func newGitCmd() (*gcli.Command, *bytes.Buffer) {
	bf := new(bytes.Buffer)

	// l0: root command
	r := &gcli.Command{
		Name: "git",
		Desc: "git usage",
		Subs: []*gcli.Command{
			// l1: sub command 1
			{
				Name: "add",
				Desc: "the add command for git",
				Config: func(c *gcli.Command) {
					c.AddArg("files", "added files", true)
				},
				Func: func(c *gcli.Command, args []string) error {
					bf.WriteString("command path: " + c.Path())
					dump.Println(c.Name, args)
					return nil
				},
			},
			// l1: sub command 2
			{
				Name:    "pull",
				Desc:    "the pull command for git",
				Aliases: []string{"pul"},
				Func: func(c *gcli.Command, args []string) error {
					bf.WriteString("command path: " + c.Path())
					dump.Println(c.Name, args)
					return nil
				},
			},
			// l1: sub command 3
			{
				Name:    "remote",
				Desc:    "remote command for git",
				Aliases: []string{"rmt"},
				Func: func(c *gcli.Command, args []string) error {
					dump.Println(c.Path())
					return nil
				},
				Subs: []*gcli.Command{
					// l2: sub command 4
					{
						Name: "add",
						Desc: "add command for git remote",
						Config: func(c *gcli.Command) {
							c.AddArg("name", "the remote name", true)
							c.AddArg("address", "the remote address", true)
						},
						Func: func(c *gcli.Command, args []string) error {
							bf.WriteString("command path: " + c.Path())
							dump.Println(c.Path(), args)
							return nil
						},
					},
					// l2: sub command 5
					{
						Name:    "set-url",
						Desc:    "set-url command for git remote",
						Aliases: []string{"su"},
						Func: func(c *gcli.Command, args []string) error {
							bf.WriteString("command path: " + c.Path())
							dump.Println(c.Path(), args)
							return nil
						},
					},
				},
			},
		},
		Func: func(c *gcli.Command, args []string) error {
			bf.WriteString("command path: " + c.Path())
			// dump.Println(c.Path(), args)
			return nil
		},
	}

	return r, bf
}

func TestCommand_MatchByPath(t *testing.T) {
	r, _ := newGitCmd()
	c := r.MatchByPath("add")

	assert.NotNil(t, c)
	assert.Eq(t, "add", c.Name)
	assert.Eq(t, "git", c.ParentName())

	c = r.MatchByPath("remote:add")
	assert.NotNil(t, c)
	assert.Eq(t, "add", c.Name)
	assert.Eq(t, "Add command for git remote", c.Desc)
	assert.Eq(t, "remote", c.Parent().Name)
	assert.Eq(t, "git", c.Root().Name)

	// empty will return self
	c = r.MatchByPath("")
	assert.NotNil(t, c)
	assert.Eq(t, "git", c.Name)

	c = r.MatchByPath("not-exist")
	assert.Nil(t, c)
}

func TestCommand_Sub(t *testing.T) {
	r, _ := newGitCmd()
	r.MatchByPath("") // use for init

	assert.True(t, r.IsRoot())
	assert.True(t, r.IsCommand("remote"))
	assert.True(t, r.IsCommand("remote"))

	c := r.Sub("add")
	assert.NotNil(t, c)
	assert.False(t, c.IsRoot())
	assert.Eq(t, "add", c.Name)
}

func TestCommand_Run_top(t *testing.T) {
	r, bf := newGitCmd()

	err := r.Run([]string{})
	assert.NoErr(t, err)
	assert.Eq(t, "command path: git", bf.String())
}

func TestCommand_Run_oneLevelSub(t *testing.T) {
	r, _ := newGitCmd()

	err := r.Run([]string{"add", "./"})
	assert.NoErr(t, err)
}

func TestCommand_Run_moreLevelSub(t *testing.T) {
	r, bf := newGitCmd()
	err := r.Run([]string{
		"remote",
		"add",
		"origin",
		"https://github.com/inhere/console",
	})

	assert.NoErr(t, err)
	assert.True(t, r.IsAlias("rmt"))
	assert.True(t, r.IsAlias("pul"))
	assert.False(t, r.IsAlias("not-exist"))
	assert.Eq(t, "remote", r.ResolveAlias("rmt"))
	assert.Eq(t, "command path: git remote add", bf.String())
}

// c0Opts 保存 newC0Cmd 绑定的选项值, 每个用例独立持有, 避免共享可变状态。
type c0Opts struct {
	int0 int
	str0 string
}

// newC0Cmd 构造一个与原共享 c0 等价的命令, 每个用例独立调用获得全新实例,
// 返回命令本身、其写入缓冲 bf 以及绑定的选项值容器。
func newC0Cmd() (*gcli.Command, *bytes.Buffer, *c0Opts) {
	bf := new(bytes.Buffer)
	opts := &c0Opts{}

	c0 := gcli.NewCommand("test", "desc for test command", func(c *gcli.Command) {
		c.IntOpt(&opts.int0, "int", "", 0, "int desc")
		c.StrOpt(&opts.str0, "str", "", "", "str desc")
		c.AddArg("arg0", "arg0 desc")
		c.AddArg("arg1", "arg1 desc")
		c.Func = func(c *gcli.Command, args []string) error {
			bf.WriteString("name=" + c.Name)
			c.Ctx.Set("name", c.Name)
			c.Ctx.Set("args", args)
			// dump.P(c.ID(), "command Func is exec")
			return nil
		}
	})

	return c0, bf, opts
}

func TestCommand_Run_emptyArgs(t *testing.T) {
	is := assert.New(t)
	c0, bf, _ := newC0Cmd()

	defer gcli.ResetGOpts()
	gcli.SetVerbose(gcli.VerbCrazy)
	defer gcli.ResetVerbose()

	is.Eq("test", c0.Name)

	err := c0.Run([]string{})

	is.NoErr(err)
	is.Eq("name=test", bf.String())
	is.Eq("int desc", c0.Opt("int").Desc)
	is.NotEmpty(c0.Args())
	is.Eq("arg0", c0.Arg("arg0").Name)
}

func TestCommand_Run_showHelp1(t *testing.T) {
	is := assert.New(t)
	c0, _, _ := newC0Cmd()

	// SetDisable 会改包级 gOpts, 跑完后重置, 避免污染其它用例(如 --version 选项)
	defer gcli.ResetGOpts()
	gcli.Config(func(opts *gcli.GlobalOpts) {
		opts.SetDisable()
	})
	err := c0.Run([]string{"-h"})
	is.NoErr(err)
}

func TestCommand_Run_showHelp2(t *testing.T) {
	is := assert.New(t)
	c0, bf, _ := newC0Cmd()
	defer gcli.ResetGOpts()

	// no color
	color.Disable()
	color.SetOutput(bf)
	defer color.ResetOptions()

	err := c0.Run([]string{"--help"})
	is.NoErr(err)
	str := bf.String()
	is.Contains(str, "Int desc")
	is.Contains(str, "--str string")
	is.Contains(str, "Str desc")
	is.Contains(str, "Display the help information")
	is.StrContains(str, "Arg0 desc")
	is.StrContains(str, "Arg1 desc")
}

func TestCommand_Run_parseOptions(t *testing.T) {
	is := assert.New(t)
	c0, _, opts := newC0Cmd()

	defer gcli.ResetGOpts()
	gcli.SetDebugMode()
	defer gcli.ResetVerbose()

	is.Eq("test", c0.Name)

	dump.P(gcli.GOpts())
	err := c0.Run([]string{"--int", "10", "--str=abc", "txt"})

	// dump.P(gcli.GOpts(), c0.Context)
	is.NoErr(err)
	is.Eq("test", c0.Ctx.Get("name"))
	is.Eq("txt", c0.Arg("arg0").String())
	is.Empty(c0.Ctx.Get("args"))

	is.Eq(10, opts.int0)
	is.Eq("abc", opts.str0)
	is.Eq([]string{"txt"}, c0.FSetArgs())
	is.Eq("txt", c0.RawArg(0))

	// var str0 string
	co := struct {
		maxSteps  int
		overwrite bool
	}{}
	var int1 int

	c1 := gcli.NewCommand("test1", "desc test", func(c *gcli.Command) {
		c.IntOpt(&int1, "int", "", 0, "desc")
		c.IntOpt(&co.maxSteps, "max-step", "", 0, "setting the max step value")
		c.AddArg("arg0", "arg0 desc")
	}).WithFunc(func(c *gcli.Command, args []string) error {
		is.Eq("txt", c.Arg("arg0").String())
		is.Empty(args)
		return nil
	})

	is.Eq("test1", c1.Name)
	err = c1.Run([]string{"--int", "10", "--max-step=100", "txt"})
	is.NoErr(err)
	is.Eq(10, int1)
	is.Eq(100, co.maxSteps)
	is.Eq("[txt]", fmt.Sprint(c0.RawArgs()))
}

func TestInts(t *testing.T) {
	is := assert.New(t)
	ints := gcli.Ints{}

	err := ints.Set("1")
	is.NoErr(err)
	err = ints.Set("3")
	is.NoErr(err)
	is.Eq("[1,3]", ints.String())
	err = ints.Set("abc")
	is.Err(err)

	ints = gcli.Ints{1, 3}
	is.Eq("[1,3]", ints.String())
}

func TestStrings(t *testing.T) {
	is := assert.New(t)
	ss := gcli.Strings{}

	err := ss.Set("1")
	is.NoErr(err)
	err = ss.Set("3")
	is.NoErr(err)
	err = ss.Set("abc")
	is.NoErr(err)
	is.Eq("1,3,abc", ss.String())
}

func TestBooleans(t *testing.T) {
	is := assert.New(t)
	val := gcli.Booleans{}

	err := val.Set("false")
	is.NoErr(err)
	is.False(val[0])
	is.Eq("[false]", val.String())

	err = val.Set("True")
	is.NoErr(err)
	is.Eq("[false,true]", val.String())

	err = val.Set("abc")
	is.Err(err)
}

func TestCommand_Copy(t *testing.T) {
	is := assert.New(t)

	c := gcli.NewCommand("orig", "desc")
	c.On(gcli.EvtCmdInit, func(ctx *gcli.HookCtx) bool { return false })
	is.True(c.HasHook(gcli.EvtCmdInit))

	nc := c.Copy()
	// 副本不应继承原命令的钩子
	is.False(nc.HasHook(gcli.EvtCmdInit))
	// 关键回归点：Copy 不能清空原命令的钩子
	is.True(c.HasHook(gcli.EvtCmdInit))
}

func TestCommand_Use(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		defer gcli.ResetGOpts()

		var order []string
		c := gcli.NewCommand("mw-order", "test middleware order")
		ret := c.Use(
			func(c *gcli.Command, args []string) error {
				order = append(order, "mw1")
				return nil
			},
			func(c *gcli.Command, args []string) error {
				order = append(order, "mw2")
				return nil
			},
		)
		// 链式：Use 返回 *Command
		assert.Eq(t, c, ret)

		c.SetFunc(func(c *gcli.Command, args []string) error {
			order = append(order, "func")
			return nil
		})

		err := c.Run([]string{})
		assert.NoErr(t, err)
		// 按注册顺序: mw1 -> mw2 -> func
		assert.Eq(t, []string{"mw1", "mw2", "func"}, order)
	})

	t.Run("abort on error", func(t *testing.T) {
		defer gcli.ResetGOpts()

		var order []string
		wantErr := fmt.Errorf("mw1 abort")

		c := gcli.NewCommand("mw-abort", "test middleware abort")
		c.Use(
			func(c *gcli.Command, args []string) error {
				order = append(order, "mw1")
				return wantErr
			},
			func(c *gcli.Command, args []string) error {
				order = append(order, "mw2")
				return nil
			},
		)
		c.SetFunc(func(c *gcli.Command, args []string) error {
			order = append(order, "func")
			return nil
		})

		err := c.Run([]string{})
		// 第一个中间件返回 error, 错误向上返回
		assert.Err(t, err)
		assert.Eq(t, wantErr, err)
		// 第二个中间件与主函数都不应执行
		assert.Eq(t, []string{"mw1"}, order)
	})

	t.Run("no middleware regression", func(t *testing.T) {
		defer gcli.ResetGOpts()

		called := false
		c := gcli.NewCommand("mw-none", "test no middleware")
		c.SetFunc(func(c *gcli.Command, args []string) error {
			called = true
			return nil
		})

		err := c.Run([]string{})
		assert.NoErr(t, err)
		assert.True(t, called)
	})
}
