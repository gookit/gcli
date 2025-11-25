package table

import "github.com/gookit/goutil/strutil"

const (
	OverflowWrap = 0
	OverflowCut = 1
)

// Options struct
type Options struct {
	Style
	// Alignment column 内容对齐方式
	Alignment strutil.PosFlag

	// ColMaxWidth column max width.
	//  - 0: auto
	//  - 超出宽度时，将对内容进行处理 OverflowFlag
	ColMaxWidth int
	// ColPadding column value padding left and right.
	//  - 默认L,R填充一个空格
	ColPadding string
	// OverflowFlag 内容溢出处理方式 0: 默认换行, 1: 截断
	OverflowFlag uint8
	// ShowRowNumber 显示行号，将会多一个列
	ShowRowNumber bool
	// ColumnWidths 自定义设置列宽. 按顺序设置，不设置时，将根据内容自动计算
	ColumnWidths []int

	// SortColumn sort rows by column index value.
	//  - -1: 不排序
	//  - 0+: 按指定列索引排序
	SortColumn int
	// SortAscending sort direction, true for ascending, false for descending
	SortAscending bool
	// TrimSpace trim spaces from cell values. default: true
	TrimSpace bool
	// CSVOutput output table in CSV format
	CSVOutput bool

	// -- control border show

	// ShowBorder show borderline
	ShowBorder bool
	// RowBorder show row border
	RowBorder bool
	// HeadBorder show head border
	HeadBorder bool
	// WrapBorder wrap(l,r,t,b) border for table
	WrapBorder bool
}

// NewOptions create default options
func NewOptions() *Options {
	return &Options{
		Style:        StyleDefault,
		ColPadding:   " ",
		SortColumn:   -1,
		SortAscending: true,
		TrimSpace:    true,
		ShowBorder:   true,
		HeadBorder:   true,
		WrapBorder:   true,
	}
}

// OptionFunc define
type OptionFunc func(opts *Options)

// WithStyle set table style
func WithStyle(style Style) OptionFunc {
	return func(opts *Options) { opts.Style = style }
}

// WithColPadding set column padding
func WithColPadding(padding string) OptionFunc {
	return func(opts *Options) { opts.ColPadding = padding }
}

// WithOverflowFlag set overflow handling flag
func WithOverflowFlag(flag uint8) OptionFunc {
	return func(opts *Options) { opts.OverflowFlag = flag }
}

// WithShowRowNumber enable/disable row number display
func WithShowRowNumber(show bool) OptionFunc {
	return func(opts *Options) { opts.ShowRowNumber = show }
}

// WithSortColumn set sort column and direction
func WithSortColumn(col int, ascending bool) OptionFunc {
	return func(opts *Options) {
		opts.SortColumn = col
		opts.SortAscending = ascending
	}
}

// WithTrimSpace enable/disable trim space
func WithTrimSpace(trim bool) OptionFunc {
	return func(opts *Options) { opts.TrimSpace = trim }
}

// WithCSVOutput enable/disable CSV output
func WithCSVOutput(csv bool) OptionFunc {
	return func(opts *Options) { opts.CSVOutput = csv }
}
