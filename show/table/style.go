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
type Style struct {
	Border  BorderStyle
	Divider DividerStyle

	TitleColor string
	HeadColor  string
	RowColor   string
}

// BorderStyle for table
type BorderStyle struct {
	TopLeft, Top, TopIntersect, TopRight rune

	Right, Center, Cell, Left rune

	BottomRight, Bottom, BottomIntersect, BottomLeft rune
}

// DividerStyle defines table divider style
type DividerStyle struct {
	Left  rune
	Right rune
	// Intersect 交叉
	Intersect rune
}

var (
	// StyleDefault - MySql-like table style:
	StyleDefault = Style{
		HeadColor: "info",
		RowColor:  "",
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

	// StyleSimple - Simple table style:
	StyleSimple = Style{
		HeadColor: "info",
		RowColor:  "",
		Border: BorderStyle{
			// Simple style without corners
			TopLeft:      ' ', // No corner
			Top:          ' ', // No top line
			TopIntersect: ' ', // No intersection
			TopRight:     ' ', // No corner
			// Body
			Right:  '|',
			Cell:   '|',
			Left:   '|',
			Center: '-',
			// Bottom
			BottomRight:     ' ', // No corner
			Bottom:          ' ', // No bottom line
			BottomLeft:      ' ', // No corner
			BottomIntersect: ' ', // No intersection
		},
		Divider: DividerStyle{
			Left:      '|',
			Right:     '|',
			Intersect: '|',
		},
	}

	// StyleMarkdown - Markdown table style:
	StyleMarkdown = Style{
		HeadColor: "info",
		RowColor:  "",
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
)
