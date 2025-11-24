package table_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gookit/gcli/v3/show/table"
)

func TestNewTable(t *testing.T) {
	tb := table.New("Table example1")
	tb.AddHead("Name", "Age", "City").
		AddRow("Tom", 25, "New York").
		AddRow("Jerry", 30, "Boston").
		AddRow("Alice", 28, "Chicago")

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

	tb.AddHead("Name", "Age", "Job").SetRows(users)

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
	tb.AddHead("ID", "Name").
		AddRow(1, "Product A").
		AddRow(2, "Product B")

	// 测试不同样式
	tb.WithOptions(func(opts *table.Options) {
		opts.Style = table.StyleSimple
		opts.HeadColor = "info"
	})

	result := tb.String()
	fmt.Println(result)

	if !strings.Contains(result, "Styled Table") {
		t.Error("Table should contain title")
	}
}

func TestTableMarkdownStyle(t *testing.T) {
	tb := table.New("Markdown Table")
	tb.AddHead("Name", "Value").
		AddRow("A", 100).
		AddRow("B", 200)

	tb.WithOptions(func(opts *table.Options) {
		opts.Style = table.StyleMarkdown
	})

	result := tb.String()
	fmt.Println(result)

	// 验证 Markdown 表格格式
	if !strings.Contains(result, "| Name | Value |") {
		t.Error("Markdown table should have header row")
	}
	if !strings.Contains(result, "| --- | --- |") {
		t.Error("Markdown table should have separator row")
	}
}
