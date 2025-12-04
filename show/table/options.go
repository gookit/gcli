package table

import "github.com/gookit/goutil/strutil"

// OverflowFlag for handling content overflow
type OverflowFlag uint8

// OverflowFlag values
const (
	OverflowAuto OverflowFlag = iota // auto: default is cut
	OverflowCut                      // 截断
	OverflowWrap                     // 换行
)

// BorderFlags for controlling table border display
const (
	BorderNone uint8 = 0
	BorderTop  uint8 = 1 << iota
	BorderBottom
	BorderLeft
	BorderRight
	BorderHeader // 显示表头边框
	BorderRows   // 显示行边框
	// BorderDefault display head, top and bottom borders
	BorderDefault = BorderTop | BorderBottom | BorderHeader
	// BorderBody display body borders
	BorderBody = BorderTop | BorderBottom | BorderLeft | BorderRight | BorderRows
	// BorderAll display all borders
	BorderAll = BorderTop | BorderBottom | BorderLeft | BorderRight | BorderHeader | BorderRows
)

// Options struct
type Options struct {
	Style
	// Alignment column 内容对齐方式
	Alignment strutil.PosFlag
	// StructTag struct tag name on input struct data. Default: json
	StructTag string

	// ColMaxWidth column max width.
	//  - 0: auto
	//  - 超出宽度时，将对内容进行处理 OverflowFlag
	ColMaxWidth int
	// CellPadding column value padding left and right.
	//  - 默认L,R填充一个空格
	CellPadding string
	// OverflowFlag 内容溢出处理方式 0: 默认截断, 1: 换行
	OverflowFlag OverflowFlag
	// ShowRowNumber 显示行号，将会多一个列
	ShowRowNumber bool
	// ColumnWidths 自定义设置列宽. 按顺序设置，不设置/为0时，将根据内容自动计算
	//
	// start index is 0
	ColumnWidths []int

	// SortColumn sort rows by column index value.
	//  - -1: 不排序
	//  - 0+: 按指定索引列排序
	SortColumn int
	// SortAscending sort direction, true for ascending, false for descending
	SortAscending bool
	// TrimSpace trim spaces from cell values. default: true
	TrimSpace bool
	// CSVOutput output table in CSV format TODO
	CSVOutput bool
}

// HasBorderFlag check border flag
func (opts *Options) HasBorderFlag(flag uint8) bool {
	return (opts.BorderFlags & flag) != 0
}

// NewOptions create default options
func NewOptions() *Options {
	return &Options{
		Style: StyleSimple,
		StructTag:   "json",
		CellPadding: " ",
		SortColumn:  -1,
		SortAscending: true,
		TrimSpace:   true,
	}
}

// OptionFunc define
type OptionFunc func(opts *Options)

// WithStyle set table style
func WithStyle(style Style) OptionFunc {
	return func(opts *Options) { opts.Style = style }
}

// WithoutBorder set border to none
func WithoutBorder() OptionFunc {
	return func(opts *Options) { opts.BorderFlags = BorderNone }
}

// WithBorderFlags set border flags
func WithBorderFlags(flags uint8) OptionFunc {
	return func(opts *Options) { opts.BorderFlags = flags }
}

// WithCellPadding set column padding
func WithCellPadding(padding string) OptionFunc {
	return func(opts *Options) { opts.CellPadding = padding }
}

// WithColMaxWidth set column max width
func WithColMaxWidth(width int) OptionFunc {
	return func(opts *Options) { opts.ColMaxWidth = width }
}

// WithColumnWidths set column widths by index
func WithColumnWidths(widths ...int) OptionFunc {
	return func(opts *Options) { opts.ColumnWidths = widths }
}

// WithOverflowFlag set overflow handling flag
func WithOverflowFlag(flag OverflowFlag) OptionFunc {
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
