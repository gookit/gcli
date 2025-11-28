package table

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/gookit/color"
	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/errorx"
	"github.com/gookit/goutil/strutil"
)

// Table a cli Table show
type Table struct {
	show.Base // use for internal
	// options ...
	opts *Options

	// Title for the table
	Title string
	// Heads the table head data
	Heads []*Cell
	// Heads []string
	// Rows table data rows
	Rows []*Row

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
	t.Heads = nil // reset before set
	for _, name := range names {
		t.AddHead(name)
	}
	return t
}

// AddHead add head column to table
func (t *Table) AddHead(name string) *Table {
	t.Heads = append(t.Heads, NewCell(name))
	return t
}

// PrependHead prepend head column
func (t *Table) PrependHead(name string) *Table {
	t.Heads = append([]*Cell{NewCell(name)}, t.Heads...)
	return t
}

// AddRow data to table
func (t *Table) AddRow(cols ...any) *Table {
	tr := &Row{
		Cells: make([]*Cell, 0, len(cols)),
	}

	for _, colVal := range cols {
		tr.Cells = append(tr.Cells, NewCell(colVal))
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
			// 从第一个 map 中提取键作为表头
			if len(v) > 0 {
				for k := range v[0] {
					t.AddHead(k)
				}
			}
		}

		for _, m := range v {
			rowData := make([]any, len(t.Heads))
			for i, head := range t.Heads {
				rowData[i] = m[head.String()]
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
	t.reset()

	t.prepare()

	t.formatHeader()

	t.formatBody()

	t.formatFooter()
}

func (t *Table) reset() {
	// 清空缓冲区
	t.Buffer().Reset()
	t.colWidths = nil

	for _, row := range t.Rows {
		for _, cell := range row.Cells {
			cell.init = false
			cell.width = 0
			cell.height = 0
		}
	}
}

func (t *Table) calcColWidth(width, i int) int {
	// 自定义列宽
	if len(t.opts.ColumnWidths) > i && t.opts.ColumnWidths[i] > 0 {
		width = t.opts.ColumnWidths[i]
	}
	// 列L,R填充
	if t.opts.CellPadding != "" {
		width += len(t.opts.CellPadding) * 2
	}

	// 列最大宽度
	if t.opts.ColMaxWidth != 0 && width > t.opts.ColMaxWidth {
		width = t.opts.ColMaxWidth
	}
	return width
}

func (t *Table) prepare() {
	// 如果需要显示行号，在表头前添加 "#"
	if t.opts.ShowRowNumber {
		t.PrependHead("#")
	}

	// 计算列数 + init row.Cells
	colCount := len(t.Heads)
	for _, row := range t.Rows {
		if len(row.Cells) > colCount {
			colCount = len(row.Cells)
		}
		for _, cell := range row.Cells {
			cell.Init(t.opts)
		}
	}

	// 初始化列宽数组
	colWidths := make([]int, colCount)

	// 计算表头列宽
	for i, head := range t.Heads {
		head.Init(t.opts)
		width := t.calcColWidth(head.Width, i)
		if width > colWidths[i] {
			colWidths[i] = width
		}
	}

	// 如果需要排序，对行进行排序
	sortColIdx := t.opts.SortColumn
	if sortColIdx >= 0 && sortColIdx < colCount {
		sort.SliceStable(t.Rows, func(i, j int) bool {
			// 确保索引有效
			if sortColIdx >= len(t.Rows[i].Cells) || sortColIdx >= len(t.Rows[j].Cells) {
				return false
			}

			valI := t.Rows[i].Cells[sortColIdx].String()
			valJ := t.Rows[j].Cells[sortColIdx].String()

			if t.opts.SortAscending {
				return valI < valJ
			}
			return valI > valJ
		})
	}

	// 如果显示行号，为每行添加行号单元格
	if t.opts.ShowRowNumber {
		for i, row := range t.Rows {
			rowNumCell := &Cell{Align: strutil.PosAuto, Val: fmt.Sprintf("%d", i)}
			row.Cells = append([]*Cell{rowNumCell}, row.Cells...)
		}
	}

	// 计算数据列宽
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			cellWidth := t.calcColWidth(cell.width, i)
			if cellWidth > colWidths[i] {
				colWidths[i] = cellWidth
			}
		}
	}

	// 为 Cell 设置宽度
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			if i < len(colWidths) {
				cell.width = colWidths[i]
			}
		}
	}

	// 保存计算后的列宽到 table 实例
	t.colWidths = colWidths
}

