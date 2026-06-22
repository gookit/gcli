package gcli_test

import (
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/x/assert"
)

// 继承: top 定义共享 --git-dir, 在子命令 sub 段写出时能被解析。
func TestSharedOpts_inherit(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var gitDir string
	var ran bool
	app.Add(&gcli.Command{
		Name: "top",
		Desc: "top command",
		Config: func(c *gcli.Command) {
			c.SharedOpts().StrOpt(&gitDir, "git-dir", "", "", "the git work dir")
		},
		Subs: []*gcli.Command{
			{
				Name: "sub",
				Desc: "sub command",
				Func: func(c *gcli.Command, _ []string) error {
					ran = true
					return nil
				},
			},
		},
	})

	code := app.Run([]string{"top", "sub", "--git-dir", "/x"})
	is.Eq(0, code)
	is.True(ran)
	is.Eq("/x", gitDir)
}

// 任意位置(配合重排): 共享选项写在 arguments 之后也能解析, arg 仍正常。
func TestSharedOpts_anyPosition(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var gitDir string
	var a0 string
	app.Add(&gcli.Command{
		Name: "top",
		Desc: "top command",
		Config: func(c *gcli.Command) {
			c.SharedOpts().StrOpt(&gitDir, "git-dir", "", "", "the git work dir")
		},
		Subs: []*gcli.Command{
			{
				Name: "sub",
				Desc: "sub command",
				Config: func(c *gcli.Command) {
					c.AddArg("arg0", "arg0 desc")
				},
				Func: func(c *gcli.Command, _ []string) error {
					a0 = c.Arg("arg0").String()
					return nil
				},
			},
		},
	})

	code := app.Run([]string{"top", "sub", "myarg", "--git-dir", "/x"})
	is.Eq(0, code)
	is.Eq("/x", gitDir)
	is.Eq("myarg", a0)
}

// 多级继承: 三层 a -> b -> c, a 的共享选项在叶子 c 可用。
func TestSharedOpts_multiLevel(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var aShared string
	var ran bool
	app.Add(&gcli.Command{
		Name: "a",
		Desc: "a command",
		Config: func(c *gcli.Command) {
			c.SharedOpts().StrOpt(&aShared, "a-opt", "", "", "shared opt from a")
		},
		Subs: []*gcli.Command{
			{
				Name: "b",
				Desc: "b command",
				Subs: []*gcli.Command{
					{
						Name: "c",
						Desc: "c command",
						Func: func(c *gcli.Command, _ []string) error {
							ran = true
							return nil
						},
					},
				},
			},
		},
	})

	code := app.Run([]string{"a", "b", "c", "--a-opt", "vvv"})
	is.Eq(0, code)
	is.True(ran)
	is.Eq("vvv", aShared)
}

// 局部优先: 子命令定义同名局部选项时, 继承被跳过, 父子变量互不写串。
func TestSharedOpts_localPriority(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var shared string // top 的共享变量
	var local string  // sub 的局部同名变量
	app.Add(&gcli.Command{
		Name: "top",
		Desc: "top command",
		Config: func(c *gcli.Command) {
			c.SharedOpts().StrOpt(&shared, "name", "", "", "shared name")
		},
		Subs: []*gcli.Command{
			{
				Name: "sub",
				Desc: "sub command",
				Config: func(c *gcli.Command) {
					// 同名局部选项, 应优先于继承
					c.StrOpt(&local, "name", "", "", "local name")
				},
				Func: func(c *gcli.Command, _ []string) error { return nil },
			},
		},
	})

	code := app.Run([]string{"top", "sub", "--name", "L"})
	is.Eq(0, code)
	// 写回的是局部变量, 共享变量不受影响
	is.Eq("L", local)
	is.Eq("", shared)
}

// 自身可用: 共享选项写在子命令名之前, 由定义它的命令段即可识别。
func TestSharedOpts_selfUsable(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var gitDir string
	var ran bool
	app.Add(&gcli.Command{
		Name: "top",
		Desc: "top command",
		Config: func(c *gcli.Command) {
			c.SharedOpts().StrOpt(&gitDir, "git-dir", "", "", "the git work dir")
		},
		Subs: []*gcli.Command{
			{
				Name: "sub",
				Desc: "sub command",
				Func: func(c *gcli.Command, _ []string) error {
					ran = true
					return nil
				},
			},
		},
	})

	code := app.Run([]string{"top", "--git-dir", "/x", "sub"})
	is.Eq(0, code)
	is.True(ran)
	is.Eq("/x", gitDir)
}

// buildRequiredApp 构造一个 top 定义 required 共享选项、sub 在 Func 标记运行的 app。
// 用独立 app 实例避免「同一命令重复 Parse」的非典型场景。
func buildRequiredApp(gitDir *string, ran *bool) *gcli.App {
	app := gcli.NewApp(gcli.NotExitOnEnd())
	app.Add(&gcli.Command{
		Name: "top",
		Desc: "top command",
		Config: func(c *gcli.Command) {
			c.SharedOpts().StrOpt2(gitDir, "git-dir", "the git work dir", gflag.WithRequired())
		},
		Subs: []*gcli.Command{
			{
				Name: "sub",
				Desc: "sub command",
				Func: func(c *gcli.Command, _ []string) error {
					*ran = true
					return nil
				},
			},
		},
	})
	return app
}

