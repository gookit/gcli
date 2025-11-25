package table

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
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
	// ColPadding column value l,r.
	//  - 默认L,R填充一个空格
	ColPadding string
	// OverflowFlag 内容溢出处理方式 0: 默认换行, 1: 截断
	OverflowFlag uint8
	// ShowRowNumber 显示行号，将会多一个列
	ShowRowNumber bool
	// ColumnWidths 自定义设置列宽. 按顺序设置，不设置时，将根据内容自动计算
	ColumnWidths []int

	// SortColumn sort rows by column index value.
	//
	//  -1: 不排序
	SortColumn int
	// SortAscending sort direction, true for ascending
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
		Style:      StyleDefault,
		ShowBorder: true,
		HeadBorder: true,
		WrapBorder: true,
	}
}

// OptionFunc define
type OptionFunc func(opts *Options)

// WithStyle set table style
func WithStyle(style Style) OptionFunc {
	return func(opts *Options) { opts.Style = style }
}

// Table a cli Table show
type Table struct {
	show.Base // use for internal
	// options ...
	opts *Options

	// Title for the table
	Title string
	// Heads the table head data
	Heads []string
	// Rows table data rows
	Rows []*Row

	// column value align type.
	// key is col index. start from 0.
	colAlign map[int]strutil.PosFlag

	// 计算后的列宽
	colWidths []int
}

// New create table
func New(title string, fns ...OptionFunc) *Table {
	t := &Table{
		Title: title,
		opts: NewOptions(),
	}

	return t.WithOptions(fns...)
}

// WithOptions for table
func (t *Table) WithOptions(fns ...OptionFunc) *Table {
	for _, fn := range fns {
		fn(t.opts)
	}
	return t
}

// SetHeads column names to table
func (t *Table) SetHeads(names ...string) *Table {
	t.Heads = names
	return t
}

// AddRow data to table
func (t *Table) AddRow(cols ...any) *Table {
	tr := &Row{
		Cells:     make([]*Cell, 0, len(cols)),
		Separator: '|',
	}

	for _, colVal := range cols {
		cell := &Cell{
			Width: 0,
			Wrap:  false,
			Align: 0,
			Val:   colVal,
		}

		tr.Cells = append(tr.Cells, cell)
	}

	t.Rows = append(t.Rows, tr)
	return t
}

// SetRows to table
func (t *Table) SetRows(rs any) *Table {
	t.Rows = nil // 清空现有行

	switch v := rs.(type) {
	case [][]any:
		// 二维切片
		for _, row := range v {
			t.AddRow(row...)
		}
	case []map[string]any:
		// map 切片，需要建立表头（如果还没有）
		if len(t.Heads) == 0 {
			if len(v) > 0 {
				// 从第一个 map 中提取键作为表头
				for k := range v[0] {
					t.Heads = append(t.Heads, k)
				}
			}
		}

		for _, m := range v {
			rowData := make([]any, len(t.Heads))
			for i, head := range t.Heads {
				rowData[i] = m[head]
			}
			t.AddRow(rowData...)
		}
	case []any:
		// 一维切片，作为单行处理
		t.AddRow(v...)
	default:
		// 尝试使用反射处理其他类型
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}

		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < rv.Len(); i++ {
				elem := rv.Index(i)
				if elem.Kind() == reflect.Interface {
					elem = elem.Elem()
				}

				switch elem.Kind() {
				case reflect.Slice, reflect.Array:
					// 二维数组/切片
					rowData := make([]any, elem.Len())
					for j := 0; j < elem.Len(); j++ {
						rowData[j] = elem.Index(j).Interface()
					}
					t.AddRow(rowData...)
				case reflect.Map:
					// map 类型，需要建立表头（如果还没有）
					if len(t.Heads) == 0 {
						mapKeys := elem.MapKeys()
						for _, key := range mapKeys {
							t.Heads = append(t.Heads, fmt.Sprintf("%v", key.Interface()))
						}
					}

					rowData := make([]any, len(t.Heads))
					for i, head := range t.Heads {
						mapVal := elem.MapIndex(reflect.ValueOf(head))
						if mapVal.IsValid() {
							rowData[i] = mapVal.Interface()
						} else {
							rowData[i] = ""
						}
					}
					t.AddRow(rowData...)
				case reflect.Struct:
					// 结构体，需要建立表头（如果还没有）
					if len(t.Heads) == 0 {
						rt := elem.Type()
						for i := 0; i < rt.NumField(); i++ {
							field := rt.Field(i)
							// 获取 json 标签或字段名
							jsonTag := field.Tag.Get("json")
							if jsonTag != "" && jsonTag != "-" {
								// 处理 json 标签中的选项，如 "name,omitempty"
								if commaPos := strings.Index(jsonTag, ","); commaPos != -1 {
									jsonTag = jsonTag[:commaPos]
								}
								t.Heads = append(t.Heads, jsonTag)
							} else if field.PkgPath == "" { // 导出的字段
								t.Heads = append(t.Heads, field.Name)
							}
						}
					}

					rowData := make([]any, len(t.Heads))
					rt := elem.Type()
					for i, head := range t.Heads {
						// 查找匹配的字段
						for j := 0; j < rt.NumField(); j++ {
							field := rt.Field(j)
							jsonTag := field.Tag.Get("json")
							fieldName := field.Name

							if jsonTag != "" && jsonTag != "-" {
								if commaPos := strings.Index(jsonTag, ","); commaPos != -1 {
									jsonTag = jsonTag[:commaPos]
								}
								if jsonTag == head {
									rowData[i] = elem.Field(j).Interface()
									break
								}
							} else if fieldName == head && field.PkgPath == "" {
								rowData[i] = elem.Field(j).Interface()
								break
							}
						}
					}
					t.AddRow(rowData...)
				default:
					// 单个元素作为一行
					t.AddRow(elem.Interface())
				}
			}
		default:
			t.SetErr(errorx.Rf("Unsupported data type: %v", reflect.TypeOf(v)))
		}
	}

	return t
}