// Format as string
func (t *Table) formatHeader() {
	if len(t.Heads) == 0 && len(t.Rows) == 0 {
		return // 没有表头和数据，直接返回
	}

	buf := t.Buffer()
	opts := t.opts
	style := t.opts.Style

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
		buf.WriteRune(style.Border.Right) // 左边框
		var coloredHead string

		for i, head := range t.Heads {
			headStr := head.String()
			// 添加列填充
			if opts.CellPadding != "" {
				headStr = opts.CellPadding + headStr + opts.CellPadding
			}

			if i < len(t.colWidths) {
				// 使用 strutil.Resize 来对齐表头内容
				resized := strutil.Resize(headStr, t.colWidths[i], opts.Alignment)
				// 应用颜色（优先使用 FirstColor 给第一列）
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
				buf.WriteString(headStr)
			}

			if i < len(t.Heads)-1 { // 不是最后一个元素
				buf.WriteRune(style.Border.Cell) // 列分隔符
			}
		}

		buf.WriteRune(style.Border.Right) // 右边框
		buf.WriteByte('\n')

		// 画表头分隔线（如果需要）
		if opts.HeadBorder {
			t.drawBorderLine(buf, style.Divider.Left, style.Border.Center, style.Divider.Intersect, style.Divider.Right)
		} else if opts.RowBorder {
			t.drawBorderLine(buf, style.Border.Right, style.Border.Center, style.Border.Cell, style.Border.Right)
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

				// 添加列填充
				if opts.CellPadding != "" {
					cellStr = opts.CellPadding + cellStr + opts.CellPadding
				}

				// 根据宽度调整内容
				if cell.width > 0 {
					// 截断模式
					if opts.OverflowFlag == OverflowCut && cell.valWidth > cell.width {
						cellStr = strutil.Utf8Truncate(cellStr, cell.width, "")
					} else {
						// 填充至 cell.width 宽度
						cellStr = strutil.Resize(cellStr, cell.width, cell.Align)
					}
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
		buf.WriteByte('\n')

		// 画行分隔线（如果需要）
		if opts.RowBorder && i < len(t.Rows)-1 {
			t.drawBorderLine(buf, style.Border.Right, style.Border.Center, style.Divider.Intersect, style.Border.Right)
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

// drawBorderLine draws a borderline with the given characters
func (t *Table) drawBorderLine(buf *bytes.Buffer, leftChar, centerChar, intersect, rightChar rune) {
	if leftChar == 0 && rightChar == 0 {
		return // 如果没有边框字符，则跳过
	}

	buf.WriteRune(leftChar)

	for i, width := range t.colWidths {
		for j := 0; j < width; j++ {
			buf.WriteRune(centerChar)
		}
		if i < len(t.colWidths)-1 { // 不是最后一个列
			buf.WriteRune(intersect) // 列间分隔符
		}
	}

	buf.WriteRune(rightChar)
	buf.WriteByte('\n')
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

	// Height is the height of the row.
	//
	// Defaults to 0 - the height of the tallest cell(最高的单元格的高度)
	Height int
	// Separator for table columns
	// Separator rune
}

// reset the row context info
func (r *Row) reset() {
	for _, cell := range r.Cells {
		cell.init = false
		cell.width = 0
		cell.height = 0
	}
}

// Cell represents a column in a row
type Cell struct {
	// Width is the width of the cell. use for calc width
	Width int
	// Wrap when true wraps the contents of the cell when the length exceeds the width
	Wrap bool
	// Align when true aligns contents to the right
	Align strutil.PosFlag

	// TODO 支持 跨列，跨行 设置

	init bool
	// Val is the cell data
	Val any
	// string cache of Val
	str string
	// val content width
	valWidth int
	// width for the cell
	width  int
	height int
}

// NewCell creates a new cell with the given value
func NewCell(val any) *Cell {
	return &Cell{Align: strutil.PosAuto, Val: val}
}

// Init for one cell
func (c *Cell) Init(opts *Options) {
	if c.init {
		return
	}
	c.init = true

	if c.Align == strutil.PosAuto {
		c.Align = opts.Alignment
	}

	// conv value to string.
	s := c.String()
	// 去除空格
	if opts.TrimSpace {
		c.str = strings.TrimSpace(s)
	}
	c.calcWH()
}

// calc width and height of the cell
func (c *Cell) calcWH() {
	c.height = 0
	c.width = c.Width

	for _, s := range strings.Split(c.String(), "\n") {
		c.height++
		w := strutil.Utf8Width(s)
		if w > c.width {
			c.width = w
		}
	}
	c.valWidth = c.width
}

// String returns the string formatted representation of the cell
func (c *Cell) String() string {
	if c.str == "" {
		c.str = c.toString()
	}
	return c.str
}

func (c *Cell) toString() string {
	if c.Val == nil {
		return ""
	}
	if s, ok := c.Val.(string); ok {
		return s
	}
	return strutil.SafeString(c.Val)
}
