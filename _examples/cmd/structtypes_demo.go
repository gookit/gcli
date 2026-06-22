package cmd

import (
	"time"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
)

// structTypesOpts 演示 D1：结构体标签绑定更丰富的字段类型。
// 原生 slice / time.Duration / map[string]string 直接可绑，enum 标签做候选+校验。
type structTypesOpts struct {
	Names []string          `flag:"name=names;shorts=n;desc=name list (repeatable)"`
	Ports []int             `flag:"name=ports;shorts=p;desc=port list (repeatable)"`
	TTL   time.Duration     `flag:"name=ttl;desc=time to live, eg: 1h30m"`
	Meta  map[string]string `flag:"name=meta;shorts=m;desc=key=value metadata (repeatable)"`
	Lang  string            `flag:"name=lang;shorts=l;desc=language;enum=go,php,java"`
}

var structTypesData = &structTypesOpts{}

// StructTypesDemo 演示 D1：结构体标签支持 slice / Duration / map / enum
var StructTypesDemo = &gcli.Command{
	Name:    "struct-types",
	Aliases: []string{"stp"},
	Desc:    "demo D1: bind richer field types from struct (slice/duration/map/enum)",
	Examples: `
  <cyan>{$fullCmd} -n a -n b -p 80 -p 443 --ttl 1h30m -m k1=v1 -m k2=v2 -l go</>
  <cyan>{$fullCmd} -l ruby</>   # enum 校验失败示例(lang 仅允许 go/php/java)`,
	Config: func(c *gcli.Command) {
		// 默认 TagRuleNamed：flag tag 内 name/shorts/desc/enum 等分号分隔。
		c.MustFromStruct(structTypesData)
	},
	Func: func(c *gcli.Command, args []string) error {
		color.Infoln("hello, in struct-types command (D1 demo)")
		color.Cyanln("Parsed struct options:")
		dump.V(structTypesData)
		color.Cyanf("TTL(readable): %s\n", structTypesData.TTL.String())
		return nil
	},
}
