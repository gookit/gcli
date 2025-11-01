package title

import (
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/gcli/v3/show/symbols"
	"github.com/gookit/goutil/strutil"
)

// Title 在终端中打印标题行
type Title struct {
	Title string
	Color string // 颜色Tag
	// Formatter 自定义格式化处理，不设置时，使用默认格式化处理
	Formatter func(t *Title) string
	// Formatter IFormatter

	// PaddingLR 是否左右填充 Char
	PaddingLR bool
	// Char 左右填充字符
	Char rune

	// 是否显示上下边框
	ShowBorder bool
	BorderChar rune
	// BorderPos 边框位置 0: 无, 1: 上, 2: 下, 4: 上下
	BorderPos gclicom.BorderPos

	Width int // 总的显示宽度
	Indent int
	Align gclicom.TextPos
}

// OptionFunc definition
type OptionFunc func(t *Title)

// New Title instance
func New(title string, fns ...OptionFunc) *Title {
	t := &Title{
		Title: title,
		Width: 80,
		Char:  symbols.Equal,
		// Indent: 2,
		Align: gclicom.TextPosLeft,
		Color: "green1",
		// Border
		ShowBorder: false,
		BorderChar: '-',
		BorderPos:  gclicom.BorderPosBottom,
		// Padding
		PaddingLR: true,
	}
	return t.WithOptionFns(fns)
}

// WithBorderTop setting the title border to top
func WithBorderTop() OptionFunc {
	return func(t *Title) {
		t.ShowBorder = true
		t.BorderPos = gclicom.BorderPosTop
	}
}

// WithBorderBottom setting the title border to bottom
func WithBorderBottom() OptionFunc {
	return func(t *Title) {
		t.ShowBorder = true
		t.BorderPos = gclicom.BorderPosBottom
	}
}

// WithBorderBoth setting the title border to both top and bottom
func WithBorderBoth() OptionFunc {
	return func(t *Title) {
		t.ShowBorder = true
		t.BorderPos = gclicom.BorderPosTB
	}
}

// WithoutBorder setting the title border to none
func WithoutBorder() OptionFunc {
	return func(t *Title) {
		t.ShowBorder = false
	}
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

// Println print title line
func (t *Title) Println() {
	color.Fprintln(gclicom.Output, t.Render())
}

func (t *Title) Render() string {
	// 如果有自定义格式化函数，优先使用
	if t.Formatter != nil {
		return t.Formatter(t)
	}

	// 计算实际可用宽度（减去缩进）
	availableWidth := t.Width - t.Indent
	if availableWidth <= 0 {
		availableWidth = t.Width
	}

	// 根据对齐方式处理标题
	var content string
	switch t.Align {
	case gclicom.TextPosCenter:
		content = t.renderCenter(availableWidth)
	case gclicom.TextPosRight:
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
		return t.renderWithBorder(content, availableWidth)
	}

	return content
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
			return string(t.Char) + " " + t.title()
		}

		remaining := width - titleLen - 2
		rightChars := make([]rune, remaining)
		for i := range rightChars {
			rightChars[i] = t.Char
		}
		return string(t.Char) + " " + t.title() + " " + string(rightChars)
	}

	// 不填充: Title CHAR
	// remaining := width - titleLen
	// chars := make([]rune, remaining)
	// for i := range chars {
	// 	chars[i] = t.Char
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
			leftChars[i] = t.Char
		}
		for i := range rightChars {
			rightChars[i] = t.Char
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
	// 	rightChars[i] = t.Char
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
			leftChars[i] = t.Char
		}
		return string(leftChars) + " " + t.title() + " " + string(t.Char)
	}

	// 不填充: CHAR Title
	return t.title()
	// remaining := width - titleLen
	// leftChars := make([]rune, remaining)
	// for i := range leftChars {
	// 	leftChars[i] = t.Char
	// }
	// return string(leftChars) + " " + t.title()
}

func (t *Title) title() string {
	return color.WrapTag(t.Title, t.Color)
}
