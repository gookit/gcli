package banner

import (
	"fmt"
	"strings"

	"github.com/gookit/gcli/v3/show/showcom"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/strutil"
	"github.com/gookit/goutil/x/termenv"
)

/*
eg:
   ╭──────────────────────────────────────────────────────────────────╮
   │                                                                  │
   │                Update available! 3.21.0 → 3.27.0.                │
   │   Changelog: https://github.com/gookit/gcli/releases/tag/v3.2.0  │
   │                Run "x y z" to update.                			  │
   │                                                                  │
   ╰──────────────────────────────────────────────────────────────────╯

style4:
  +────+
  | hi |
  +────+

style6:
 +=======+
 |   hi  |
 +=======+
*/

// Banner 在终端中绘制横幅样式的信息
type Banner struct {
	// use for internal
	showcom.Base
	// Options banner options
	Options
	// Contents 横幅显示的内容
	Contents []string

	// context data 存储渲染过程中计算的数据
	termWidth  int
	innerWidth int
	totalWidth int // boxWidth = innerWidth + 2 + padding*2
	// contentLines 存储处理后的内容行
	contentLines []string
	tbHorizontal string // top, bottom line
	// margin padding string
	marginStr string
}

// New 创建新的 Banner 实例
func New(content any, fns ...OptionFunc) *Banner {
	b := &Banner{
		Options: Options{
			Padding:      1,
			OverflowFlag: showcom.OverflowWrap,
			Alignment:    comdef.Left, // 默认左对齐
			// 默认使用圆角边框样式
			BorderStyle: RoundedBorderStyle,
		},
	}
	b.FormatFn = b.Format
	return b.WithContents(content).WithOptionFns(fns)
}

// WithOptionFn 设置 Banner 的选项
func (b *Banner) WithOptionFn(fns ...OptionFunc) *Banner { return b.WithOptionFns(fns) }

// WithOptionFns 设置 Banner 的选项
func (b *Banner) WithOptionFns(fns []OptionFunc) *Banner {
	for _, fn := range fns {
		fn(&b.Options)
	}
	return b
}

// WithContents 设置横幅显示的内容
//
//	content: string, []string, []any, ...
func (b *Banner) WithContents(content any) *Banner {
	var contents []string
	switch v := content.(type) {
	case string:
		contents = append(contents, v)
	case []string:
		contents = append(contents, v...)
	case []any:
		for _, item := range v {
			contents = append(contents, fmt.Sprint(item))
		}
	default:
		contents = append(contents, fmt.Sprint(content))
	}
	b.Contents = contents
	return b
}

// Render 绘制横幅样式的信息
func (b *Banner) Render() string {
	b.Format()
	return b.Buf.String()
}

// Format 绘制横幅
func (b *Banner) Format() {
	if len(b.Contents) == 0 {
		return
	}

	b.InitBuffer()

	// 准备渲染数据
	b.prepare()

	// 顶部边框
	b.renderTop()

	// 内容行
	b.renderBody()

	// 底部边框
	b.renderBottom()
}

