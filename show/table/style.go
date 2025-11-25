package table

// Style for table
/*
	━━━┯━━━━━━━┯━━━━━━━━━━━━━━━━━┯━━━━━━━━━━┯━━━━━━━━━━
	 # │ pid   │ name            │ status   │ cpu
	───┼───────┼─────────────────┼──────────┼──────────
	 0 │   992 │ chrome          │ Sleeping │ 6.988768
	 2 │ 13973 │ qemu-system-x86 │ Sleeping │ 4.996551
	━━━┷━━━━━━━┷━━━━━━━━━━━━━━━━━┷━━━━━━━━━━┷━━━━━━━━━━

	+-----+------------+-----------+--------+-----------------------------+
	|   # | FIRST NAME | LAST NAME | SALARY |                             |
	+-----+------------+-----------+--------+-----------------------------+
	|   1 | Arya       | Stark     |   3000 |                             |
	|  20 | Jon        | Snow      |   2000 | You know nothing, Jon Snow! |
	| 300 | Tyrion     | Lannister |   5000 |                             |
	+-----+------------+-----------+--------+-----------------------------+
	|     |            | TOTAL     |  10000 |                             |
	+-----+------------+-----------+--------+-----------------------------+
*/

// Style for table
type Style struct {
	Border  BorderStyle
	Divider DividerStyle

	// TitleColor - table title color.
	//  allow: red, green, blue, yellow, magenta, cyan, white, gray, black
	TitleColor string
	HeadColor  string
	RowColor   string
	// FirstColor - first column color
	FirstColor string
}

// BorderStyle for table
type BorderStyle struct {
	TopLeft, Top, TopIntersect, TopRight rune

	// Right - right border char
	Right  rune
	Center rune
	// Cell - column separator char. 列分隔符
	Cell rune
	Left rune

	BottomRight, Bottom, BottomIntersect, BottomLeft rune
}

// DividerStyle defines table divider style. 定义表格分隔器样式
type DividerStyle struct {
	Left  rune
	Right rune
	// Intersect 交叉
	Intersect rune
}

var (
	// StyleDefault - MySql-like table style
	StyleDefault = Style{
		HeadColor:  "info",
		Border: BorderStyle{
			// Top
			TopLeft:      '+',
			Top:          '-',
			TopIntersect: '+',
			TopRight:     '+',
			// Body
			Right:  '|',
			Cell:   '|',
			Left:   '|',
			Center: '-',
			// Bottom
			BottomRight:     '+',
			Bottom:          '-',
			BottomLeft:      '+',
			BottomIntersect: '+',
		},
		Divider: DividerStyle{
			Left:      '+',
			Right:     '+',
			Intersect: '+',
		},
	}

	// StyleSimple - Simple table style
	StyleSimple = Style{
		HeadColor:  "info",
		Border: BorderStyle{
			// Simple style without corners
			TopLeft:      0, // No corner
			Top:          0, // No top line
			TopIntersect: 0, // No intersection
			TopRight:     0, // No corner
			// Body
			Right:  '|',
			Cell:   '|',
			Left:   '|',
			Center: '-',
			// Bottom
			BottomRight:     0, // No corner
			Bottom:          0, // No bottom line
			BottomLeft:      0, // No corner
			BottomIntersect: 0, // No intersection
		},
		Divider: DividerStyle{
			Left:      '|',
			Right:     '|',
			Intersect: '+',
		},
	}

	// StyleMarkdown - Markdown table style
	StyleMarkdown = Style{
		HeadColor:  "info",
		Border: BorderStyle{
			// Markdown doesn't have corners since it uses text characters
			TopLeft:      0, // Not used in markdown
			Top:          0, // Not used in markdown
			TopIntersect: 0, // Not used in markdown
			TopRight:     0, // Not used in markdown
			// Body
			Right:  '|',
			Cell:   '|',
			Left:   '|',
			Center: '-', // Used for header separator
			// Bottom
			BottomRight:     0, // Not used in markdown
			Bottom:          0, // Not used in markdown
			BottomLeft:      0, // Not used in markdown
			BottomIntersect: 0, // Not used in markdown
		},
		Divider: DividerStyle{
			Left:      '|',
			Right:     '|',
			Intersect: '|',
		},
	}

	// StyleBold - Bold table style with thick borders:
	StyleBold = Style{
		HeadColor:  "info",
		Border: BorderStyle{
			// Top
			TopLeft:      '┏',
			Top:          '━',
			TopIntersect: '┳',
			TopRight:     '┓',
			// Body
			Right:  '┃',
			Cell:   '┃',
			Left:   '┃',
			Center: '━',
			// Bottom
			BottomRight:     '┛',
			Bottom:          '━',
			BottomLeft:      '┗',
			BottomIntersect: '┻',
		},
		Divider: DividerStyle{
			Left:      '┣',
			Right:     '┫',
			Intersect: '╋',
		},
	}

	// StyleRounded - Rounded corner table style:
	StyleRounded = Style{
		HeadColor:  "info",
		Border: BorderStyle{
			// Top
			TopLeft:      '╭',
			Top:          '─',
			TopIntersect: '┬',
			TopRight:     '╮',
			// Body
			Right:  '│',
			Cell:   '│',
			Left:   '│',
			Center: '─',
			// Bottom
			BottomRight:     '╯',
			Bottom:          '─',
			BottomLeft:      '╰',
			BottomIntersect: '┴',
		},
		Divider: DividerStyle{
			Left:      '├',
			Right:     '┤',
			Intersect: '┼',
		},
	}

	// StyleDouble - Double line table style:
	StyleDouble = Style{
		HeadColor:  "info",
		Border: BorderStyle{
			// Top
			TopLeft:      '╔',
			Top:          '═',
			TopIntersect: '╦',
			TopRight:     '╗',
			// Body
			Right:  '║',
			Cell:   '║',
			Left:   '║',
			Center: '═',
			// Bottom
			BottomRight:     '╝',
			Bottom:          '═',
			BottomLeft:      '╚',
			BottomIntersect: '╩',
		},
		Divider: DividerStyle{
			Left:      '╠',
			Right:     '╣',
			Intersect: '╬',
		},
	}

	// StyleMinimal - Minimal table style with light borders:
	StyleMinimal = Style{
		HeadColor:  "info",
		Border: BorderStyle{
			// Top
			TopLeft:      '┌',
			Top:          '─',
			TopIntersect: '┬',
			TopRight:     '┐',
			// Body
			Right:  '│',
			Cell:   '│',
			Left:   '│',
			Center: '─',
			// Bottom
			BottomRight:     '┘',
			Bottom:          '─',
			BottomLeft:      '└',
			BottomIntersect: '┴',
		},
		Divider: DividerStyle{
			Left:      '├',
			Right:     '┤',
			Intersect: '┼',
		},
	}
)
