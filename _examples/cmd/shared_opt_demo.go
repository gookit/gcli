package cmd

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
)

var sharedAskOpts = struct {
	token string
}{}

// SharedAskDemo 演示共享选项的交互收集(CliOpt.Question)。
//
// 共享选项的值可以写在任意命令段，中间命令解析时看不到叶子段的取值，因此收集延后到
// 实际执行的命令处只做一次(见 Command.prepareSharedOpts)：命令行任意段给了值都不提问。
// 对应 PTY 冒烟用例: scripts/ptysmoke/shared_ask.py
var SharedAskDemo = &gcli.Command{
	Name: "shared-ask",
	Desc: "demo: interactive collect of a shared option, asked once at the executed command",
	Examples: `
  <cyan>{$fullCmd} run</>                  # 不带 --token -> 在执行命令处提问一次
  <cyan>{$fullCmd} run --token abc123</>   # 叶子段给值 -> 不提问
  <cyan>{$fullCmd} --token abc123 run</>   # 父命令段给值 -> 不提问`,
	Config: func(c *gcli.Command) {
		c.SharedOpts().StrOpt2(&sharedAskOpts.token, "token", "the api token (asked when empty)",
			gflag.WithQuestion("Please input your api token: "))
	},
	Subs: []*gcli.Command{
		{
			Name: "run",
			Desc: "leaf command inheriting the shared --token",
			Func: func(c *gcli.Command, _ []string) error {
				color.Cyanln("leaf: token=" + sharedAskOpts.token)
				return nil
			},
		},
	},
}
