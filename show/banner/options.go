package banner

import (
	"github.com/gookit/gcli/v3/show/showcom"
	"github.com/gookit/goutil/comdef"
)

// 特殊margin定义banner 居中，居右
const (
	AtCenter = -1
	AtRight  = -2
)

// Options banner options
//
// 宽度设置优先级: MinWidth > Width > PercentWidth > contentWidth
type Options struct {
	// Padding 内边距 default: 1
	Padding int
	// MarginL 左外边距
	//  - >0 表示左外边距
	//  - 0 表示无外边距 居左显示banner(默认)
	//  - -1 表示居中显示banner
	//  - -2 表示居右显示banner
	MarginL int
	MarginR int // TODO
	// MarginTB 上下外边距 TODO
	MarginTB []int
	// Width 横幅宽度
	//  - 0 表示自动计算内容宽度
	Width int
	// Height 内容高度行数，默认为内容行数
	Height int
	// MinWidth 最小宽度
	//  0: 表示不限制最小宽度
	MinWidth int
	// PercentWidth 使用终端宽度的百分比宽度 (1-100)
	//  0 表示不使用百分比宽度
	PercentWidth int
	// OverflowFlag 内容溢出处理 default: wrap
	OverflowFlag showcom.OverflowFlag
	// Alignment 内容对齐方式
	Alignment comdef.Align
	// TextColor 文本颜色 tag
	TextColor string
	// BorderStyle 边框样式
	BorderStyle BorderStyle
}

// OptionFunc definition
type OptionFunc func(b *Options)

// WithHeight 设置内容高度行数
func WithHeight(height int) OptionFunc {
	return func(b *Options) {
		b.Height = height
	}
}

// WithMinWidth 设置最小宽度
func WithMinWidth(minWidth int) OptionFunc {
	return func(b *Options) {
		b.MinWidth = minWidth
	}
}

// WithPercentWidth 使用终端宽度的百分比宽度
func WithPercentWidth(percent int) OptionFunc {
	return func(b *Options) {
		b.PercentWidth = percent
	}
}

// WithOverflowFlag 设置内容溢出处理方式
func WithOverflowFlag(flag showcom.OverflowFlag) OptionFunc {
	return func(b *Options) {
		b.OverflowFlag = flag
	}
}

// WithAlignment 设置内容对齐方式
func WithAlignment(alignment comdef.Align) OptionFunc {
	return func(b *Options) {
		b.Alignment = alignment
	}
}

// WithWidth 设置横幅宽度
func WithWidth(width int) OptionFunc {
	return func(b *Options) {
		b.Width = width
	}
}

// WithMarginLeft 设置左边距
func WithMarginLeft(margin int) OptionFunc {
	return func(b *Options) {
		b.MarginL = margin
	}
}

// WithBannerCenter 设置banner居中
func WithBannerCenter() OptionFunc {
	return func(b *Options) {
		b.MarginL = AtCenter
	}
}

// WithBannerRight 设置banner居右
func WithBannerRight() OptionFunc {
	return func(b *Options) {
		b.MarginL = AtRight
	}
}

// WithMarginTopBottom 添加上下边距
func WithMarginTopBottom(marginTop, marginBottom int) OptionFunc {
	return func(b *Options) {
		b.MarginTB = []int{marginTop, marginBottom}
	}
}

// BorderStyle 边框样式
type BorderStyle struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal rune // top, bottom
	Vertical   rune // left, right
	// Color 边框颜色 tag
	Color string
}

// 预定义的边框样式
var (
	/*
		SimpleBorderStyle example:

		 +-----+
		 | hi  |
		 +-----+
	*/
	SimpleBorderStyle  = BorderStyle{TopLeft: '+', TopRight: '+', BottomLeft: '+', BottomRight: '+', Horizontal: '-', Vertical: '|'}
	/*
		RoundedBorderStyle example
		 ╭─────╮
		 │  hi │
		 ╰─────╯
	*/
	RoundedBorderStyle = BorderStyle{TopLeft: '╭', TopRight: '╮', BottomLeft: '╰', BottomRight: '╯', Horizontal: '─', Vertical: '│'}
	/*
		SharpBorderStyle example:
		 ┌─────┐
		 │  hi │
		 └─────┘
	*/
	SharpBorderStyle   = BorderStyle{TopLeft: '┌', TopRight: '┐', BottomLeft: '└', BottomRight: '┘', Horizontal: '─', Vertical: '│'}
	/*
		DoubleBorderStyle example:
		 ╔═════╗
		 ║ hi  ║
		 ╚═════╝
	*/
	DoubleBorderStyle = BorderStyle{TopLeft: '╔', TopRight: '╗', BottomLeft: '╚', BottomRight: '╝', Horizontal: '═', Vertical: '║'}
)
