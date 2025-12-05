package main

import (
	"fmt"
	"os"

	"github.com/gookit/gcli/v3/show/banner"
	"github.com/gookit/gcli/v3/show/showcom"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/x/termenv"
	"golang.org/x/term"
)

// RUN: go run ./_examples/show_banner.go
func main() {
	// 示例1：使用 Height 选项
	b1 := banner.New([]string{"Line 1", "Line 2"},
		banner.WithHeight(4))
	b1.Println()

	// 示例2：使用 MinWidth 选项
	b2 := banner.New([]string{"Short"},
		banner.WithMinWidth(30))
	b2.Println()

	// 示例3：使用 PercentWidth 选项
	b3 := banner.New([]string{"50% width banner"},
		banner.WithPercentWidth(50))
	b3.Println()

	// OverflowFlag 自动换行
	b4 := banner.New([]string{"This is a very long line that should wrap"})
	b4.Width = 18
	b4.OverflowFlag = showcom.OverflowWrap
	b4.Println()

	// 示例：使用 OverflowFlag 选项（截断）
	b4 = banner.New([]string{"This is a very long line that should be truncated"})
	b4.Width = 20
	b4.OverflowFlag = showcom.OverflowCut // 截断模式
	b4.Println()

	// 示例：使用 Alignment 选项（居中对齐）
	b5 := banner.New([]string{"Centered Text"},
		banner.WithWidth(30),
		banner.WithAlignment(comdef.Center))
	b5.Println()

	fmt.Println("IsTerminal:", termenv.IsTerminal())
	fmt.Println(termenv.GetTermSize())
	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		fmt.Printf("Get size error: %v\n", err)
	}
	fmt.Printf("Terminal size: %d x %d\n", width, height)

	// 示例6：使用 Margin 选项（左边距）
	b6 := banner.New([]string{"Indented Banner"}, banner.WithMarginLeft(10))
	b6.Padding = 2
	b6.Println()

	// 示例7：使用居中对齐
	b7 := banner.New([]string{"Centered Banner"}, banner.WithBannerCenter())
	b7.Println()

	// 示例8：使用右对齐
	b8 := banner.New([]string{"Right Aligned"}, banner.WithBannerRight())
	b8.Println()
}