// String format as string
func (t *Table) String() string {
	t.Format()
	return t.Buffer().String()
}

// Print formatted message
func (t *Table) Print() {
	t.Format()
	t.Base.Print()
}

// Println formatted message with newline
func (t *Table) Println() {
	t.Format()
	t.Base.Println()
}

// Render formatted message with newline
func (t *Table) Render() {
	t.Format()
	t.Base.Println()
}

// Format as string
func (t *Table) Format() {
	// 清空缓冲区
	t.Buffer().Reset()

	t.prepare()

	t.formatHeader()

	t.formatBody()

	t.formatFooter()
}

func (t *Table) prepare() {
	// 计算列数
	colCount := len(t.Heads)
	for _, row := range t.Rows {
		if len(row.Cells) > colCount {
			colCount = len(row.Cells)
		}
		row.Separator = t.opts.Border.Cell
	}

	// 初始化列宽数组
	colWidths := make([]int, colCount)

	// 计算表头列宽
	for i, head := range t.Heads {
		width := strutil.Utf8Width(head)
		if t.opts.ColMaxWidth != 0 && width > t.opts.ColMaxWidth {
			width = t.opts.ColMaxWidth
		}
		if width > colWidths[i] {
			colWidths[i] = width
		}
	}

	// 计算数据列宽
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			cellWidth := cell.MaxWidth()
			if t.opts.ColMaxWidth != 0 && cellWidth > t.opts.ColMaxWidth {
				cellWidth = t.opts.ColMaxWidth
			}
			if cellWidth > colWidths[i] {
				colWidths[i] = cellWidth
			}
		}
	}

	// 为 Cell 设置宽度
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			if i < len(colWidths) {
				cell.Width = colWidths[i]
			}
		}
	}

	// 保存计算后的列宽到 table 实例
	t.colWidths = colWidths
}

