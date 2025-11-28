package table_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3/show/table"
	"github.com/gookit/goutil/testutil/assert"
)

func TestNewTable(t *testing.T) {
	tb := table.New("Table example1")
	tb.SetHeads("Name", "Age", "City").
		AddRow("Tom", 25, "New York").
		AddRow("Jerry", 30, "Boston").
		AddRow("Alice", 28, "Chicago").
		WithOptions(
			table.WithShowRowNumber(false),
			table.WithOverflowFlag(0),
			table.WithTrimSpace(true),
			table.WithSortColumn(-1, false),
			table.WithCSVOutput(false),
		)

	result := tb.String()

	// 验证表格包含基本元素
	if !strings.Contains(result, "Table example1") {
		t.Error("Table should contain title")
	}
	if !strings.Contains(result, "Name") || !strings.Contains(result, "Age") {
		t.Error("Table should contain headers")
	}
	if !strings.Contains(result, "Tom") || !strings.Contains(result, "Jerry") {
		t.Error("Table should contain data")
	}

	// 输出结果供检查
	fmt.Println(result)
}

func TestTableSetRowsWithSlice(t *testing.T) {
	tb := table.New("User Table")

	users := [][]any{
		{"Tom", 25, "Engineer"},
		{"Jerry", 30, "Designer"},
	}

	tb.SetHeads("Name", "Age", "Job").SetRows(users).
		WithOptions(table.WithStyle(table.StyleBold), table.WithCellPadding(" "))

	result := tb.String()

	if !strings.Contains(result, "User Table") {
		t.Error("Table should contain title")
	}
	if !strings.Contains(result, "Tom") || !strings.Contains(result, "Jerry") {
		t.Error("Table should contain user data")
	}

	fmt.Println(result)
}

func TestTableWithStyle(t *testing.T) {
	tb := table.New("Styled Table")
	tb.SetHeads("ID", "Name", "Age", "Description").
		AddRow(1, "Product A", 10, "This is a description").
		AddRow(2, "Product B", 20, "This is another description")

	// 测试不同样式
	tb.WithOptions(func(opts *table.Options) {
		opts.Style = table.StyleDefault
		opts.HeadColor = "info"
	})

	t.Run("StyleDefault", func(t *testing.T) {
		result := tb.String()
		fmt.Println(result)
		assert.StrContains(t, result, "Styled Table")
	})

	t.Run("StyleSimple", func(t *testing.T) {
		tb.WithOptions(table.WithStyle(table.StyleSimple))
		tb.Println()
	})
	t.Run("StyleMarkdown", func(t *testing.T) {
		tb.WithOptions(table.WithStyle(table.StyleMarkdown))
		tb.Println()
	})
	t.Run("StyleBold", func(t *testing.T) {
		tb.WithOptions(table.WithStyle(table.StyleBold))
		tb.Println()
	})
	t.Run("StyleBoldBorder", func(t *testing.T) {
		tb.WithOptions(table.WithStyle(table.StyleBoldBorder))
		tb.Println()
	})
	t.Run("StyleRounded", func(t *testing.T) {
		tb.WithOptions(table.WithStyle(table.StyleRounded))
		tb.Println()
	})
	t.Run("StyleDouble", func(t *testing.T) {
		tb.WithOptions(table.WithStyle(table.StyleDouble))
		tb.Println()
	})
	t.Run("StyleMinimal", func(t *testing.T) {
		tb.WithOptions(table.WithStyle(table.StyleMinimal))
		tb.Println()
	})
}

func TestTableMarkdownStyle(t *testing.T) {
	tb := table.New("Markdown Table")
	tb.SetHeads("Name", "Value").
		AddRow("A", 100).
		AddRow("B", 200)

	tb.WithOptions(table.WithStyle(table.StyleMarkdown))

	result := tb.String()
	fmt.Println(result)

	// 验证 Markdown 表格格式
	assert.StrContainsAll(t, result, []string{"|----|-----|", "|Name|Value|"})
}

// test for ColMaxWidth and OverflowFlag
func TestTableColMaxWidthAndOverflowFlag(t *testing.T) {
	tb := table.New("ColMaxWidth and OverflowFlag")
	tb.SetHeads("Name", "Description").
		AddRow("Long Name", "This is a long description that exceeds the column width.").
		AddRow("Short Name", "This is a short description.")
	tb.WithOptions(
		table.WithColMaxWidth(30),
		table.WithColumnWidths(20, 50),
		table.WithOverflowFlag(table.OverflowCut),
		func(opts *table.Options) {
			opts.RowBorder = true
		},
	)

	result := tb.String()
	fmt.Println(result)
	assert.StrContainsAll(t, result, []string{"Name", "Description"})
}
