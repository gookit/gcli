package banner_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3/show/banner"
	"github.com/gookit/gcli/v3/show/showcom"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/envutil"
	"github.com/gookit/goutil/testutil"
	"github.com/gookit/goutil/testutil/assert"
	"github.com/gookit/goutil/x/termenv"
)

func TestBanner_Render(t *testing.T) {
	b1 := banner.New("Hello World")
	text := b1.Render()
	fmt.Println(text)

	// multi lines
	b1.Contents = []string{"Hello", "World"}
	fmt.Println(b1.Render())

	// change border style
	b1.WithContents([]any{"Hello", 123})
	b1.BorderStyle = banner.SimpleBorderStyle
	fmt.Println(b1.Render())

	t.Run("Other type", func(t *testing.T) {
		b := banner.New(123456)
		fmt.Println(b.Render())
	})
}

func TestBanner_Render_cn(t *testing.T) {
	// 使用中文内容
	b1 := banner.New("你好，世界")
	b1.BorderStyle = banner.SimpleBorderStyle
	fmt.Println(b1.Render())

	// 使用 中文+英文混合
	b1.Contents = []string{"你好，World"}
	b1.BorderStyle = banner.SharpBorderStyle
	fmt.Println(b1.Render())
}

// go test -v -run ^TestBanner_WithNewOptions$ ./show/banner/...
func TestBanner_WithNewOptions(t *testing.T) {
	// 测试 Height 选项
	t.Run("Height gt lines", func(t *testing.T) {
		b1 := banner.New([]string{"Line 1", "Line 2"}, banner.WithHeight(3))
		result1 := b1.Render()
		assert.NotEmpty(t, result1)
		fmt.Println(result1)
	})
	t.Run("Height lt lines", func(t *testing.T) {
		b1 := banner.New([]string{"Line 1", "Line 2", "Line 3"}).WithOptionFn(banner.WithHeight(2))
		result1 := b1.Render()
		assert.NotEmpty(t, result1)
		fmt.Println(result1)
	})

	// 测试 MinWidth 选项
	t.Run("MinWidth", func(t *testing.T) {
		b2 := banner.New([]string{"MinWidth"}, banner.WithMinWidth(30))
		result2 := b2.Render()
		assert.NotEmpty(t, result2)
		fmt.Println(result2)
	})

	// 测试 OverflowFlag 选项
	t.Run("OverflowFlag wrap", func(t *testing.T) {
		b := banner.New([]string{"This is a very long line that should be wrapped"}).
			WithOptionFn(banner.WithWidth(20), banner.WithOverflowFlag(showcom.OverflowWrap))
		result := b.Render()
		assert.NotEmpty(t, result)
		fmt.Println(result)
	})
	t.Run("OverflowFlag cut", func(t *testing.T) {
		b := banner.New([]string{"This is a very long line that should be truncated"}).
			WithOptionFn(banner.WithWidth(20), banner.WithOverflowFlag(showcom.OverflowCut))
		result := b.Render()
		assert.NotEmpty(t, result)
		fmt.Println(result)
	})

	// 测试 Alignment 选项
	t.Run("Alignment center", func(t *testing.T) {
		b := banner.New([]string{"Centered Text"}, banner.WithWidth(30), banner.WithAlignment(comdef.Center))
		result := b.Render()
		assert.NotEmpty(t, result)
		fmt.Println(result)
	})

	t.Run("Alignment right", func(t *testing.T) {
		b := banner.New([]string{"Right Text"}, banner.WithWidth(30), banner.WithAlignment(comdef.Right))
		result := b.Render()
		assert.NotEmpty(t, result)
		fmt.Println(result)
	})
}

// go test -v -run ^TestBanner_MarginOption$ ./show/banner/...
func TestBanner_MarginOption(t *testing.T) {
	// 测试 Margin 选项
	t.Run("Margin left", func(t *testing.T) {
		b := banner.New([]string{"Margin left"}, banner.WithMarginLeft(5))
		result := b.Render()
		assert.NotEmpty(t, result)
		fmt.Println(result)
	})

	// TIP: return is not a terminal on run go test
	fmt.Print("IsTerminal: ")
	fmt.Println(termenv.IsTerminal())

	// use env mock
	testutil.MockEnvValues(map[string]string{"COLUMNS": "80", "LINES": "24"}, func() {
		fmt.Print("MOCK ENV ")
		fmt.Println(envutil.GetMulti("COLUMNS", "LINES"))
		// fmt.Println(termenv.GetTermSize())
		// 测试 居中对齐
		t.Run("Margin center", func(t *testing.T) {
			b := banner.New([]string{"Centered Banner"}, banner.WithBannerCenter())
			result := b.Render()
			assert.NotEmpty(t, result)
			assert.StrContains(t, result, "                              │ Centered Banner │")
			fmt.Println(result)
		})

		// 测试居右
		t.Run("Margin Right", func(t *testing.T) {
			b := banner.New([]string{"Right Banner"}, banner.WithBannerRight())
			result := b.Render()
			assert.NotEmpty(t, result)
			assert.StrContains(t, result, strings.Repeat(" ", 64)+"│ Right Banner │")
			fmt.Println(result)
		})

		// 测试 PercentWidth 选项
		t.Run("PercentWidth", func(t *testing.T) {
			b3 := banner.New([]string{"Test PercentWidth"}, banner.WithPercentWidth(50))
			result3 := b3.Render()
			assert.NotEmpty(t, result3)
			assert.StrContains(t, result3, "│ Test PercentWidth                        │")
			fmt.Println(result3)
		})
	})

}
