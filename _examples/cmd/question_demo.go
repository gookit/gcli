package cmd

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/gcli/v3/gflag"
	"github.com/gookit/goutil/dump"
)

var questionOpts = struct {
	token string
	name  string
}{}

// QuestionDemo 演示 B7：CliOpt.Question 声明式交互收集
var QuestionDemo = &gcli.Command{
	Name:    "ask-demo",
	Aliases: []string{"qd"},
	Desc:    "demo B7: declarative interactive collect by CliOpt.Question",
	Examples: `
  <cyan>{$fullCmd}</>                  # 不带 --token 运行 -> 自动提问收集输入
  <cyan>{$fullCmd} --token abc123</>   # 已提供值则不提问`,
	Config: func(c *gcli.Command) {
		// WithQuestion: 当选项值为空时，自动用该问题交互收集输入(内置默认 Collector)。
		// 若同时设置了 WithCollector，则 Collector 优先、Question 被忽略。
		c.StrOpt2(&questionOpts.token, "token", "the access token (will ask if empty)",
			gflag.WithQuestion("Please input your access token: "))
		c.StrOpt(&questionOpts.name, "name", "n", "guest", "the user name")
	},
	Func: func(c *gcli.Command, args []string) error {
		color.Infoln("hello, in ask-demo command (B7 demo)")
		color.Cyanln("Collected options:")
		dump.V(questionOpts)
		return nil
	},
}
