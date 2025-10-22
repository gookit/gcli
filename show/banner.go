package show

import (
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/goutil/strutil"
)

/*
eg: TODO
   ╭──────────────────────────────────────────────────────────────────╮
   │                                                                  │
   │                Update available! 3.21.0 → 3.27.0.                │
   │   Changelog: https://github.com/gookit/gcli/releases/tag/v3.2.0  │
   │                Run "x y z" to update.                			  │
   │                                                                  │
   ╰──────────────────────────────────────────────────────────────────╯
*/

/*
style1:
 ╭─────╮
 │  hi │
 ╰─────╯
style2:
  ┌────┐
  │ hi │
  └────┘
style3:
 ╔═════╗
 ║ hi  ║
 ╚═════╝
style4:
  +────+
  | hi |
  +────+
style5:
 +-----+
 | hi  |
 +-----+
style6:
 +=======+
 |   hi  |
 +=======+
*/

// Banner 在终端中绘制横幅样式的信息
type Banner struct {
	// Contents 横幅显示的内容
	Contents []string
	// Padding 内边距
	Padding int
	// Margin 外边距
	Margin int
	// Width 横幅宽度
	//  - 0 表示自动计算
	Width int
	// TextColor 文本颜色 tag
	TextColor string
	// BorderStyle 边框样式
	BorderStyle BorderStyle
}

// BannerOpFunc definition
type BannerOpFunc func(*Banner)

// BorderStyle 边框样式
type BorderStyle struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal  rune
	Vertical    rune
	// Color 边框颜色 tag
	Color string
}

// 预定义的边框样式
var (
	SimpleBorderStyle  = BorderStyle{TopLeft: '+', TopRight: '+', BottomLeft: '+', BottomRight: '+', Horizontal: '-', Vertical: '|'}
	RoundedBorderStyle = BorderStyle{TopLeft: '╭', TopRight: '╮', BottomLeft: '╰', BottomRight: '╯', Horizontal: '─', Vertical: '│'}
	SharpBorderStyle   = BorderStyle{TopLeft: '┌', TopRight: '┐', BottomLeft: '└', BottomRight: '┘', Horizontal: '─', Vertical: '│'}
)

// NewBanner1 创建新的 Banner 实例
func NewBanner1(content string, fns ...BannerOpFunc) *Banner {
	return NewBanner([]string{content}, fns...)
}

// NewBanner 创建新的 Banner 实例
func NewBanner(content []string, fns ...BannerOpFunc) *Banner {
	b := &Banner{
		Contents: content,
		Padding:  1,
		Margin:   0,
		Width:    0, // 0 表示自动计算
		// 默认使用圆角边框样式
		BorderStyle: RoundedBorderStyle,
	}
	return b.WithOptionFns(fns)
}

// WithOptionFn 设置 Banner 的选项
func (b *Banner) WithOptionFn(fns ...BannerOpFunc) *Banner {
	return b.WithOptionFns(fns)
}

// WithOptionFns 设置 Banner 的选项
func (b *Banner) WithOptionFns(fns []BannerOpFunc) *Banner {
	for _, fn := range fns {
		fn(b)
	}
	return b
}

func (b *Banner) Println() {
	color.Fprintln(Output, b.Render())
}

// Render 绘制横幅样式的信息
func (b *Banner) Render() string {
	if len(b.Contents) == 0 {
		return ""
	}

	// 计算最大内容宽度
	maxContentWidth := 0
	for _, line := range b.Contents {
		lineWidth := strutil.TextWidth(line)
		if lineWidth > maxContentWidth {
			maxContentWidth = lineWidth
		}
	}

	// 计算横幅总宽度
	contentWidth := maxContentWidth
	if b.Width > 0 && b.Width > contentWidth {
		contentWidth = b.Width
	}

	// +2: 左右边框宽度
	totalWidth := contentWidth + 2 + b.Padding*2

	// 构建横幅
	var lines []string
	tbHorizontal := strings.Repeat(string(b.BorderStyle.Horizontal), totalWidth-2)

	// 顶部边框
	topLine := string(b.BorderStyle.TopLeft) + tbHorizontal + string(b.BorderStyle.TopRight)
	lines = append(lines, topLine)

	// 内容行
	vBorderChar := string(b.BorderStyle.Vertical)
	for _, line := range b.Contents {
		lineWidth := strutil.TextWidth(line)
		// 计算左右填充
		leftPad := strings.Repeat(" ", b.Padding)
		rightPad := strings.Repeat(" ", b.Padding+(contentWidth-lineWidth))

		fullLine := vBorderChar + leftPad + line + rightPad + vBorderChar
		lines = append(lines, fullLine)
	}

	// 底部边框
	bottomLine := string(b.BorderStyle.BottomLeft) + tbHorizontal + string(b.BorderStyle.BottomRight)
	lines = append(lines, bottomLine)

	// 添加外边距
	if b.Margin > 0 {
		marginLine := strings.Repeat(" ", b.Margin)
		for i := range lines {
			lines[i] = marginLine + lines[i]
		}
	}

	return strings.Join(lines, "\n")
}
