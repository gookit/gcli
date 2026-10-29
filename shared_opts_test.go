package gcli_test

import (
	"io"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/cliutil"
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

// 共享选项的 Collector/Validator/Handler 延后到执行命令时统一处理: 各只触发一次, 且用最终取值。
// 依据: 共享值可写在叶子命令段, 中间命令解析时看不到它(见 mergeSharedOpts/prepareSharedOpts)。
func TestSharedOpts_deferredCallbacks(t *testing.T) {
	is := assert.New(t)

	var collectCalls, validateCalls, handleCalls int
	var token, handled string

	newApp := func() *gcli.App {
		collectCalls, validateCalls, handleCalls = 0, 0, 0
		token, handled = "", ""

		app := gcli.NewApp(gcli.NotExitOnEnd())
		app.Add(&gcli.Command{
			Name: "top",
			Desc: "top command",
			Config: func(c *gcli.Command) {
				c.SharedOpts().StrOpt2(&token, "token", "the token",
					gflag.WithCollector(func() (string, error) { collectCalls++; return "collected", nil }),
					gflag.WithValidator(func(string) error { validateCalls++; return nil }),
					gflag.WithHandler(func(val string) error { handleCalls++; handled = val; return nil }),
				)
			},
			Subs: []*gcli.Command{
				{
					Name: "sub",
					Desc: "sub command",
					Func: func(c *gcli.Command, _ []string) error { return nil },
				},
			},
		})
		return app
	}

	t.Run("value from leaf segment: no prompt, each once", func(t *testing.T) {
		code := newApp().Run([]string{"top", "sub", "--token=abc"})
		is.Eq(0, code)
		is.Eq("abc", token)
		is.Eq(0, collectCalls)  // 已传值: 不应弹交互收集
		is.Eq(1, validateCalls) // 校验只跑一次(不再是父/子各一次)
		is.Eq(1, handleCalls)
		is.Eq("abc", handled)
	})

	t.Run("no value: collected once before running", func(t *testing.T) {
		code := newApp().Run([]string{"top", "sub"})
		is.Eq(0, code)
		is.Eq("collected", token)
		is.Eq(1, collectCalls)
		is.Eq(1, validateCalls)
		is.Eq(1, handleCalls)
		is.Eq("collected", handled)
	})
}

// askText is the question used by the injected-input case; it is counted in the
// captured prompt output.
const askText = "Please input the api token: "

// readCounter counts how often the injected stdin was consumed; a case that must
// not ask has to leave it at zero.
type readCounter struct {
	r     io.Reader
	reads int
}

func (c *readCounter) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.reads++
	return n, err
}

// 交互提问的取值路径用注入式输入测: CliOpt.Question 走 goutil/cliutil.ReadLine,
// 而 cliutil.Input/Output 可替换(gflag 层的 opts_test 已这么测单选项)。
// 真终端的语义(raw/cooked、回显、按键)由 scripts/ptysmoke 与 scripts/wincons 覆盖,
// 这里只锁“值从哪来、问几次、什么时候问”, 跨平台且无 PTY 依赖。
func TestSharedOpts_questionFromInjectedInput(t *testing.T) {
	is := assert.New(t)

	// 应用的 help/普通输出不需要进测试日志, 提示符输出单独用 cliutil.Output 捕获
	color.Disable()
	color.SetOutput(io.Discard)
	defer color.ResetOptions()

	type result struct {
		token  string
		prompt string
		reads  int
	}
	run := func(stdin string, args ...string) result {
		var token string
		in := &readCounter{r: strings.NewReader(stdin)}
		out := byteutil.NewBuffer()
		oldIn, oldOut := cliutil.Input, cliutil.Output
		cliutil.Input, cliutil.Output = in, out
		defer func() { cliutil.Input, cliutil.Output = oldIn, oldOut }()

		app := gcli.NewApp(gcli.NotExitOnEnd())
		app.Add(&gcli.Command{
			Name: "top",
			Config: func(c *gcli.Command) {
				c.SharedOpts().StrOpt2(&token, "token", "the token", gflag.WithQuestion(askText))
			},
			Subs: []*gcli.Command{
				{Name: "run", Func: func(*gcli.Command, []string) error { return nil }},
			},
		})
		is.Eq(0, app.Run(args))
		return result{token: token, prompt: out.String(), reads: in.reads}
	}

	t.Run("empty value: asked once at the executed command", func(t *testing.T) {
		res := run("injected-token\n", "top", "run")
		is.Eq("injected-token", res.token)
		is.Eq(1, strings.Count(res.prompt, askText)) // 不是父/子各问一次
		is.True(res.reads > 0)
	})

	t.Run("value in the leaf segment: stdin is never read", func(t *testing.T) {
		res := run("not-used\n", "top", "run", "--token=cli-token")
		is.Eq("cli-token", res.token)
		is.NotContains(res.prompt, askText)
		is.Eq(0, res.reads)
	})

	t.Run("help: stdin is never read", func(t *testing.T) {
		res := run("not-used\n", "top", "run", "-h")
		is.NotContains(res.prompt, askText)
		is.Eq(0, res.reads)
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
		is.StrContains(s, "--git-dir")         // ShowHelp 合并后自身共享选项可见
		is.NotContains(s, "Inherited Options") // 自身定义的不标"Inherited"
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
