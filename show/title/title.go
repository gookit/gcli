package title

import (
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/gcli/v3/show/showcom"
	"github.com/gookit/gcli/v3/show/symbols"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/strutil"
)

// Title 在终端中打印标题行
type Title struct {
	showcom.Base
	Options
	Prefix string
	Title string
}

// New Title instance
func New(title string, fns ...OptionFunc) *Title {
	t := &Title{
		Title: title,
		Options: Options{
			Width:       80,
			PaddingChar: symbols.Equal,
			// Indent: 2,
			Align: comdef.Left,
			Color: "green1",
			// Border
			ShowBorder: false,
			BorderChar: '-',
			BorderPos:  gclicom.BorderPosBottom,
			// Padding
			PaddingLR: true,
		},
	}

	t.FormatFn = t.Format
	return t.WithOptionFns(fns)
}

// WithOptionFns 设置选项
func (t *Title) WithOptionFns(fns []OptionFunc) *Title {
	for _, fn := range fns {
		fn(t)
	}
	return t
}

// SetTitle set title text
func (t *Title) SetTitle(title string) *Title {
	t.Title = title
	return t
}

// ShowNew set new title and print
func (t *Title) ShowNew(title string) {
	t.SetTitle(title).Println()
}

// Render 渲染标题
func (t *Title) Render() string {
	t.Format()
	return t.Buffer().String()
}

// Format 格式化标题信息
func (t *Title) Format() {
	t.InitBuffer()

	// 计算实际可用宽度（减去缩进）
	availableWidth := t.Width - t.Indent
	if availableWidth <= 0 {
		availableWidth = t.Width
	}

	// 根据对齐方式处理标题
	var content string
	switch t.Align {
	case comdef.Center:
		content = t.renderCenter(availableWidth)
	case comdef.Right:
		content = t.renderRight(availableWidth)
	default: // left
		content = t.renderLeft(availableWidth)
	}

	// 添加缩进
	if t.Indent > 0 {
		indent := make([]rune, t.Indent)
		for i := range indent {
			indent[i] = ' '
		}
		content = string(indent) + content
	}

	// 处理边框显示
	if t.ShowBorder {
		content = t.renderWithBorder(content, availableWidth)
	}
	t.Buf.WriteString(content)
}

// renderWithBorder 添加边框处理
func (t *Title) renderWithBorder(content string, width int) string {
	// 创建边框线
	borderLine := strings.Repeat(string(t.BorderChar), width)

	// 根据边框位置添加边框
	switch t.BorderPos {
	case gclicom.BorderPosTop:
		return borderLine + "\n" + content
	case gclicom.BorderPosBottom:
		return content + "\n" + borderLine
	case gclicom.BorderPosTB: // Top & Bottom
		return borderLine + "\n" + content + "\n" + borderLine
	default:
		return content
	}
}

// renderLeft 左对齐渲染
func (t *Title) renderLeft(width int) string {
	titleLen := strutil.TextWidth(t.Title)
	if titleLen >= width {
		return t.title()
	}

	if t.PaddingLR {
		// 填充左右: CHAR Title CHAR
		if titleLen >= width-2 {
			return string(t.PaddingChar) + " " + t.title()
		}

		remaining := width - titleLen - 2
		rightChars := make([]rune, remaining)
		for i := range rightChars {
			rightChars[i] = t.PaddingChar
		}
		return string(t.PaddingChar) + " " + t.title() + " " + string(rightChars)
	}

	// 不填充: Title CHAR
	// remaining := width - titleLen
	// chars := make([]rune, remaining)
	// for i := range chars {
	// 	chars[i] = t.PaddingChar
	// }
	return t.title() // + " " + string(chars)
}

// renderCenter 居中渲染
func (t *Title) renderCenter(width int) string {
	titleLen := strutil.TextWidth(t.Title)
	if titleLen >= width {
		return t.title()
	}

	if t.PaddingLR {
		// 填充左右: CHAR Title CHAR 居中
		totalPadding := width - titleLen - 2
		leftPadding := totalPadding / 2
		rightPadding := totalPadding - leftPadding

		leftChars := make([]rune, leftPadding)
		rightChars := make([]rune, rightPadding)
		for i := range leftChars {
			leftChars[i] = t.PaddingChar
		}
		for i := range rightChars {
			rightChars[i] = t.PaddingChar
		}

		return string(leftChars) + " " + t.title() + " " + string(rightChars)
	}

	// 不填充: Title 居中
	totalPadding := width - titleLen
	leftPadding := totalPadding / 2
	// rightPadding := totalPadding - leftPadding

	leftChars := make([]rune, leftPadding)
	// rightChars := make([]rune, rightPadding)
	for i := range leftChars {
		leftChars[i] = ' '
	}
	// for i := range rightChars {
	// 	rightChars[i] = t.PaddingChar
	// }

	return string(leftChars) + " " + t.title() // + " " + string(rightChars)
}

// renderRight 右对齐渲染
func (t *Title) renderRight(width int) string {
	titleLen := strutil.TextWidth(t.Title)
	if titleLen >= width {
		return t.title()
	}

	if t.PaddingLR {
		// 填充左右: CHAR Title CHAR
		remaining := width - titleLen - 2
		leftChars := make([]rune, remaining)
		for i := range leftChars {
			leftChars[i] = t.PaddingChar
		}
		return string(leftChars) + " " + t.title() + " " + string(t.PaddingChar)
	}

	// 不填充: CHAR Title
	return t.title()
	// remaining := width - titleLen
	// leftChars := make([]rune, remaining)
	// for i := range leftChars {
	// 	leftChars[i] = t.PaddingChar
	// }
	// return string(leftChars) + " " + t.title()
}

func (t *Title) title() string {
	return color.WrapTag(t.Title, t.Color)
}
