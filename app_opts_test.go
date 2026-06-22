package gcli_test

import (
	"testing"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/byteutil"
	"github.com/gookit/goutil/x/assert"
)

// D2.6: 每-App 解析状态(AppOptions)相互独立, 全局配置(GlobalOpts)仍共享。
func TestApp_AppOpts_perAppIsolation(t *testing.T) {
	is := assert.New(t)

	mk := func() *gcli.App {
		a := gcli.NewApp(gcli.NotExitOnEnd())
		a.Version = "1.0.0"
		a.Add(&gcli.Command{
			Name: "demo", Desc: "demo command",
			Func: func(c *gcli.Command, _ []string) error { return nil },
		})
		return a
	}

	app1, app2 := mk(), mk()

	// 每-App 的 AppOptions 是独立实例
	is.True(app1.AppOpts() != app2.AppOpts())
	// 全局 GlobalOpts 仍是共享的包级单例(向后兼容: app.Opts()==GOpts())
	is.True(app1.Opts() == app2.Opts())
	is.True(app1.Opts() == gcli.GOpts())

	// 抑制 --version 横幅输出
	b := byteutil.NewBuffer()
	color.Disable()
	color.SetOutput(b)
	defer color.ResetOptions()

	// app1 解析 --version 写的是 app1 自己的状态, 不污染 app2
	app1.Run([]string{"--version"})
	is.True(app1.AppOpts().ShowVersion)
	is.False(app2.AppOpts().ShowVersion)
}
