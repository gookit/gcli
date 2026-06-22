package cmd

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

var reorderOpts = struct {
	name  string
	age   int
	force bool
}{}

// ReorderDemo 演示 args 自动重排：写在 arguments 之后的 options 仍能被正确解析。
//
// 该特性由 gflag Config.DisableReorderArgs 控制(默认开启)，
// 可用 c.ParserCfg().DisableReorderArgs = true 关闭以恢复严格顺序。
var ReorderDemo = &gcli.Command{
	Name:    "reorder-args",
	Aliases: []string{"rda"},
	Desc:    "demo: auto-reorder mixed input to canonical --options... arguments",
	Examples: `
  <cyan>{$fullCmd} src dst --name tom -f</>        # options 写在 arguments 之后, 仍被解析
  <cyan>{$fullCmd} src --name tom dst --age 18</>  # options/arguments 交错混写
  <cyan>{$fullCmd} --name tom -f src dst</>        # 标准顺序, 结果与上面一致`,
	Config: func(c *gcli.Command) {
		c.StrOpt(&reorderOpts.name, "name", "n", "", "the name option (take value)")
		c.IntOpt(&reorderOpts.age, "age", "a", 0, "the age option (take value)")
		c.BoolOpt(&reorderOpts.force, "force", "f", false, "the force option (bool)")

		c.AddArg("src", "the source argument", true) // required
		c.AddArg("dst", "the destination argument")
	},
	Func: func(c *gcli.Command, _ []string) error {
		color.Infoln("hello, in reorder-args command (auto-reorder demo)")

		color.Cyanln("Parsed options:")
		dump.V(reorderOpts)

		color.Cyanln("Parsed arguments:")
		dump.V(map[string]string{
			"src": c.Arg("src").String(),
			"dst": c.Arg("dst").String(),
		})
		return nil
	},
}
