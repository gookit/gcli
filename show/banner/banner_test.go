package banner_test

import (
	"fmt"
	"testing"

	"github.com/gookit/gcli/v3/show/banner"
)

func TestBanner_Render(t *testing.T) {
	b1 := banner.NewBanner1("Hello World")
	text := b1.Render()
	fmt.Println(text)

	// multi lines
	b1.Contents = []string{"Hello", "World"}
	fmt.Println(b1.Render())

	// change border style
	b1.Contents = []string{"Hello, World"}
	b1.BorderStyle = banner.SimpleBorderStyle
	fmt.Println(b1.Render())
}

func TestBanner_Render_cn(t *testing.T) {
	// 使用中文内容
	b1 := banner.NewBanner1("你好，世界")
	b1.BorderStyle = banner.SimpleBorderStyle
	fmt.Println(b1.Render())

	// 使用 中文+英文混合
	b1.Contents = []string{"你好，World"}
	b1.BorderStyle = banner.SharpBorderStyle
	fmt.Println(b1.Render())
}
