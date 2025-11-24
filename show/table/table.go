package table

import (
	"fmt"
	"io"
	"strings"

	"github.com/gookit/gcli/v3/show"
	"github.com/gookit/goutil/comdef"
	"github.com/gookit/goutil/strutil"
)

// Options struct
type Options struct {
	Style
	HeadColor string

	Alignment   strutil.PosFlag
	ColMaxWidth int
	LineNumber  bool
	WrapContent bool

	// HasBorder show borderline
	HasBorder bool
	// RowBorder show row border
	RowBorder bool
	// HeadBorder show head border
	HeadBorder bool
	// WrapBorder wrap border for table
	WrapBorder bool
}

// OpFunc define
type OpFunc func(opts *Options)

// Table a cli Table show
type Table struct {
	show.Base // use for internal
	// options ...
	opts *Options
	out  comdef.ByteStringWriter

	// Title for the table
	Title string
	// Heads the table head data
	Heads []string
	// Rows table data rows
	Rows []*Row

	// column value align type.
	// key is col index. start from 0.
	colAlign map[int]strutil.PosFlag
}

// New create table
func New(title string, fns ...OpFunc) *Table {
	t := &Table{
		Title: title,
		opts: &Options{
			Style: StyleDefault,
		},
	}

	return t.WithOptions(fns...)
}

// WithOptions for table
func (t *Table) WithOptions(fns ...OpFunc) *Table {
	for _, fn := range fns {
		fn(t.opts)
	}
	return t
}

// AddHead column names to table
func (t *Table) AddHead(names ...string) *Table {
	t.Heads = names
	return t
}

// AddRow data to table
func (t *Table) AddRow(cols ...any) {
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
}

// SetRows to table
func (t *Table) SetRows(rs any) *Table {

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
	t.prepare()

	t.formatHeader()

	t.formatBody()

	t.formatFooter()

	panic("implement me")
}

func (t *Table) prepare() {

	// determine the width for each column (cell in a row)
	var colWidths []int
	for _, row := range t.Rows {
		for i, cell := range row.Cells {
			// resize colwidth array
			if i+1 > len(colWidths) {
				colWidths = append(colWidths, 0)
			}

			cellWidth := cell.MaxWidth()
			if t.opts.ColMaxWidth != 0 && cellWidth > t.opts.ColMaxWidth {
				cellWidth = t.opts.ColMaxWidth
			}

			if cellWidth > colWidths[i] {
				colWidths[i] = cellWidth
			}
		}
	}
}

// Format as string
func (t *Table) formatHeader() {
	panic("implement me")
}

// Format as string
func (t *Table) formatBody() {
	for _, row := range t.Rows {
		fmt.Println(row)
	}

	panic("implement me")
}

// Format as string
func (t *Table) formatFooter() {
	panic("implement me")
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