// Format as string
func (t *Table) formatHeader() {
	buf := t.Buffer()

	if len(t.Heads) == 0 && len(t.Rows) == 0 {
		return // 没有表头和数据，直接返回
	}

	opts := t.opts
	style := opts.Style

	// 如果有标题，先打印标题
	if t.Title != "" {
		buf.WriteString(t.Title + "\n")
	}

	// 画顶部边框（如果需要）
	if opts.ShowBorder && style.Border.TopLeft != 0 {
		t.drawBorderLine(t.Buffer(), style.Border.TopLeft, style.Border.Top, style.Border.TopIntersect, style.Border.TopRight)
	}

	// 打印表头
	if len(t.Heads) > 0 {
		// 特殊处理 Markdown 样式
		if style == StyleMarkdown {
			// 对于 Markdown 样式，先打印表头内容
			buf.WriteRune(style.Border.Right) // 左边框

			for i, head := range t.Heads {
				if i < len(t.colWidths) {
					// 使用 strutil.Resize 来对齐表头内容
					resized := strutil.Resize(head, t.colWidths[i], t.opts.Alignment)
					// 应用颜色（优先使用 FirstColor 给第一列）
					var coloredHead string
					if i == 0 && opts.FirstColor != "" {
						// 表头第一列使用 FirstColor
						coloredHead = color.Sprintf("<%s>%s</>", opts.FirstColor, resized)
					} else if opts.HeadColor != "" {
						// 其他列使用 HeadColor
						coloredHead = color.Sprintf("<%s>%s</>", opts.HeadColor, resized)
					} else {
						coloredHead = resized
					}
					buf.WriteString(coloredHead)
				} else {
					buf.WriteString(head)
				}

				if i < len(t.Heads)-1 { // 不是最后一个元素
					buf.WriteRune(style.Border.Cell) // 列分隔符
				}
			}

			buf.WriteRune(style.Border.Right) // 右边框
			buf.WriteString("\n")

			// 然后打印 Markdown 风格的分隔行
			buf.WriteRune(style.Border.Right) // 左边框
			for i := 0; i < len(t.Heads); i++ {
				if i < len(t.colWidths) {
					// Markdown 表格分隔符，通常为至少3个连字符
					sepWidth := t.colWidths[i]
					if sepWidth < 3 {
						sepWidth = 3
					}
					buf.WriteString(strings.Repeat("-", sepWidth))
				}

				if i < len(t.Heads)-1 { // 不是最后一个元素
					buf.WriteRune(style.Border.Cell) // 列分隔符
				}
			}

			buf.WriteRune(style.Border.Right) // 右边框
			buf.WriteString("\n")
		} else {
			// 普通表格样式
			buf.WriteRune(style.Border.Right) // 左边框

			for i, head := range t.Heads {
				if i < len(t.colWidths) {
					// 使用 strutil.Resize 来对齐表头内容
					resized := strutil.Resize(head, t.colWidths[i], t.opts.Alignment)
					// 应用颜色（优先使用 FirstColor 给第一列）
					var coloredHead string
					if i == 0 && opts.FirstColor != "" {
						// 表头第一列使用 FirstColor
						coloredHead = color.Sprintf("<%s>%s</>", opts.FirstColor, resized)
					} else if opts.HeadColor != "" {
						// 其他列使用 HeadColor
						coloredHead = color.Sprintf("<%s>%s</>", opts.HeadColor, resized)
					} else {
						coloredHead = resized
					}
					buf.WriteString(coloredHead)
				} else {
					buf.WriteString(head)
				}

				if i < len(t.Heads)-1 { // 不是最后一个元素
					buf.WriteRune(style.Border.Cell) // 列分隔符
				}
			}

			buf.WriteRune(style.Border.Right) // 右边框
			buf.WriteString("\n")

			// 画表头分隔线（如果需要）
			if opts.HeadBorder {
				t.drawBorderLine(buf, style.Divider.Left, style.Border.Center, style.Divider.Intersect, style.Divider.Right)
			} else if opts.RowBorder {
				t.drawBorderLine(buf, style.Border.Right, style.Border.Center, style.Border.Cell, style.Border.Right)
			}
		}
	} else if len(t.Heads) == 0 && len(t.Rows) > 0 {
		// 没有表头但有数据，仍可能需要画分隔线
		if opts.HeadBorder {
			t.drawBorderLine(buf, style.Divider.Left, style.Border.Center, style.Divider.Intersect, style.Divider.Right)
		}
	}
}

