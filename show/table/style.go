package table

/*
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

	// BorderFlags control which borders to show using bit flags
	// e.g. BorderFlags = BorderTop | BorderBottom | BorderHeader
	BorderFlags uint8

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
	Center rune // eg: ─
	// Cell - column separator char. 列分隔符
	Cell rune
	Left rune

	BottomLeft, Bottom, BottomIntersect, BottomRight rune
}

// DividerStyle defines table divider style. 定义表格分隔器样式
type DividerStyle struct {
	Left  rune // eg: ├
	Right rune // eg: ┤
	// Intersect 交叉 eg: ┼
	Intersect rune
}

var (
	BorderStyleDefault = BorderStyle{
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
	}
)

var (
	/*
		StyleMySql - MySql-like table style
		+-----+------------+-----------+--------+
		|   # | FIRST NAME | LAST NAME | SALARY |
		+-----+------------+-----------+--------+
		|   1 | Arya       | Stark     |   3000 |
		|  20 | Jon        | Snow      |   2000 |
		| 300 | Tyrion     | Lannister |   5000 |
		+-----+------------+-----------+--------+
	*/
	StyleMySql = Style{
		HeadColor:  "info",
		BorderFlags: BorderAll,
		Border:      BorderStyleDefault,
		Divider: DividerStyle{
			Left:      '+',
			Right:     '+',
			Intersect: '+',
		},
	}

	// StyleSimple - Simple table style
	StyleSimple = Style{
		HeadColor:  "info",
		BorderFlags: BorderDefault,
		Border:      BorderStyleDefault,
		Divider: DividerStyle{
			Left:      '|',
			Right:     '|',
			Intersect: '+',
		},
	}

	// StyleMarkdown - Markdown table style
	StyleMarkdown = Style{
		HeadColor:  "info",
		BorderFlags: BorderAll,
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
		BorderFlags: BorderDefault,
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
		BorderFlags: BorderDefault,
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
		BorderFlags: BorderDefault,
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
		BorderFlags: BorderAll,
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

/*
StyleBoldBorder - table style with bold top and bottom lines:
━━━┯━━━━━━━┯━━━━━━━━━━━━━━━━━┯━━━━━━━━━━┯━━━━━━━━━━
 # │ pid   │ name            │ status   │ cpu
───┼───────┼─────────────────┼──────────┼──────────
 0 │   992 │ chrome          │ Sleeping │ 6.988768
 2 │ 13973 │ qemu-system-x86 │ Sleeping │ 4.996551
━━━┷━━━━━━━┷━━━━━━━━━━━━━━━━━┷━━━━━━━━━━┷━━━━━━━━━━
*/
var StyleBoldBorder = Style{
	HeadColor: "info",
	BorderFlags: BorderDefault,
	Border: BorderStyle{
		// Top - 使用粗线字符
		TopLeft:      '┏',
		Top:          '━',
		TopIntersect: '┯',
		TopRight:     '┓',
		// Body - 使用普通线字符
		Right:  '┃',
		Cell:   '│',
		Left:   '┃',
		Center: '─',
		// Bottom - 使用粗线字符
		BottomRight:     '┛',
		Bottom:          '━',
		BottomLeft:      '┗',
		BottomIntersect: '┷',
	},
	Divider: DividerStyle{
		Left:      '┃',
		Right:     '┃',
		Intersect: '┼',
	},
}