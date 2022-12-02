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

	HeadColor string
	RowColor  string
}

// BorderStyle for table
type BorderStyle struct {
	TopLeft, Top, TopIntersect, TopRight rune

	Right, Center, Cell, Left rune

	BottomRight, Bottom, BottomIntersect, BottomLeft rune
}

// DividerStyle defines table divider style
type DividerStyle struct {
	Left      rune
	Right     rune
	Intersect rune
}

var (
	// StyleDefault - MySql-like table style:
	StyleDefault = Style{
		HeadColor: "",
		RowColor:  "",
		Border: BorderStyle{
			// Top
			TopLeft:      '+',
			Top:          '-',
			TopIntersect: '+',
			TopRight:     '+',
			// Body
			Right: '|',
			Cell:  '|',
			// Bottom
			BottomRight:     '+',
			Bottom:          '-',
			BottomLeft:      '+',
			BottomIntersect: '+',
		},
		Divider: DividerStyle{},
	}

	StyleSimple   = Style{}
	StyleMarkdown = Style{}
)
