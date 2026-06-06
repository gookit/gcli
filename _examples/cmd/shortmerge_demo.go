package cmd

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

var shortMergeOpts = struct {
	all     bool
	upload  bool
	extract bool
	output  string
}{}

// ShortMergeDemo 演示 B4+B5：EnhanceShort POSIX 短选项合并
var ShortMergeDemo = &gcli.Command{
	Name:    "short-merge",
	Aliases: []string{"smd"},
	Desc:    "demo B4+B5: POSIX short option merge by EnhanceShort",
	Examples: `
  <cyan>{$fullCmd} -aux</>            # = -a -u -x (三个 bool 短选项合并，仅全 bool 才拆)
  <cyan>{$fullCmd} -Ostdout</>        # = -O stdout (取值短选项紧贴写法，需 level2)
  <cyan>{$fullCmd} -au -O file.txt</> # 混合写法`,
	Config: func(c *gcli.Command) {
		// 开启 EnhanceShortAttach(level2)：
		//   Merge(level1)  -> 全为 bool 的短选项组合拆分(-aux => -a -u -x)
		//   Attach(level2) -> 额外支持取值短选项紧贴写法(-Ostdout => -O stdout)
		c.ParserCfg().EnhanceShort = gcli.EnhanceShortAttach

		c.BoolOpt(&shortMergeOpts.all, "all", "a", false, "the all option (bool)")
		c.BoolOpt(&shortMergeOpts.upload, "upload", "u", false, "the upload option (bool)")
		c.BoolOpt(&shortMergeOpts.extract, "extract", "x", false, "the extract option (bool)")
		c.StrOpt(&shortMergeOpts.output, "output", "O", "", "the output option (take value)")
	},
	Func: func(c *gcli.Command, args []string) error {
		color.Infoln("hello, in short-merge command (B4+B5 demo)")
		color.Cyanln("Parsed options:")
		dump.V(shortMergeOpts)
		return nil
	},
}