// Format as string
func (t *Table) formatBody() {
	buf := t.Buffer()
	opts := t.opts
	style := opts.Style

	for i, row := range t.Rows {
		// 表格样式
		buf.WriteRune(style.Border.Right) // 左边框

		// 处理每列
		for j := 0; j < len(t.colWidths); j++ {
			if j < len(row.Cells) {
				cell := row.Cells[j]
				cellStr := cell.String()

				// 应用对齐方式
				var align strutil.PosFlag
				if cell.Align != 0 {
					align = cell.Align
				} else {
					align = opts.Alignment
				}

				// 根据宽度调整内容
				if cell.Width > 0 {
					cellStr = strutil.Resize(cellStr, cell.Width, align)
				}

				// 应用颜色（如果设置了行颜色或首列颜色）
				var coloredCell string
				if j == 0 && opts.FirstColor != "" {
					// 首列使用 FirstColor
					coloredCell = color.Sprintf("<%s>%s</>", opts.FirstColor, cellStr)
				} else if opts.RowColor != "" {
					// 其他列使用 RowColor
					coloredCell = color.Sprintf("<%s>%s</>", opts.RowColor, cellStr)
				} else {
					coloredCell = cellStr
				}
				buf.WriteString(coloredCell)
			} else {
				// 如果这一行没有足够的列，使用空格填充
				if j < len(t.colWidths) {
					buf.WriteString(strings.Repeat(" ", t.colWidths[j]))
				}
			}

			if j < len(t.colWidths)-1 { // 不是最后一个元素
				buf.WriteRune(style.Border.Cell) // 列分隔符
			}
		}

		buf.WriteRune(style.Border.Right) // 右边框
		buf.WriteString("\n")

		// 画行分隔线（如果需要）
		if opts.RowBorder && i < len(t.Rows)-1 {
			t.drawBorderLine(buf, style.Border.Right, style.Border.Center, style.Border.Cell, style.Border.Right)
		}

	}
}

// Format as string
func (t *Table) formatFooter() {
	buf := t.Buffer()
	opts := t.opts
	style := opts.Style

	// 画底部边框（如果需要）
	if opts.ShowBorder && opts.WrapBorder && style.Border.BottomLeft != 0 {
		t.drawBorderLine(buf, style.Border.BottomLeft, style.Border.Bottom, style.Border.BottomIntersect, style.Border.BottomRight)
	}
}

// drawBorderLine draws a border line with the given characters
func (t *Table) drawBorderLine(buf *bytes.Buffer, leftChar, centerChar, intersectChar, rightChar rune) {
	if leftChar == 0 && rightChar == 0 {
		return // 如果没有边框字符，则跳过
	}

	buf.WriteRune(leftChar)

	for i, width := range t.colWidths {
		for j := 0; j < width; j++ {
			buf.WriteRune(centerChar)
		}
		if i < len(t.colWidths)-1 { // 不是最后一个列
			buf.WriteRune(intersectChar) // 列间分隔符
		}
	}

	buf.WriteRune(rightChar)
	buf.WriteString("\n")
}

// WriteTo format table to string and write to w.
func (t *Table) WriteTo(w io.Writer) (int64, error) {
	t.Format()
	return t.Buffer().WriteTo(w)
}

// Row represents a row in a table
type Row struct {
	// Cells is the group of cell for the row
	Cells []*Cell

	// Separator for table columns
	Separator rune
}

// Cell represents a column in a row
type Cell struct {
	// Width is the width of the cell
	Width int
	// Wrap when true wraps the contents of the cell when the length exceeds the width
	Wrap bool
	// Align when true aligns contents to the right
	Align strutil.PosFlag

	// Val is the cell data
	Val any
	str string // string cache of Val
}

// MaxWidth returns the max width of all the lines in a cell
func (c *Cell) MaxWidth() int {
	width := 0
	for _, s := range strings.Split(c.String(), "\n") {
		w := strutil.Utf8Width(s)
		if w > width {
			width = w
		}
	}

	return width
}

// String returns the string formatted representation of the cell
func (c *Cell) String() string {
	if c.Val == nil {
		return strutil.PadLeft(" ", " ", c.Width)
		// return c.str
	}

	s := strutil.SafeString(c.Val)
	if c.Width == 0 {
		return s
	}

	if c.Wrap && len(s) > c.Width {
		return strutil.WordWrap(s, c.Width)
	}
	return strutil.Resize(s, c.Width, c.Align)
}
