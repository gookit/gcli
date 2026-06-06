package cmd

import (
	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

// baseFlags 匿名嵌套结构体，演示 B6 的匿名字段展开：
// 它的导出字段会被并入外层命令的选项中。
type baseFlags struct {
	// 选项名取字段名 SnakeCase = "verbose"，short = "v"
	Verbose bool `flag:"v" desc:"enable verbose output"`
}

// structFlagOpts 演示 B6 TagRuleField 标签规则：
// 以字段名(SnakeCase)做选项名，元数据从独立 tag 键(flag/desc/default/required)读取。
type structFlagOpts struct {
	baseFlags // 匿名嵌套 -> 自动展开 verbose 选项

	UserName string `flag:"u" desc:"the user name" required:"true"`
	Age      int    `desc:"the user age" default:"18"`
	Email    string `flag:"e" desc:"the user email"`

	internal string // 未导出且无 tag -> 自动跳过，不会成为选项
}

var structFlagData = &structFlagOpts{}

// StructFlagDemo 演示 B6：TagRuleField 标签规则 + 匿名字段展开
var StructFlagDemo = &gcli.Command{
	Name:    "struct-flag",
	Aliases: []string{"sfd"},
	Desc:    "demo B6: bind options from struct by TagRuleField + anonymous field expand",
	Examples: `
  <cyan>{$fullCmd} --user-name tom --age 22 -e tom@example.com -v</>
  <cyan>{$fullCmd} -u tom</>   # age 用默认值 18，user-name 为必填项`,
	Config: func(c *gcli.Command) {
		// TagRuleField: 选项名=字段名 SnakeCase，shorts 来自 flag tag，
		// desc/default/required 各用独立 tag；匿名 baseFlags.Verbose 被展开为 -v 选项。
		c.MustFromStruct(structFlagData, gcli.TagRuleField)
	},
	Func: func(c *gcli.Command, args []string) error {
		color.Infoln("hello, in struct-flag command (B6 demo)")
		color.Cyanln("Parsed struct options:")
		dump.V(structFlagData)
		return nil
	},
}