// Required 共享: 共享选项设 Required, 缺失时执行命令报错, Func 不应运行; 提供后正常执行。
func TestSharedOpts_required(t *testing.T) {
	is := assert.New(t)

	t.Run("missing required", func(t *testing.T) {
		var gitDir string
		var ran bool
		app := buildRequiredApp(&gitDir, &ran)
		// 缺失 required 共享选项: parseOptions 校验失败, sub.Func 不执行
		app.Run([]string{"top", "sub"})
		is.False(ran)
	})

	t.Run("provided required", func(t *testing.T) {
		var gitDir string
		var ran bool
		app := buildRequiredApp(&gitDir, &ran)
		code := app.Run([]string{"top", "sub", "--git-dir", "/x"})
		is.Eq(0, code)
		is.True(ran)
		is.Eq("/x", gitDir)
	})
}

// Required 共享(int 类型): 校验需类型感知判空 —— 缺失时值为 "0" 也应判为未提供。
func TestSharedOpts_requiredInt(t *testing.T) {
	is := assert.New(t)

	build := func(count *int, ran *bool) *gcli.App {
		app := gcli.NewApp(gcli.NotExitOnEnd())
		app.Add(&gcli.Command{
			Name: "top",
			Desc: "top command",
			Config: func(c *gcli.Command) {
				c.SharedOpts().IntOpt(count, "count", "", 0, "the count", gflag.WithRequired())
			},
			Subs: []*gcli.Command{
				{
					Name: "sub",
					Desc: "sub command",
					Func: func(c *gcli.Command, _ []string) error { *ran = true; return nil },
				},
			},
		})
		return app
	}

	t.Run("missing int required", func(t *testing.T) {
		var count int
		var ran bool
		build(&count, &ran).Run([]string{"top", "sub"})
		is.False(ran) // 缺失(值=0)应判为未提供 -> Func 不执行
	})

	t.Run("provided int required", func(t *testing.T) {
		var count int
		var ran bool
		code := build(&count, &ran).Run([]string{"top", "sub", "--count", "3"})
		is.Eq(0, code)
		is.True(ran)
		is.Eq(3, count)
	})
}

// help 分组(D2.5): 继承选项在命令 help 里归到 "Inherited Options" 分组, 与本地选项分开。
func TestSharedOpts_helpGrouping(t *testing.T) {
	is := assert.New(t)

	newApp := func() *gcli.App {
		app := gcli.NewApp(gcli.NotExitOnEnd())
		var gitDir, local string
		app.Add(&gcli.Command{
			Name: "top",
			Desc: "top command",
			Config: func(c *gcli.Command) {
				c.SharedOpts().StrOpt(&gitDir, "git-dir", "", "", "the git work dir")
			},
			Subs: []*gcli.Command{
				{
					Name: "sub",
					Desc: "sub command",
					Config: func(c *gcli.Command) {
						c.StrOpt(&local, "name", "n", "", "the local name")
					},
					Func: func(c *gcli.Command, _ []string) error { return nil },
				},
			},
		})
		return app
	}

	b := byteutil.NewBuffer()
	color.Disable()
	color.SetOutput(b)
	defer color.ResetOptions()

	// 叶子 sub 的 help: top 的共享选项是「祖先继承」, 应归入 "Inherited Options" 分组
	t.Run("inherited grouped on leaf (cmd -h)", func(t *testing.T) {
		b.Reset()
		newApp().Run([]string{"top", "sub", "-h"})
		s := b.String()
		is.StrContains(s, "Inherited Options") // 继承分组标题
		is.StrContains(s, "--git-dir")         // 祖先继承选项
		is.StrContains(s, "--name")            // 本地选项
	})

	// help 命令(单级)展示 top 自身 help: 自身共享选项应显示(验证 ShowHelp 也会合并),
	// 但作为命令自身选项, 不归入 "Inherited Options" 分组。
	t.Run("self shared shown via help command", func(t *testing.T) {
		b.Reset()
		newApp().Run([]string{"help", "top"})
		s := b.String()
		is.StrContains(s, "--git-dir")            // ShowHelp 合并后自身共享选项可见
		is.NotContains(s, "Inherited Options")    // 自身定义的不标"Inherited"
	})
}

// 未用共享回归: 不定义共享选项时, 行为与之前完全一致(sharedFs==nil, 合并是 no-op)。
func TestSharedOpts_noSharedRegression(t *testing.T) {
	is := assert.New(t)
	app := gcli.NewApp(gcli.NotExitOnEnd())

	var name string
	var a0 string
	app.Add(&gcli.Command{
		Name: "top",
		Desc: "top command",
		Subs: []*gcli.Command{
			{
				Name: "sub",
				Desc: "sub command",
				Config: func(c *gcli.Command) {
					c.StrOpt(&name, "name", "n", "", "name opt")
					c.AddArg("arg0", "arg0 desc")
				},
				Func: func(c *gcli.Command, _ []string) error {
					a0 = c.Arg("arg0").String()
					return nil
				},
			},
		},
	})

	code := app.Run([]string{"top", "sub", "av", "--name", "tom"})
	is.Eq(0, code)
	is.Eq("tom", name)
	is.Eq("av", a0)
}