// prepare 准备渲染所需的数据
func (b *Banner) prepare() {
	// 获取终端宽度
	b.termWidth, _ = termenv.GetTermSize()

	// 预处理内容 - 拆分有换行的line
	fmtContents := make([]string, 0, len(b.Contents))
	for _, content := range b.Contents {
		content = strings.TrimSpace(content)
		fmtContents = append(fmtContents, strings.Split(content, "\n")...)
	}

	// 计算最大内容宽度
	maxContentWidth := 0
	lineWidths := make([]int, 0, len(fmtContents))
	for _, line := range fmtContents {
		lineWidth := strutil.TextWidth(line)
		if lineWidth > maxContentWidth {
			maxContentWidth = lineWidth
		}
		lineWidths = append(lineWidths, lineWidth)
	}

	// 计算横幅总宽度 - 默认为内容宽度
	b.innerWidth = maxContentWidth

	// 应用百分比宽度
	if b.PercentWidth > 0 && b.termWidth > 0 {
		percentWidth := b.termWidth * b.PercentWidth / 100
		if percentWidth > 0 {
			b.innerWidth = percentWidth
		}
	}

	// 应用固定宽度
	if b.Width > 0 {
		b.innerWidth = b.Width
	}
	// 应用最小宽度
	if b.MinWidth > 0 && b.innerWidth < b.MinWidth {
		b.innerWidth = b.MinWidth
	}

	// +2: 左右边框宽度
	b.totalWidth = b.innerWidth + 2 + b.Padding*2
	// 构建横幅水平线
	b.tbHorizontal = strings.Repeat(string(b.BorderStyle.Horizontal), b.totalWidth-2)

	// 处理内容行 - 换行/截断
	b.contentLines = nil
	for i, line := range fmtContents {
		// 处理内容溢出
		if lineWidths[i] > b.innerWidth {
			if b.OverflowFlag == showcom.OverflowCut { // 截断模式
				line = strutil.TextTruncate(line, b.innerWidth, "")
			} else { // 换行模式
				splitLines := strutil.TextSplit(line, b.innerWidth)
				b.contentLines = append(b.contentLines, splitLines...)
				continue
			}
		}
		b.contentLines = append(b.contentLines, line)
	}

	// 应用高度设置
	if b.Height > 0 && len(b.contentLines) < b.Height {
		// 填充空行到指定高度
		for i := len(b.contentLines); i < b.Height; i++ {
			b.contentLines = append(b.contentLines, "")
		}
	} else if b.Height > 0 && len(b.contentLines) > b.Height {
		// 截断到指定高度
		b.contentLines = b.contentLines[:b.Height]
	}

	// 预处理外边距
	marginWidth := b.MarginL
	// 居中对齐
	if b.MarginL == AtCenter && b.termWidth > 0 {
		marginWidth = (b.termWidth - b.totalWidth) / 2
	} else if b.MarginL == AtRight && b.termWidth > 0 {
		marginWidth = b.termWidth - b.totalWidth
	}
	if marginWidth > 0 {
		b.marginStr = strings.Repeat(" ", marginWidth)
	}
}

// renderTop 渲染顶部边框
func (b *Banner) renderTop() {
	if b.marginStr != "" {
		b.Buf.WriteString(b.marginStr)
	}
	b.Buf.WriteRune(b.BorderStyle.TopLeft)
	b.Buf.WriteString(b.tbHorizontal)
	b.Buf.WriteRune(b.BorderStyle.TopRight)
	b.Buf.WriteByte('\n')
}

// renderBody 渲染内容行
func (b *Banner) renderBody() {
	// 渲染内容行
	for _, lineText := range b.contentLines {
		lineWidth := strutil.TextWidth(lineText)

		// 根据对齐方式计算左右填充
		var leftPad, rightPad string
		totalPad := b.innerWidth - lineWidth

		switch b.Alignment {
		case comdef.Center:
			leftPadCount := totalPad / 2
			rightPadCount := totalPad - leftPadCount
			leftPad = strings.Repeat(" ", b.Padding+leftPadCount)
			rightPad = strings.Repeat(" ", b.Padding+rightPadCount)
		case comdef.Right:
			leftPad = strings.Repeat(" ", b.Padding+totalPad)
			rightPad = strings.Repeat(" ", b.Padding)
		default: // Left
			leftPad = strings.Repeat(" ", b.Padding)
			rightPad = strings.Repeat(" ", b.Padding+totalPad)
		}

		// 应用外边距
		if b.marginStr != "" {
			b.Buf.WriteString(b.marginStr)
		}

		b.Buf.WriteRune(b.BorderStyle.Vertical)
		b.Buf.WriteString(leftPad)
		b.Buf.WriteString(lineText)
		b.Buf.WriteString(rightPad)
		b.Buf.WriteRune(b.BorderStyle.Vertical)
		b.Buf.WriteByte('\n')
	}
}

// renderBottom 渲染底部边框
func (b *Banner) renderBottom() {
	if b.marginStr != "" {
		b.Buf.WriteString(b.marginStr)
	}
	b.Buf.WriteRune(b.BorderStyle.BottomLeft)
	b.Buf.WriteString(b.tbHorizontal)
	b.Buf.WriteRune(b.BorderStyle.BottomRight)
}
