package title

import (
	"github.com/gookit/gcli/v3/gclicom"
	"github.com/gookit/goutil/comdef"
)

// Options title options
type Options struct {
	// Color 颜色Tag
	Color string
	// PaddingLR 是否左右填充 PaddingChar
	PaddingLR bool
	// PaddingChar 左右填充字符
	PaddingChar rune

	// 是否显示上下边框
	ShowBorder bool
	BorderChar rune
	// BorderPos 边框位置 0: 无, 1: 上, 2: 下, 4: 上下
	BorderPos gclicom.BorderPos

	// 总的显示宽度
	Width int
	Indent int
	Align comdef.Align
}

// OptionFunc definition
type OptionFunc func(t *Title)

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

// WithAlignRight setting the title align to right
func WithAlignRight() OptionFunc {
	return func(t *Title) {
		t.Align = comdef.Right
	}
}

// WithAlignCenter setting the title align to center
func WithAlignCenter() OptionFunc {
	return func(t *Title) {
		t.Align = comdef.Center
	}
}
